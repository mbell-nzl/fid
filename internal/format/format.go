package format

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
)

const (
	JSON     = "json"
	TOON     = "toon"
	Markdown = "markdown"
	YAML     = "yaml"
)

type UnsupportedFormatError struct {
	Format          string
	KnownFormats    []string
	NotYetSupported bool
}

func (err *UnsupportedFormatError) Error() string {
	if err == nil {
		return "<nil>"
	}
	if err.NotYetSupported {
		return fmt.Sprintf("format %q is not yet supported", err.Format)
	}
	if len(err.KnownFormats) == 0 {
		return fmt.Sprintf("unsupported format %q", err.Format)
	}
	return fmt.Sprintf("unsupported format %q; known formats: %s", err.Format, strings.Join(err.KnownFormats, ", "))
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
	case TOON:
		if err := writeTOON(writer, value); err != nil {
			return fmt.Errorf("encode toon output: %w", err)
		}
		return nil
	case Markdown:
		if err := writeMarkdown(writer, value); err != nil {
			return fmt.Errorf("encode markdown output: %w", err)
		}
		return nil
	case YAML:
		if err := writeYAML(writer, value); err != nil {
			return fmt.Errorf("encode yaml output: %w", err)
		}
		return nil
	default:
		return &UnsupportedFormatError{
			Format:       normalized,
			KnownFormats: knownFormats(),
		}
	}
}

func genericStructure(value any) (any, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal generic structure: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()

	var result any
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("decode generic structure: %w", err)
	}
	return result, nil
}

func knownFormats() []string {
	formats := []string{JSON, TOON, Markdown, YAML}
	return slices.Clone(formats)
}
