package cli

import "testing"

func TestExitCodeForCategory(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		category Category
		want     int
	}{
		{CategoryInvalidInput, 2},
		{CategoryNotFound, 3},
		{CategoryConflict, 4},
		{CategoryPrecondition, 4},
		{CategoryAuthority, 5},
		{CategoryIntegrity, 6},
		{CategoryInternal, 7},
		{CategoryGateUnsuccessful, 8},
	} {
		test := test
		t.Run(string(test.category), func(t *testing.T) {
			t.Parallel()
			if got := ExitCodeForCategory(test.category); got != test.want {
				t.Fatalf("ExitCodeForCategory(%q) = %d, want %d", test.category, got, test.want)
			}
		})
	}
}

func TestFailureErrorCarriesCodeAndCategory(t *testing.T) {
	t.Parallel()

	failure := Failure{Code: CodeCapabilityUnavailable, Category: CategoryAuthority}
	err := NewFailureError(failure)
	if err.Code != CodeCapabilityUnavailable || err.Category != CategoryAuthority {
		t.Fatalf("FailureError = %#v, want code/category from failure", err)
	}
	if got := ExitCode(err); got != 5 {
		t.Fatalf("ExitCode(FailureError) = %d, want 5", got)
	}
}
