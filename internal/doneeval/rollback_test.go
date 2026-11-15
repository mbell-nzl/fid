package doneeval

import "testing"

func TestRollbackHeading(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		body string
		want bool
	}{
		{"## Rollback", true}, {"   ##\tRollback ##  \r\nRevert", true},
		{"# Rollback", false}, {"### Rollback", false}, {"##Rollback", false}, {"## Rollback plan", false}, {"## rollback", false},
		{"    ## Rollback", false}, {"\t## Rollback", false}, {"> ## Rollback", false},
		{"```markdown\n## Rollback\n```", false}, {"~~~\n## Rollback\n~~~", false},
		{"````\n```\n## Rollback\n````", false}, {"~~~\n```\n## Rollback\n~~~", false},
		{"```\n## Rollback\n``` still fenced\n## Rollback", false},
		{"```\nexample\n```\n## Rollback", true}, {"~~~\nexample\n~~~~\n## Rollback", true},
	} {
		t.Run(tc.body, func(t *testing.T) {
			t.Parallel()
			s := completeSnapshot("assured", false)
			s.GitHub.PR.Body = tc.body
			got := Evaluate(s)
			if (got.Gates[RollbackGate].State == Satisfied) != tc.want {
				t.Fatalf("gate=%+v", got.Gates[RollbackGate])
			}
		})
	}
}
