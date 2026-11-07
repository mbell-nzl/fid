package backup

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
)

// Commit identity for backups. fid sets it itself so a backup never depends
// on the user's git config, which CI and cron may not have.
const (
	commitName  = "fid backup"
	commitEmail = "backup@fid.invalid"
)

// repositoryVariables name another repository or change how git reads one.
// They are git's local environment variables (git rev-parse
// --local-env-vars), plus quarantine and namespace. fid may run where git
// set them, such as a hook, so its own git commands drop them.
var repositoryVariables = map[string]bool{
	"GIT_ALTERNATE_OBJECT_DIRECTORIES": true, "GIT_CONFIG": true, "GIT_CONFIG_PARAMETERS": true, "GIT_CONFIG_COUNT": true,
	"GIT_OBJECT_DIRECTORY": true, "GIT_DIR": true, "GIT_WORK_TREE": true, "GIT_IMPLICIT_WORK_TREE": true, "GIT_GRAFT_FILE": true,
	"GIT_INDEX_FILE": true, "GIT_NO_REPLACE_OBJECTS": true, "GIT_REPLACE_REF_BASE": true, "GIT_PREFIX": true,
	"GIT_SHALLOW_FILE": true, "GIT_COMMON_DIR": true, "GIT_QUARANTINE_PATH": true, "GIT_NAMESPACE": true,
}

// gitEnvironment is the process environment less repository variables,
// with fid's commit identity and no terminal prompts: a backup run by the
// loop or cron fails rather than waiting for a password. Credential helpers,
// SSH and URL rewriting from the user's git config still apply. fid's values
// come last, and os/exec uses the last of duplicate keys, so they win over
// inherited ones.
func gitEnvironment(environment []string) []string {
	result := make([]string, 0, len(environment)+5)
	for _, entry := range environment {
		if name, _, _ := strings.Cut(entry, "="); !repositoryVariables[name] {
			result = append(result, entry)
		}
	}
	return append(result, "GIT_TERMINAL_PROMPT=0",
		"GIT_AUTHOR_NAME="+commitName, "GIT_AUTHOR_EMAIL="+commitEmail,
		"GIT_COMMITTER_NAME="+commitName, "GIT_COMMITTER_EMAIL="+commitEmail)
}

// clone runs git in fid's bare clone. Every command names the clone
// explicitly, so git never searches for a repository, and overrides the
// settings in the user's git config that could run code or detach work:
// hooks, and background garbage collection.
type clone struct {
	gitDir      string
	workingDir  string
	environment []string
}

// gitError is a failed git command with what it wrote to stderr.
type gitError struct {
	arguments []string
	stdout    string
	stderr    string
	err       error
}

func (err *gitError) Error() string {
	message := condense(err.stderr)
	if message == "" {
		message = err.err.Error()
	}
	return message
}

func (err *gitError) Unwrap() error { return err.err }

func (repository clone) run(stdin io.Reader, arguments ...string) (string, error) {
	return runGit(repository.workingDir, repository.environment, stdin, append([]string{"--git-dir=" + repository.gitDir}, arguments...)...)
}

// runGit runs git with the overrides every backup command needs.
func runGit(directory string, environment []string, stdin io.Reader, arguments ...string) (string, error) {
	arguments = append([]string{"-c", "core.hooksPath=/dev/null", "-c", "gc.autoDetach=false", "-c", "maintenance.autoDetach=false"}, arguments...)
	command := exec.Command("git", arguments...)
	command.Dir = directory
	command.Env = environment
	command.Stdin = stdin
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return stdout.String(), &gitError{arguments: arguments, stdout: stdout.String(), stderr: stderr.String(), err: err}
	}
	return stdout.String(), nil
}

// text runs git and returns its output less one trailing newline.
func (repository clone) text(arguments ...string) (string, error) {
	return repository.output(nil, arguments...)
}

// output is text with stdin.
func (repository clone) output(stdin io.Reader, arguments ...string) (string, error) {
	output, err := repository.run(stdin, arguments...)
	return strings.TrimSuffix(output, "\n"), err
}

// condense keeps git's messages short: its lines without hints, joined.
func condense(stderr string) string {
	lines := []string{}
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "hint:") {
			lines = append(lines, line)
		}
	}
	return bounded(strings.Join(lines, "; "))
}

// rejectedPush returns git push --porcelain's account of a rejected ref, or
// "" when the push failed before the remote could judge it.
func rejectedPush(err error) string {
	var failed *gitError
	if !errors.As(err, &failed) {
		return ""
	}
	for _, line := range strings.Split(failed.stdout, "\n") {
		if fields := strings.SplitN(line, "\t", 3); len(fields) == 3 && fields[0] == "!" {
			remote := []string{}
			for _, message := range strings.Split(failed.stderr, "\n") {
				if strings.HasPrefix(message, "remote:") {
					remote = append(remote, strings.TrimSpace(strings.TrimPrefix(message, "remote:")))
				}
			}
			return bounded(strings.TrimSpace(strings.Join(append([]string{fields[2]}, remote...), "; ")))
		}
	}
	return ""
}
