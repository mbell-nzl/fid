package claims

import (
	"testing"
	"time"
)

func TestGenericClaimExpiresAtEqualityAndNeverBefore(t *testing.T) {
	expiresAt := time.Date(2026, time.September, 5, 12, 0, 0, 0, time.UTC)
	claim := GenericClaim{ExpiresAt: expiresAt.Format(time.RFC3339Nano)}
	for _, testCase := range []struct {
		name string
		now  time.Time
		want bool
	}{
		{name: "before expiry", now: expiresAt.Add(-time.Nanosecond), want: false},
		{name: "at expiry", now: expiresAt, want: true},
		{name: "after expiry", now: expiresAt.Add(time.Nanosecond), want: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			if got := claim.expiredAt(testCase.now); got != testCase.want {
				t.Fatalf("expiredAt(%s) = %t, want %t", testCase.now.Format(time.RFC3339Nano), got, testCase.want)
			}
		})
	}
}
