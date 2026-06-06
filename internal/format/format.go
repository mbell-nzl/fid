package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	JSON     = "json"
	TOON     = "toon"
	Markdown = "markdown"
	YAML     = "yaml"
)

type UnsupportedFormatError struct {
	Format string
}

func (err *UnsupportedFormatError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("format %q is not yet supported", err.Format)
}

func Normalize(value string) string {
	normalized := strings.TrimSpace(strings.ToLower(value))
	if normalized == "" {
		return JSON
	}
	return normalized
}

func Write(writer io.Writer, formatName string, value any) error {
	if writer == nil {
		return fmt.Errorf("write format output: writer is nil")
	}

	switch normalized := Normalize(formatName); normalized {
	case JSON:
		encoder := json.NewEncoder(writer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(value); err != nil {
			return fmt.Errorf("encode json output: %w", err)
		}
		return nil
	case TOON, Markdown, YAML:
		return &UnsupportedFormatError{Format: normalized}
	default:
		return &UnsupportedFormatError{Format: normalized}
	}
}
