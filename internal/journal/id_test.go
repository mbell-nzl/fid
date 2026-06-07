package journal

import (
	"regexp"
	"testing"
	"time"
)

func TestGenerateIDUsesDefaultSixHexChars(t *testing.T) {
	got, err := GenerateID("FID", nil)
	if err != nil {
		t.Fatalf("GenerateID() error = %v", err)
	}

	if !regexp.MustCompile(`^FID-[0-9a-f]{6}$`).MatchString(got) {
		t.Fatalf("GenerateID() = %q, want FID-<6 hex chars>", got)
	}
}

func TestGenerateIDWithLengthSupportsEightHexChars(t *testing.T) {
	got, err := GenerateIDWithLength("FID", 8, nil)
	if err != nil {
		t.Fatalf("GenerateIDWithLength() error = %v", err)
	}

	if !regexp.MustCompile(`^FID-[0-9a-f]{8}$`).MatchString(got) {
		t.Fatalf("GenerateIDWithLength() = %q, want FID-<8 hex chars>", got)
	}
}

func TestGenerateIDRetriesOnCollision(t *testing.T) {
	calls := 0
	restore := stubHashSuffixGenerator(func(length int) (string, error) {
		switch length {
		case 6:
			if calls == 0 {
				calls++
				return "abc123", nil
			}
			return "def456", nil
		default:
			t.Fatalf("hash length = %d, want 6", length)
			return "", nil
		}
	})
	defer restore()

	existsCalls := 0
	got, err := GenerateID("FID", func(candidate string) (bool, error) {
		existsCalls++
		return candidate == "FID-abc123", nil
	})
	if err != nil {
		t.Fatalf("GenerateID() error = %v", err)
	}

	if got != "FID-def456" {
		t.Fatalf("GenerateID() = %q, want %q", got, "FID-def456")
	}
	if existsCalls != 2 {
		t.Fatalf("exists predicate calls = %d, want 2", existsCalls)
	}
}

func TestFormatEventIDUsesUTCDateAndZeroPaddedSequence(t *testing.T) {
	occurredAt := time.Date(2026, time.June, 1, 1, 30, 0, 0, time.FixedZone("NZST", 12*60*60))

	got, err := FormatEventID(occurredAt, 42)
	if err != nil {
		t.Fatalf("FormatEventID() error = %v", err)
	}

	if got != "EVT-20260531-000042" {
		t.Fatalf("FormatEventID() = %q, want %q", got, "EVT-20260531-000042")
	}
}

func stubHashSuffixGenerator(stub func(int) (string, error)) func() {
	original := hashSuffixGenerator
	hashSuffixGenerator = stub
	return func() {
		hashSuffixGenerator = original
	}
}
