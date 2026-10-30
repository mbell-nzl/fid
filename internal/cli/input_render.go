package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func InputHelp(descriptor InputDescriptor) string {
	var output strings.Builder
	output.WriteString(descriptor.Purpose + "\n\nInput fields:\n")
	writeInputHelpFields(&output, descriptor.Fields, "input")
	output.WriteString("\nTemplate: fid " + strings.ReplaceAll(descriptor.Operation, ".", " ") + " --example\n")
	return output.String()
}

func writeInputHelpFields(output *strings.Builder, fields []InputField, prefix string) {
	for _, field := range fields {
		typeName := field.Type
		if field.Type == "array" && field.Element != nil {
			typeName += "<" + field.Element.Type + ">"
		}
		_, _ = fmt.Fprintf(output, "  %s.%s: %s, %s", prefix, field.Name, typeName, field.Requirement)
		if field.Format != "" {
			_, _ = fmt.Fprintf(output, ", %s", field.Format)
		}
		if len(field.Members) > 0 {
			_, _ = fmt.Fprintf(output, ", allowed: %s", strings.Join(field.Members, "|"))
		}
		if field.Element != nil && len(field.Element.Members) > 0 {
			_, _ = fmt.Fprintf(output, ", elements allowed: %s", strings.Join(field.Element.Members, "|"))
		}
		if field.Element != nil && field.Element.Format != "" && len(field.Element.Members) == 0 {
			_, _ = fmt.Fprintf(output, ", element format: %s", field.Element.Format)
		}
		if field.Element != nil && (field.Element.MinInt != nil || field.Element.MaxInt != nil || field.Element.MinLen != nil || field.Element.MaxLen != nil) {
			_, _ = fmt.Fprintf(output, ", element %s", inputBounds(*field.Element))
		}
		if field.MinInt != nil || field.MaxInt != nil || field.MinLen != nil || field.MaxLen != nil {
			_, _ = fmt.Fprintf(output, ", %s", inputBounds(field))
		}
		_, _ = fmt.Fprintf(output, " — %s\n", field.Meaning)
		if field.Type == "object" {
			writeInputHelpFields(output, field.Children, prefix+"."+field.Name)
		}
		if len(field.Variants) > 0 {
			keys := make([]string, 0, len(field.Variants))
			for key := range field.Variants {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				_, _ = fmt.Fprintf(output, "    when %s=%s:\n", inputDiscriminator(field), key)
				variantPath := prefix + "." + field.Name
				if field.Type == "array" {
					variantPath += "[i]"
				}
				writeInputHelpFields(output, field.Variants[key], variantPath)
			}
		}
		if field.Type == "array" && field.Element != nil && field.Element.Type == "object" {
			writeInputHelpFields(output, field.Element.Children, prefix+"."+field.Name+"[i]")
		}
	}
}

func InputExample(descriptor InputDescriptor) ([]byte, error) {
	var output bytes.Buffer
	if err := writeInputExampleObject(&output, descriptor.Fields, 0, descriptor); err != nil {
		return nil, err
	}
	output.WriteByte('\n')
	return output.Bytes(), nil
}

func writeInputExampleObject(output *bytes.Buffer, fields []InputField, depth int, descriptor InputDescriptor) error {
	output.WriteByte('{')
	first := true
	for _, field := range fields {
		if field.Requirement == "optional" {
			continue
		}
		if first {
			output.WriteByte('\n')
		} else {
			output.WriteString(",\n")
		}
		first = false
		output.WriteString(strings.Repeat("  ", depth+1))
		name, _ := json.Marshal(field.Name)
		output.Write(name)
		output.WriteString(": ")
		if err := writeInputExampleValue(output, field, depth+1, descriptor); err != nil {
			return err
		}
	}
	if !first {
		output.WriteByte('\n')
		output.WriteString(strings.Repeat("  ", depth))
	}
	output.WriteByte('}')
	return nil
}

func writeInputExampleValue(output *bytes.Buffer, field InputField, depth int, descriptor InputDescriptor) error {
	if field.Requirement == "nullable" {
		output.WriteString("null")
		return nil
	}
	switch field.Type {
	case "object":
		children := exampleInputChildren(field, descriptor)
		return writeInputExampleObject(output, children, depth, descriptor)
	case "array":
		output.WriteByte('[')
		if field.MinLen != nil && *field.MinLen > 0 {
			element := *field.Element
			if len(field.Variants) > 0 {
				element.Children = exampleInputChildren(field, descriptor)
			}
			if err := writeInputExampleValue(output, element, depth, descriptor); err != nil {
				return err
			}
		}
		output.WriteByte(']')
	case "string":
		value := field.Placeholder
		if len(field.Members) > 0 {
			value = field.Members[0]
		} else if value == "" {
			value = "<" + strings.ReplaceAll(field.Name, "_", "-") + ">"
		}
		output.WriteString(strconv.Quote(value))
	case "integer":
		value := int64(0)
		if field.MinInt != nil {
			value = *field.MinInt
		}
		_, _ = fmt.Fprint(output, value)
	case "boolean":
		output.WriteString("false")
	default:
		return fmt.Errorf("unsupported input field %s type %s", field.Name, field.Type)
	}
	return nil
}

func exampleInputChildren(field InputField, descriptor InputDescriptor) []InputField {
	selector := ""
	for _, candidate := range descriptor.Fields {
		if candidate.Name == inputDiscriminator(field) && len(candidate.Members) > 0 {
			selector = candidate.Members[0]
		}
	}
	return inputVariantChildren(field, selector)
}
