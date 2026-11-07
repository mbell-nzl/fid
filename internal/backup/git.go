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

// gitEnvironment is the process environment less repository variables and
// inherited commit details, with fid's commit identity and no terminal
// prompts: a backup run by the loop or cron fails rather than waiting for a
// password. Credential helpers, SSH and URL rewriting from the user's git
// config still apply. Every inherited GIT_AUTHOR_ and GIT_COMMITTER_
// variable goes, not just those fid sets: GIT_AUTHOR_DATE and
// GIT_COMMITTER_DATE would otherwise date fid's commit, or stop it when they
// aren't dates. fid's values come last, and os/exec uses the last of
// duplicate keys, so they win over inherited ones.
func gitEnvironment(environment []string) []string {
	result := make([]string, 0, len(environment)+5)
	for _, entry := range environment {
		name, _, _ := strings.Cut(entry, "=")
		if !repositoryVariables[name] && !strings.HasPrefix(name, "GIT_AUTHOR_") && !strings.HasPrefix(name, "GIT_COMMITTER_") {
			result = append(result, entry)
		}
	}
	return append(result, "GIT_TERMINAL_PROMPT=0",
		"GIT_AUTHOR_NAME="+commitName, "GIT_AUTHOR_EMAIL="+commitEmail,
		"GIT_COMMITTER_NAME="+commitName, "GIT_COMMITTER_EMAIL="+commitEmail)
}

// clone runs git in a bare repository fid owns: its clone of the backup
// repository, or the temporary one a restore fetches into. Every command
// names the repository explicitly, so git never searches for one, and
// overrides the settings in the user's git config that could run code or
// detach work: hooks, and background garbage collection.
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

// stream runs git and writes its output to stdout as it comes, for blobs
// too large to hold in memory.
func (repository clone) stream(stdout io.Writer, arguments ...string) error {
	return streamGit(repository.workingDir, repository.environment, nil, stdout, append([]string{"--git-dir=" + repository.gitDir}, arguments...)...)
}

// runGit runs git with the overrides every backup command needs.
func runGit(directory string, environment []string, stdin io.Reader, arguments ...string) (string, error) {
	var stdout bytes.Buffer
	err := streamGit(directory, environment, stdin, &stdout, arguments...)
	var failed *gitError
	if errors.As(err, &failed) {
		failed.stdout = stdout.String()
	}
	return stdout.String(), err
}

// streamGit is runGit writing git's output to stdout.
func streamGit(directory string, environment []string, stdin io.Reader, stdout io.Writer, arguments ...string) error {
	arguments = append([]string{"-c", "core.hooksPath=/dev/null", "-c", "gc.autoDetach=false", "-c", "maintenance.autoDetach=false"}, arguments...)
	command := exec.Command("git", arguments...)
	command.Dir = directory
	command.Env = environment
	command.Stdin = stdin
	var stderr bytes.Buffer
	command.Stdout = stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return &gitError{arguments: arguments, stderr: stderr.String(), err: err}
	}
	return nil
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
