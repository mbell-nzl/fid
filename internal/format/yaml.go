package format

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

func writeYAML(writer io.Writer, value any) error {
	structure, err := genericStructure(value)
	if err != nil {
		return err
	}

	encoded, err := yaml.Marshal(yamlReadyValue(structure))
	if err != nil {
		return fmt.Errorf("marshal yaml output: %w", err)
	}

	if _, err := writer.Write(encoded); err != nil {
		return fmt.Errorf("write yaml output: %w", err)
	}
	return nil
}

func yamlReadyValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, nested := range typed {
			result[key] = yamlReadyValue(nested)
		}
		return result
	case []any:
		result := make([]any, 0, len(typed))
		for _, nested := range typed {
			result = append(result, yamlReadyValue(nested))
		}
		return result
	case json.Number:
		return yamlScalarFromNumber(typed)
	default:
		return value
	}
}

func yamlScalarFromNumber(number json.Number) any {
	if intValue, err := number.Int64(); err == nil {
		return intValue
	}
	if floatValue, err := number.Float64(); err == nil {
		return floatValue
	}
	return number.String()
}
