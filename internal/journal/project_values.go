package journal

import (
	"fmt"
	"math"
	"sort"
	"time"
)

func validatedIdentifier(identifier string) (string, error) {
	if !identifierPattern.MatchString(identifier) {
		return "", fmt.Errorf("invalid SQL identifier %q", identifier)
	}
	return identifier, nil
}

func sortedMapKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func cloneMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return map[string]any{}
	}

	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

func valuesEqual(expected, actual any) bool {
	expected = normalizeValue(expected)
	actual = normalizeValue(actual)

	expectedNumber, expectedIsNumber := asFloat(expected)
	actualNumber, actualIsNumber := asFloat(actual)
	if expectedIsNumber && actualIsNumber {
		return math.Abs(expectedNumber-actualNumber) < 1e-9
	}

	return fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", actual)
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case []byte:
		return string(typed)
	case bool:
		if typed {
			return int64(1)
		}
		return int64(0)
	case time.Time:
		if typed.Hour() == 0 && typed.Minute() == 0 && typed.Second() == 0 && typed.Nanosecond() == 0 {
			return typed.UTC().Format("2006-01-02")
		}
		return typed.UTC().Format(time.RFC3339)
	default:
		return value
	}
}

func asFloat(value any) (float64, bool) {
	switch typed := value.(type) {
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case float32:
		return float64(typed), true
	case float64:
		return typed, true
	default:
		return 0, false
	}
}

func replaceValue(values []string, oldValue, newValue string) []string {
	replaced := make([]string, len(values))
	for index, value := range values {
		if value == oldValue {
			replaced[index] = newValue
			continue
		}
		replaced[index] = value
	}

	return replaced
}
