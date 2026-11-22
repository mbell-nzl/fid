package journal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

func validateCurrentJSONText(raw string) error {
	return validateTypedJSONText(raw, nil)
}

// Use the same wire type at admission and stored-Event decoding. Map entries
// stay case-sensitive; only objects decoded into structs alias field names.
func validateEventJSONText(raw string, schema int) error {
	target := reflect.TypeFor[eventV2Wire]()
	if schema == EventSchemaV1 {
		target = reflect.TypeFor[eventV1Wire]()
	}
	return validateTypedJSONText(raw, target)
}

func validateTypedJSONText(raw string, target reflect.Type) error {
	if !utf8.ValidString(raw) {
		return fmt.Errorf("invalid UTF-8 JSON")
	}
	decoder := json.NewDecoder(bytes.NewBufferString(raw))
	decoder.UseNumber()
	if err := scanCurrentJSONValue(decoder, target); err != nil {
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	return nil
}

func decodeCurrentJSON(raw string) (any, error) {
	if err := validateCurrentJSONText(raw); err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewBufferString(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	return value, nil
}

// Stored retry readers allow unknown fields, but use the same key rules as
// other embedded JSON: struct aliases are refused and map keys match exactly.
func unmarshalCurrentJSON(raw string, target any) error {
	if err := validateTypedJSONText(raw, reflect.TypeOf(target)); err != nil {
		return err
	}
	return json.Unmarshal([]byte(raw), target)
}

func scanCurrentJSONValue(decoder *json.Decoder, target reflect.Type) error {
	for target != nil && target.Kind() == reflect.Pointer {
		target = target.Elem()
	}
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		// Keep Token's original integer precision, but admit the same finite
		// number domain as the archive's encoding/json value check.
		if number, numeric := token.(json.Number); numeric {
			var finite float64
			if err := json.Unmarshal([]byte(number), &finite); err != nil {
				return err
			}
		}
		return nil
	}
	switch delimiter {
	case '{':
		keys := map[string]bool{}
		for decoder.More() {
			keyValue, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyValue.(string)
			if !ok {
				return fmt.Errorf("object key is not text")
			}
			name := key
			if target != nil && target.Kind() == reflect.Struct {
				name = foldJSONKey(key)
			}
			if keys[name] {
				return fmt.Errorf("duplicate key %q", key)
			}
			keys[name] = true
			if err := scanCurrentJSONValue(decoder, currentJSONFieldType(target, key)); err != nil {
				return err
			}
		}
		_, err := decoder.Token()
		return err
	case '[':
		var element reflect.Type
		if target != nil && (target.Kind() == reflect.Slice || target.Kind() == reflect.Array) {
			element = target.Elem()
		}
		for decoder.More() {
			if err := scanCurrentJSONValue(decoder, element); err != nil {
				return err
			}
		}
		_, err := decoder.Token()
		return err
	default:
		return fmt.Errorf("unexpected JSON delimiter %q", delimiter)
	}
}

func currentJSONFieldType(target reflect.Type, key string) reflect.Type {
	if target == nil {
		return nil
	}
	if target.Kind() == reflect.Map {
		return target.Elem()
	}
	if target.Kind() != reflect.Struct {
		return nil
	}
	for i := range target.NumField() {
		field := target.Field(i)
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" || !field.IsExported() {
			continue
		}
		if name == "" && field.Anonymous {
			embedded := field.Type
			if embedded.Kind() == reflect.Pointer {
				embedded = embedded.Elem()
			}
			if found := currentJSONFieldType(embedded, key); found != nil {
				return found
			}
			continue
		}
		if name == "" {
			name = field.Name
		}
		if foldJSONKey(name) == foldJSONKey(key) {
			return field.Type
		}
	}
	return nil
}

func foldJSONKey(key string) string {
	return strings.Map(func(r rune) rune {
		folded := r
		for next := unicode.SimpleFold(r); next != r; next = unicode.SimpleFold(next) {
			folded = min(folded, next)
		}
		return folded
	}, key)
}
