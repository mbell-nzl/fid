//go:build integration && linux

package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

// The keyed corpus must exercise an actual person's terminal, including retries
// after CHANGES returned the task to the builder. The integration binary itself
// is built without coverage; no helper test process is launched here.
func approvalCorpusCommand(t *testing.T, workspace, binary string, args ...string) (string, string, error) {
	t.Helper()
	task, action := "", "approve"
	for i, arg := range args {
		if (arg == "approve" || arg == "waive") && i+1 < len(args) {
			action = arg
			task = args[i+1]
			break
		}
	}
	if task == "" {
		t.Fatal("approval corpus needs a task")
	}
	leader, err := os.OpenFile("/dev/ptmx", os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = leader.Close() }()
	if err := unix.IoctlSetPointerInt(int(leader.Fd()), unix.TIOCSPTLCK, 0); err != nil {
		t.Fatal(err)
	}
	number, err := unix.IoctlGetUint32(int(leader.Fd()), unix.TIOCGPTN)
	if err != nil {
		t.Fatal(err)
	}
	follower, err := os.OpenFile(fmt.Sprintf("/dev/pts/%d", number), os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = follower.Close() }()
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Dir = workspace
	cmd.Stdin = follower
	cmd.Stderr = follower
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	reached := make(chan []byte, 1)
	go func() {
		var read bytes.Buffer
		sent := false
		chunk := make([]byte, 4096)
		for {
			n, err := leader.Read(chunk)
			if n > 0 {
				read.Write(chunk[:n])
				if !sent && strings.Contains(read.String(), "Type "+task+" to "+action+": ") {
					sent = true
					_, _ = io.WriteString(leader, task+"\n")
				}
			}
			if err != nil {
				break
			}
		}
		reached <- read.Bytes()
	}()
	err = cmd.Start()
	_ = follower.Close()
	if err == nil {
		err = cmd.Wait()
	}
	prompt := string(bytes.ReplaceAll(<-reached, []byte("\r\n"), []byte("\n")))
	stderr := prompt
	if err == nil {
		if prompt != "" && ((action == "approve" && !strings.Contains(prompt, "Diffstat:")) || !strings.Contains(prompt, "Type "+task+" to "+action+": ")) {
			t.Fatalf("unexpected approval prompt: %q", prompt)
		}
		// The prompt is a person's terminal output, separate from the JSON envelope
		// which the shared corpus captures and compares byte for byte.
		stderr = ""
	}
	p24CaptureKeyedCommand(workspace, args, stdout.String(), stderr, err)
	return stdout.String(), stderr, err
}
