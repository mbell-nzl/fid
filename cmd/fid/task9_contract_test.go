package main

import (
	"strings"
	"testing"
)

func TestTask9RetainedInventorySurface(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "legacy policy command is absent",
			args:    []string{"policy", "show"},
			want:    "unknown command",
			wantErr: true,
		},
		{
			name: "workflow define accepts standalone actor and idempotency",
			args: []string{"--project", workspace, "--json", "workflow", "define", "--name", "task-9", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--actor", "agent:task-9", "--idempotency-key", "task-9-workflow"},
			want: `"idempotency_key":"task-9-workflow"`,
		},
		{
			name: "AC add accepts standalone actor and idempotency",
			args: func() []string {
				issueID := createTask8Issue(t, workspace, "Task 9 acceptance criterion")
				return []string{"--project", workspace, "--json", "issue", "ac", "add", issueID, "--text", "has an attributed inventory row", "--ordinal", "1", "--actor", "agent:task-9", "--idempotency-key", "task-9-ac"}
			}(),
			want: `"idempotency_key":"task-9-ac"`,
		},
		{
			name: "label add accepts standalone actor and idempotency",
			args: []string{"--project", workspace, "--json", "label", "add", "task-9-label", "--actor", "agent:task-9", "--idempotency-key", "task-9-label"},
			want: `"idempotency_key":"task-9-label"`,
		},
		{
			name: "issue label add accepts standalone actor and idempotency",
			args: func() []string {
				issueID := createTask8Issue(t, workspace, "Task 9 label attribution")
				return []string{"--project", workspace, "--json", "issue", "label", "add", issueID, "task-9-label", "--actor", "agent:task-9", "--idempotency-key", "task-9-issue-label"}
			}(),
			want: `"idempotency_key":"task-9-issue-label"`,
		},
		{
			name: "link add accepts standalone actor and idempotency",
			args: func() []string {
				fromID := createTask8Issue(t, workspace, "Task 9 link source")
				toID := createTask8Issue(t, workspace, "Task 9 link target")
				return []string{"--project", workspace, "--json", "link", "add", fromID, "blocks", toID, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:task-9", "--idempotency-key", "task-9-link"}
			}(),
			want: `"idempotency_key":"task-9-link"`,
		},
		{
			name:    "requirement sync is removed",
			args:    []string{"--project", workspace, "--json", "requirement", "sync"},
			want:    "unknown flag: --idempotency-key",
			wantErr: true,
		},
		{
			name:    "AC verification is fail closed",
			args:    []string{"--project", workspace, "--json", "issue", "ac", "verify", "FID-1", "AC-1", "--actor", "agent:task-9"},
			want:    `"code":"capability-unavailable"`,
			wantErr: true,
		},
		{
			name:    "concern raise is fail closed",
			args:    []string{"--project", workspace, "--json", "concern", "raise", "--on", "FID-1", "--actor", "agent:task-9"},
			want:    `"code":"capability-unavailable"`,
			wantErr: true,
		},
		{
			name:    "maintenance vacuum is fail closed",
			args:    []string{"--project", workspace, "--json", "maintenance", "vacuum", "--actor", "agent:task-9"},
			want:    `"code":"capability-unavailable"`,
			wantErr: true,
		},
		{
			name: "stats groups by historical actor attribution",
			args: []string{"--project", workspace, "--json", "stats", "tokens", "--group-by", "actor"},
			want: `"items"`,
		},
		{
			name: "activity normalizes retained event names",
			args: []string{"--project", workspace, "--json", "activity", "list", "--event", "issue-label-added"},
			want: `"issue-label-added"`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := executeCommandForTest(workspace, testCase.args...)
			text := result.stderr + result.stdout
			if result.err != nil {
				text += result.err.Error()
			}
			if (result.err != nil) != testCase.wantErr || !strings.Contains(text, testCase.want) {
				t.Fatalf("%q result = (stdout=%q, stderr=%q, err=%v), want %q", testCase.args, result.stdout, result.stderr, result.err, testCase.want)
			}
		})
	}
}
