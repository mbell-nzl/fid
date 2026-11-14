package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/leaktest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// V2 / check:leak-test covers stored text, the entire help tree, and live
// human commands including refusals. JSON and verbose diagnostics are exempt.
func TestV2HumanOutputNamesNoKernelWords(t *testing.T) {
	t.Parallel()
	t.Run("goldens", checkTextGoldensForKernelWords)
	t.Run("help", func(t *testing.T) {
		root := newRootCommandAt(t.TempDir())
		root.InitDefaultHelpCmd()
		root.InitDefaultCompletionCmd()
		checkHelpForKernelWords(t, root, nil)
	})
	t.Run("journey", checkHumanJourneyForKernelWords)
}

func checkTextGoldensForKernelWords(t *testing.T) {
	t.Helper()
	count := 0
	err := filepath.WalkDir("testdata", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".txt" {
			return nil
		}
		count++
		t.Run(filepath.ToSlash(path), func(t *testing.T) {
			output, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			leaktest.AssertNoKernelWords(t, string(output))
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("no text goldens checked")
	}
	t.Logf("checked %d text goldens", count)
}

func checkHelpForKernelWords(t *testing.T, command *cobra.Command, path []string) {
	t.Helper()
	t.Run(command.Name(), func(t *testing.T) {
		command.InitDefaultHelpFlag()
		for _, field := range []struct{ name, text string }{
			{"Use", command.Use}, {"Short", command.Short}, {"Long", command.Long}, {"Example", command.Example},
		} {
			t.Run(field.name, func(t *testing.T) {
				leaktest.AssertNoKernelWords(t, helpWithoutClaimFlag(field.text))
			})
		}
		for _, flags := range []*pflag.FlagSet{command.LocalFlags(), command.InheritedFlags()} {
			flags.VisitAll(func(flag *pflag.Flag) {
				t.Run("flag/"+flag.Name, func(t *testing.T) {
					leaktest.AssertNoKernelWords(t, helpWithoutClaimFlag(flag.Usage))
				})
			})
		}
		t.Run("rendered", func(t *testing.T) {
			got := executeCommandForTest(t.TempDir(), append(append([]string{}, path...), "--help")...)
			if got.err != nil || got.stdout == "" || got.stderr != "" {
				t.Fatalf("%s --help = %#v", command.CommandPath(), got)
			}
			leaktest.AssertNoKernelWords(t, helpWithoutClaimFlag(got.stdout))
			leaktest.AssertNoKernelWords(t, helpWithoutClaimFlag(got.stderr))
		})
		for _, child := range command.Commands() {
			checkHelpForKernelWords(t, child, append(append([]string{}, path...), child.Name()))
		}
	})
}

// Only the exact flag token is exempt (ADR-0014). In particular, this does
// not hide "claim" in surrounding prose or in a different flag spelling.
func helpWithoutClaimFlag(text string) string {
	return helpTokens.ReplaceAllStringFunc(text, func(token string) string {
		if token == "--claim" {
			return "--take"
		}
		return token
	})
}

var helpTokens = regexp.MustCompile(`[[:alnum:]_-]+`)

func TestHelpExemptsOnlyTheExactClaimFlag(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ input, want string }{
		{"--claim", "--take"},
		{"Use '--claim' or --claim=true; --claim --claim.", "Use '--take' or --take=true; --take --take."},
		{"--claim takes a claim", "--take takes a claim"},
		{"claim claims --claims --claim-extra ---claim prefix--claim --claim_suffix", "claim claims --claims --claim-extra ---claim prefix--claim --claim_suffix"},
	} {
		t.Run(row.input, func(t *testing.T) {
			t.Parallel()
			if got := helpWithoutClaimFlag(row.input); got != row.want {
				t.Fatalf("masked help = %q, want %q", got, row.want)
			}
		})
	}
}

func assertHumanCommand(t *testing.T, got commandExecutionResult, exit int) {
	t.Helper()
	leaktest.AssertNoKernelWords(t, got.stdout)
	leaktest.AssertNoKernelWords(t, got.stderr)
	code := 0
	if got.err != nil {
		code = cli.ExitCode(got.err)
	}
	if code != exit {
		t.Fatalf("command = %#v, want exit %d (got %d)", got, exit, code)
	}
	if exit != 0 && !strings.HasPrefix(got.stderr, "error: ") {
		t.Fatalf("refusal has no human error: %#v", got)
	}
}
