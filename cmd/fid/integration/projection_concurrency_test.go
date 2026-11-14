//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestProjectionReadsDuringConcurrentWrites(t *testing.T) {
	t.Parallel()
	workspace, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
	p21Run(t, binary, workspace, "--json", "init", "--prefix", "FID")
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	const writes, readers, reads = 40, 3, 45
	start := make(chan struct{})
	writer := make(chan error, 1)
	reader := make(chan error, readers)
	go func() {
		<-start
		for i := range writes {
			_, err := concurrentProjectionFID(ctx, binary, workspace, "issue", "create", "--kind", "task",
				"--title", fmt.Sprintf("Concurrent issue %d", i), "--actor", "human:test",
				"--idempotency-key", fmt.Sprintf("concurrent-%d", i))
			if err != nil {
				writer <- err
				return
			}
		}
		writer <- nil
	}()
	for range readers {
		go func() {
			<-start
			for range reads {
				if _, err := concurrentProjectionFID(ctx, binary, workspace, "status"); err != nil {
					reader <- err
					return
				}
			}
			reader <- nil
		}()
	}
	close(start)
	if err := <-writer; err != nil {
		t.Errorf("writer: %v", err)
	}
	for range readers {
		if err := <-reader; err != nil {
			t.Errorf("reader: %v", err)
		}
	}
	if t.Failed() {
		return
	}
	listed := p21Run(t, binary, workspace, "--json", "issue", "list", "--kind", "task")
	issues := p20Objects(t, p20Object(t, listed["result"])["issues"])
	if len(issues) != writes {
		t.Fatalf("final issues = %d, want %d", len(issues), writes)
	}
	titles := make(map[string]bool, writes)
	for _, issue := range issues {
		titles[p21PathString(t, issue, "title")] = true
	}
	for i := range writes {
		if title := fmt.Sprintf("Concurrent issue %d", i); !titles[title] {
			t.Errorf("final state missing %q", title)
		}
	}
}

func concurrentProjectionFID(ctx context.Context, binary, workspace string, args ...string) (map[string]any, error) {
	command := exec.CommandContext(ctx, binary, append([]string{"--json"}, args...)...)
	command.Dir = workspace
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("fid %v: %w: %s", args, err, output)
	}
	var envelope map[string]any
	if err := json.Unmarshal(output, &envelope); err != nil {
		return nil, fmt.Errorf("decode fid %v: %w: %s", args, err, output)
	}
	if envelope["ok"] != true {
		return nil, fmt.Errorf("fid %v failed: %s", args, output)
	}
	return envelope, nil
}
