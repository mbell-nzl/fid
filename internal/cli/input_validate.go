package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// InputError is a structural invalid-argument with a public field path.
type InputError struct {
	Path    string
	Value   any
	Allowed []string
}

func (e *InputError) Error() string { return "invalid " + e.Path }

func (e *InputError) Details() map[string]any {
	value := e.Value
	if value != nil {
		if _, ok := value.(string); !ok {
			encoded, err := json.Marshal(value)
			if err != nil {
				value = fmt.Sprint(value)
			} else {
				value = string(encoded)
			}
		}
	}
	return map[string]any{"argument": e.Path, "value": value, "allowed": e.Allowed}
}

// ValidateInput checks a single closed object before decoding or domain work.
func ValidateInput(payload []byte, descriptor InputDescriptor) error {
	if len(payload) > 1_048_576 || !utf8.Valid(payload) || !json.Valid(payload) {
		return &InputError{Path: "input", Allowed: []string{"single UTF-8 JSON object"}}
	}
	var root json.RawMessage
	if err := json.Unmarshal(payload, &root); err != nil || jsonKind(root) != "object" {
		return &InputError{Path: "input", Value: jsonKind(root), Allowed: []string{"object"}}
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal(payload, &values); err != nil {
		return &InputError{Path: "input", Allowed: []string{"object"}}
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	if err := scanInputDuplicates(decoder, "input", descriptor.Fields, values); err != nil {
		return err
	}
	return validateInputObject(values, descriptor.Fields, "input", values)
}

func scanInputDuplicates(decoder *json.Decoder, location string, fields []InputField, root map[string]json.RawMessage) error {
	token, err := decoder.Token()
	if err != nil {
		return &InputError{Path: location, Allowed: []string{"valid JSON"}}
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	if delim == '[' {
		for i := 0; decoder.More(); i++ {
			if err := scanInputDuplicates(decoder, fmt.Sprintf("%s[%d]", location, i), fields, root); err != nil {
				return err
			}
		}
		_, err = decoder.Token()
		return err
	}
	if delim != '{' {
		return &InputError{Path: location, Allowed: []string{"object"}}
	}
	seen := map[string]bool{}
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return err
		}
		key := keyToken.(string)
		field := findInputField(fields, key)
		if seen[key] {
			value := any(key)
			if field != nil && field.Sensitive {
				value = nil
			}
			return &InputError{Path: location + "." + key, Value: value, Allowed: inputNames(fields)}
		}
		seen[key] = true
		children := []InputField(nil)
		if field != nil {
			children = selectedInputChildren(*field, root)
		}
		if err := scanInputDuplicates(decoder, location+"."+key, children, root); err != nil {
			return err
		}
	}
	_, err = decoder.Token()
	return err
}

func validateInputObject(values map[string]json.RawMessage, fields []InputField, location string, root map[string]json.RawMessage) error {
	keys := make([]string, 0, len(values))
	for key := range values {
		if findInputField(fields, key) == nil {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	if len(keys) > 0 {
		return &InputError{Path: location + "." + keys[0], Value: keys[0], Allowed: inputNames(fields)}
	}
	for _, field := range fields {
		location := location + "." + field.Name
		value, exists := values[field.Name]
		if !exists {
			if field.Requirement == "optional" {
				continue
			}
			return &InputError{Path: location, Allowed: []string{field.Name}}
		}
		if err := validateInputField(value, field, location, root); err != nil {
			return err
		}
	}
	return nil
}

func validateInputField(value json.RawMessage, field InputField, location string, root map[string]json.RawMessage) error {
	kind := jsonKind(value)
	if kind == "null" && field.Requirement == "nullable" {
		return nil
	}
	if kind != field.Type {
		return inputValueError(field, location, kind, []string{field.Type})
	}
	var item any
	_ = json.Unmarshal(value, &item)
	if len(field.Members) > 0 {
		member := fmt.Sprint(item)
		if !containsInputMember(field.Members, member) {
			return inputValueError(field, location, item, field.Members)
		}
	}
	if field.Format != "" && field.Type != "array" && !matchesInputFormat(field.Format, item) {
		return inputValueError(field, location, item, []string{field.Format})
	}
	if field.Type == "string" {
		length := len([]byte(item.(string)))
		if field.scalarBounds {
			length = utf8.RuneCountInString(item.(string))
		}
		if field.MinLen != nil && length < *field.MinLen || field.MaxLen != nil && length > *field.MaxLen {
			return inputValueError(field, location, item, []string{inputBounds(field)})
		}
	}
	if field.Type == "integer" {
		var number int64
		if err := json.Unmarshal(value, &number); err != nil {
			return inputValueError(field, location, item, []string{"integer"})
		}
		if field.MinInt != nil && number < *field.MinInt || field.MaxInt != nil && number > *field.MaxInt {
			return inputValueError(field, location, item, []string{inputBounds(field)})
		}
	}
	if field.Type == "object" {
		var members map[string]json.RawMessage
		_ = json.Unmarshal(value, &members)
		children := selectedInputChildren(field, root)
		return validateInputObject(members, children, location, root)
	}
	if field.Type == "array" {
		var entries []json.RawMessage
		_ = json.Unmarshal(value, &entries)
		if field.MinLen != nil && len(entries) < *field.MinLen || field.MaxLen != nil && len(entries) > *field.MaxLen {
			return inputValueError(field, location, entries, []string{inputBounds(field)})
		}
		for i, entry := range entries {
			element := *field.Element
			element.Sensitive = element.Sensitive || field.Sensitive
			entryPath := fmt.Sprintf("%s[%d]", location, i)
			if len(field.Variants) > 0 {
				element.Children = selectedInputChildren(field, root)
				if inputVariantMismatch(entry, field, element.Children) {
					return &InputError{Path: entryPath, Allowed: inputNames(element.Children)}
				}
			}
			if err := validateInputField(entry, element, entryPath, root); err != nil {
				return err
			}
		}
		if field.Format == "sorted-unique" && !inputSortedUnique(entries) {
			return inputValueError(field, location, entries, []string{"sorted-unique"})
		}
	}
	return nil
}

func inputValueError(field InputField, location string, value any, allowed []string) error {
	if inputHasSensitive(field) {
		value = nil
	}
	return &InputError{Path: location, Value: value, Allowed: allowed}
}

func selectedInputChildren(field InputField, root map[string]json.RawMessage) []InputField {
	var selector string
	_ = json.Unmarshal(root[inputDiscriminator(field)], &selector)
	return inputVariantChildren(field, selector)
}

func inputVariantMismatch(entry json.RawMessage, field InputField, selected []InputField) bool {
	var members map[string]json.RawMessage
	if json.Unmarshal(entry, &members) != nil {
		return false
	}
	for _, variant := range field.Variants {
		for _, candidate := range variant {
			if _, exists := members[candidate.Name]; exists && findInputField(selected, candidate.Name) == nil {
				return true
			}
		}
	}
	return false
}

func inputHasSensitive(field InputField) bool {
	if field.Sensitive {
		return true
	}
	if field.Element != nil && inputHasSensitive(*field.Element) {
		return true
	}
	children := append([]InputField(nil), field.Children...)
	for _, variant := range field.Variants {
		children = append(children, variant...)
	}
	for _, child := range children {
		if inputHasSensitive(child) {
			return true
		}
	}
	return false
}

func jsonKind(value []byte) string {
	if len(value) == 0 {
		return "unknown"
	}
	switch value[0] {
	case '{':
		return "object"
	case '[':
		return "array"
	case '"':
		return "string"
	case 't', 'f':
		return "boolean"
	case 'n':
		return "null"
	default:
		if bytes.ContainsAny(value, ".eE") {
			return "number"
		}
		return "integer"
	}
}

func inputNames(fields []InputField) []string {
	names := make([]string, len(fields))
	for i := range fields {
		names[i] = fields[i].Name
	}
	sort.Strings(names)
	return names
}

func findInputField(fields []InputField, name string) *InputField {
	for i := range fields {
		if fields[i].Name == name {
			return &fields[i]
		}
	}
	return nil
}

func containsInputMember(members []string, value string) bool {
	for _, member := range members {
		if member == value {
			return true
		}
	}
	return false
}

func inputBounds(field InputField) string {
	if field.MinInt != nil || field.MaxInt != nil {
		return "integer range " + inputIntBound(field.MinInt, "unbounded") + ".." + inputIntBound(field.MaxInt, "unbounded")
	}
	return "length " + inputLengthBound(field.MinLen, "0") + ".." + inputLengthBound(field.MaxLen, "unbounded")
}

func inputIntBound(value *int64, fallback string) string {
	if value == nil {
		return fallback
	}
	return fmt.Sprint(*value)
}

func inputLengthBound(value *int, fallback string) string {
	if value == nil {
		return fallback
	}
	return fmt.Sprint(*value)
}

var inputHex = regexp.MustCompile(`^[0-9a-f]+$`)

func matchesInputFormat(format string, value any) bool {
	text, ok := value.(string)
	if !ok {
		return false
	}
	switch format {
	case "ID":
		if len(text) == 0 || len(text) > 128 {
			return false
		}
		for _, r := range text {
			if r < 0x20 || r > 0x7e {
				return false
			}
		}
		return true
	case "text":
		return len(text) <= 4096
	case "digest":
		return strings.HasPrefix(text, "sha256:") && len(text) == 71 && inputHex.MatchString(text[7:])
	case "OID", "oid":
		return len(text) == 40 && inputHex.MatchString(text)
	case "relative path":
		return text != "" && !strings.HasPrefix(text, "/") && path.Clean(text) == text && !strings.HasPrefix(text, "..") && !strings.ContainsRune(text, 0)
	case "timestamp":
		_, err := time.Parse(time.RFC3339Nano, text)
		return err == nil
	case "environment name":
		return regexp.MustCompile(`^[A-Z_][A-Z0-9_]{0,127}$`).MatchString(text)
	case "cwd path":
		return text == "." || componentInputPath(text)
	case "artifact path":
		return componentInputPath(text)
	}
	return false
}

func componentInputPath(text string) bool {
	return text != "" && text != "." && len(text) <= 512 && !strings.HasPrefix(text, "/") &&
		!strings.Contains(text, "\\") && path.Clean(text) == text && !strings.HasPrefix(text, "../") && text != ".." && !strings.ContainsRune(text, 0)
}

func inputSortedUnique(entries []json.RawMessage) bool {
	last := ""
	var lastNumber int64
	for i, entry := range entries {
		if jsonKind(entry) == "integer" {
			var number int64
			if json.Unmarshal(entry, &number) != nil || i > 0 && number <= lastNumber {
				return false
			}
			lastNumber = number
			continue
		}
		var item string
		if err := json.Unmarshal(entry, &item); err != nil || item <= last {
			return false
		}
		last = item
	}
	return true
}
