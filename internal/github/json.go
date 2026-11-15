package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// encoding/json accepts duplicate keys and matches struct fields without regard
// to case. Either can hide a failure, so reject duplicate names under the same
// Unicode case folding before decoding into structs.
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
			key = foldedJSONKey(key)
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

func foldedJSONKey(key string) string {
	return strings.Map(func(r rune) rune {
		for {
			next := unicode.SimpleFold(r)
			if next <= r {
				return next
			}
			r = next
		}
	}, key)
}
