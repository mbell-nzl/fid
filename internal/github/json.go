package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// encoding/json accepts duplicate keys by keeping the last value. At this
// boundary that could hide a failure, so require one unambiguous JSON value.
func validateJSON(data []byte) error {
	if !json.Valid(data) {
		return fmt.Errorf("invalid JSON")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := jsonValue(decoder); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return fmt.Errorf("trailing JSON data")
	}
	return nil
}

func jsonValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("read JSON: %w", err)
	}
	delimiter, container := token.(json.Delim)
	if !container {
		return nil
	}
	keys := make(map[string]bool)
	for decoder.More() {
		if delimiter == '{' {
			token, err := decoder.Token()
			if err != nil {
				return fmt.Errorf("read JSON key: %w", err)
			}
			key, ok := token.(string)
			if !ok || keys[key] {
				return fmt.Errorf("invalid or duplicate JSON key")
			}
			keys[key] = true
		}
		if err := jsonValue(decoder); err != nil {
			return err
		}
	}
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("close JSON: %w", err)
	}
	return nil
}
