//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"sync"
)

type p24CapturedKeyedCommand struct {
	args   []string
	input  []byte
	stdout string
}

var p24KeyedCaptures = struct {
	sync.Mutex
	byWorkspace map[string]map[string]p24CapturedKeyedCommand
}{byWorkspace: map[string]map[string]p24CapturedKeyedCommand{}}

func p24StartKeyedCapture(workspace string) {
	p24KeyedCaptures.Lock()
	defer p24KeyedCaptures.Unlock()
	p24KeyedCaptures.byWorkspace[workspace] = map[string]p24CapturedKeyedCommand{}
}

func p24EndKeyedCapture(workspace string) map[string]p24CapturedKeyedCommand {
	p24KeyedCaptures.Lock()
	defer p24KeyedCaptures.Unlock()
	got := p24KeyedCaptures.byWorkspace[workspace]
	delete(p24KeyedCaptures.byWorkspace, workspace)
	return got
}

func p24CaptureKeyedCommand(workspace string, args []string, stdout, stderr string, runErr error) {
	if runErr != nil || stderr != "" {
		return
	}
	p24KeyedCaptures.Lock()
	defer p24KeyedCaptures.Unlock()
	commands, active := p24KeyedCaptures.byWorkspace[workspace]
	if !active {
		return
	}
	var envelope struct {
		Operation string `json:"operation"`
		OK        bool   `json:"ok"`
	}
	if json.Unmarshal([]byte(stdout), &envelope) != nil || !envelope.OK || envelope.Operation == "" {
		return
	}
	inputPath, keyed := "", false
	for i, arg := range args {
		if arg == "--idempotency-key" && i+1 < len(args) {
			keyed = true
		}
		if arg == "--input" && i+1 < len(args) {
			inputPath = args[i+1]
		}
	}
	if !keyed || commands[envelope.Operation].args != nil {
		return
	}
	var input []byte
	if inputPath != "" {
		input, _ = os.ReadFile(inputPath)
	}
	commands[envelope.Operation] = p24CapturedKeyedCommand{args: append([]string(nil), args...), input: input, stdout: stdout}
}
