package cli

import (
	"fmt"
	"reflect"
	"strings"
)

// InputField describes one member of a closed JSON input object.
type InputField struct {
	Name          string
	Type          string
	Requirement   string // required, optional, or nullable
	Meaning       string
	Placeholder   string
	Format        string
	Members       []string
	MinInt        *int64
	MaxInt        *int64
	MinLen        *int
	MaxLen        *int
	Children      []InputField
	Element       *InputField
	Variants      map[string][]InputField // selected by a sibling discriminator
	Discriminator string                  // defaults to code for existing object variants
	VariantTypes  map[string]reflect.Type // actual element types of raw array variants
	Sensitive     bool
}

type InputDescriptor struct {
	Operation string
	Purpose   string
	Type      reflect.Type
	Fields    []InputField
}

// CheckInputType compares the public JSON tree, including optional and
// nullable members, with the request type. RawMessage is a closed variant
// whose members are selected by the descriptor's discriminator.
func CheckInputType(descriptor InputDescriptor) error {
	typeOf := descriptor.Type
	for typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}
	return checkInputFields(descriptor.Fields, typeOf, "input")
}

func checkInputFields(fields []InputField, typeOf reflect.Type, path string) error {
	if typeOf.Kind() != reflect.Struct {
		return fmt.Errorf("%s: Go type %s is not an object", path, typeOf)
	}
	seen := make(map[string]bool, len(fields))
	for _, field := range fields {
		if seen[field.Name] {
			return fmt.Errorf("%s.%s: duplicate descriptor field", path, field.Name)
		}
		seen[field.Name] = true
		var member reflect.StructField
		found := false
		for i := range typeOf.NumField() {
			candidate := typeOf.Field(i)
			if strings.Split(candidate.Tag.Get("json"), ",")[0] == field.Name {
				member, found = candidate, true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s.%s: no Go JSON field", path, field.Name)
		}
		goType := member.Type
		optional := strings.Contains(member.Tag.Get("json"), ",omitempty")
		if optional != (field.Requirement == "optional") {
			return fmt.Errorf("%s.%s: optional mismatch", path, field.Name)
		}
		if goType.Kind() == reflect.Pointer {
			if field.Requirement != "nullable" && field.Requirement != "optional" {
				return fmt.Errorf("%s.%s: nullable mismatch", path, field.Name)
			}
			goType = goType.Elem()
		}
		switch field.Type {
		case "object":
			if goType == reflect.TypeFor[[]byte]() || goType.String() == "json.RawMessage" {
				if len(field.Variants) == 0 {
					return fmt.Errorf("%s.%s: raw object needs variants", path, field.Name)
				}
			} else if err := checkInputFields(field.Children, goType, path+"."+field.Name); err != nil {
				return err
			}
		case "array":
			if goType.String() == "json.RawMessage" || goType == reflect.TypeFor[[]byte]() {
				if err := checkRawInputArray(field, path+"."+field.Name); err != nil {
					return err
				}
				continue
			}
			if goType.Kind() != reflect.Slice || field.Element == nil {
				return fmt.Errorf("%s.%s: array mismatch", path, field.Name)
			}
			if field.Element.Type == "object" {
				if err := checkInputFields(field.Element.Children, goType.Elem(), path+"."+field.Name+"[]"); err != nil {
					return err
				}
			} else if !inputScalarMatches(field.Element.Type, goType.Elem()) {
				return fmt.Errorf("%s.%s: element mismatch", path, field.Name)
			}
		default:
			if !inputScalarMatches(field.Type, goType) {
				return fmt.Errorf("%s.%s: type mismatch", path, field.Name)
			}
		}
	}
	count := 0
	for i := range typeOf.NumField() {
		field := typeOf.Field(i)
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if field.IsExported() && name != "-" {
			count++
		}
	}
	if count != len(seen) {
		return fmt.Errorf("%s: descriptor has %d fields, Go type has %d", path, len(seen), count)
	}
	return nil
}

func checkRawInputArray(field InputField, path string) error {
	if field.Element == nil || field.Element.Type != "object" || len(field.Variants) == 0 || len(field.VariantTypes) != len(field.Variants) {
		return fmt.Errorf("%s: raw array needs typed variants", path)
	}
	for name, fields := range field.Variants {
		typeOf := field.VariantTypes[name]
		if typeOf == nil {
			return fmt.Errorf("%s: variant %s has no Go type", path, name)
		}
		if err := checkInputFields(fields, typeOf, path+"[] ("+name+")"); err != nil {
			return err
		}
	}
	return nil
}

func inputDiscriminator(field InputField) string {
	if field.Discriminator != "" {
		return field.Discriminator
	}
	return "code"
}

func inputVariantChildren(field InputField, selector string) []InputField {
	if len(field.Variants) > 0 {
		return field.Variants[selector]
	}
	if field.Element != nil {
		return field.Element.Children
	}
	return field.Children
}

func inputScalarMatches(kind string, typeOf reflect.Type) bool {
	switch kind {
	case "string":
		return typeOf.Kind() == reflect.String
	case "integer":
		return typeOf.Kind() == reflect.Int || typeOf.Kind() == reflect.Int64
	case "boolean":
		return typeOf.Kind() == reflect.Bool
	}
	return false
}
