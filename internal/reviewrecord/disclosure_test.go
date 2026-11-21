package reviewrecord

import "testing"

func TestNoOtherModelReview(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, task, stream string
		current, other     bool
		want               bool
	}{
		{"same family", "task", "primary", true, false, true},
		{"primary other family", "task", "primary", true, true, false},
		{"separate other family", "task", "separate", true, true, false},
		{"stale other family", "task", "primary", false, true, true},
		{"another task", "elsewhere", "primary", true, true, true},
		{"another stream", "task", "elsewhere", true, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rounds := []RoundView{{Round: Round{TaskID: tc.task, Stream: tc.stream}, Current: tc.current, OtherModel: tc.other}}
			if got := NoOtherModelReview("task", rounds); got != tc.want {
				t.Fatalf("no other-model review=%t, want %t", got, tc.want)
			}
		})
	}
}
