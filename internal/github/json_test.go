package github

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestCIRejectsConflictingFieldAliases(t *testing.T) {
	t.Parallel()
	success := sampleRun(1, 1, 10, "completed", "success")
	page := string(ciPage(t, 1, []WorkflowRun{success}))
	for _, tt := range []struct {
		name, original, replacement string
	}{
		{"count alias last", `"total_count":1`, `"total_count":2,"TOTAL_COUNT":1`},
		{"count alias first", `"total_count":1`, `"TOTAL_COUNT":2,"total_count":1`},
		{"conclusion alias last", `"conclusion":"success"`, `"conclusion":"failure","Conclusion":"success"`},
		{"conclusion alias first", `"conclusion":"success"`, `"Conclusion":"failure","conclusion":"success"`},
		{"Unicode conclusion alias", `"conclusion":"success"`, `"conclusion":"failure","concluſion":"success"`},
		{"Unicode workflow alias", `"workflow_id":10`, `"workflow_id":10,"worKflow_id":11`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := strings.Replace(page, tt.original, tt.replacement, 1)
			client := clientForTest(t, func(context.Context, []string, []string) Result {
				return Result{Stdout: []byte(data)}
			})
			got, err := client.MergeCommitCI(context.Background(), mergeSHA)
			if err == nil {
				t.Fatalf("conflicting response accepted: CI = %#v", got)
			}
			assertUnavailable(t, err, MalformedResponse)
		})
	}
}

func TestPRRejectsNestedFieldAliases(t *testing.T) {
	t.Parallel()
	data := strings.Replace(string(prFixture(t, nil)), fmt.Sprintf(`"oid":%q`, mergeSHA), fmt.Sprintf(`"oid":%q,"OID":%q`, headSHA, mergeSHA), 1)
	client := clientForTest(t, func(context.Context, []string, []string) Result {
		return Result{Stdout: []byte(data)}
	})
	_, err := client.PR(context.Background(), prURL)
	assertUnavailable(t, err, MalformedResponse)
}
