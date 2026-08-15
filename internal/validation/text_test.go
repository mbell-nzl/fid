package validation

import (
	"strings"
	"testing"
)

func TestNormalizeTextAppliesClosedByteAndCharacterDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		maximum int
		want    string
		wantErr string
	}{
		{name: "trims Unicode whitespace", value: " \tvalue\n", maximum: 5, want: "value"},
		{name: "accepts exact byte maximum", value: strings.Repeat("x", 5), maximum: 5, want: "xxxxx"},
		{name: "rejects empty", value: " \t\n", maximum: 5, wantErr: "field is required"},
		{name: "rejects over byte maximum", value: strings.Repeat("x", 6), maximum: 5, wantErr: "field exceeds 5 bytes"},
		{name: "rejects invalid UTF-8", value: string([]byte{0xff}), maximum: 5, wantErr: "field must be valid UTF-8"},
		{name: "rejects NUL", value: "x\x00y", maximum: 5, wantErr: "field contains ASCII control character"},
		{name: "rejects unit separator", value: "x\x1fy", maximum: 5, wantErr: "field contains ASCII control character"},
		{name: "rejects DEL", value: "x\x7fy", maximum: 5, wantErr: "field contains ASCII control character"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NormalizeText("field", test.value, test.maximum)
			if test.wantErr != "" {
				if err == nil || err.Error() != test.wantErr {
					t.Fatalf("NormalizeText() error = %v, want %q", err, test.wantErr)
				}
				return
			}
			if err != nil || got != test.want {
				t.Fatalf("NormalizeText() = (%q, %v), want (%q, nil)", got, err, test.want)
			}
		})
	}
}
