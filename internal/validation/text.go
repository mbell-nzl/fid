package validation

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func NormalizeText(field, value string, maximumBytes int) (string, error) {
	if !utf8.ValidString(value) {
		return "", fmt.Errorf("%s must be valid UTF-8", field)
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s is required", field)
	}
	if len(trimmed) > maximumBytes {
		return "", fmt.Errorf("%s exceeds %d bytes", field, maximumBytes)
	}
	for _, character := range []byte(trimmed) {
		if character < 0x20 || character == 0x7f {
			return "", fmt.Errorf("%s contains ASCII control character", field)
		}
	}
	return trimmed, nil
}
