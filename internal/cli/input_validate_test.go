package cli

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func rawArrayDescriptor() InputDescriptor {
	minimum := 1
	return InputDescriptor{Type: reflect.TypeFor[struct {
		ComponentType string          `json:"component_type"`
		Content       json.RawMessage `json:"content"`
	}](), Fields: []InputField{
		{Name: "component_type", Type: "string", Requirement: "required", Members: []string{"authority-profile", "gate-definitions"}},
		{Name: "content", Type: "array", Requirement: "required", MinLen: &minimum, Discriminator: "component_type", VariantTypes: map[string]reflect.Type{
			"authority-profile": reflect.TypeFor[struct {
				ActorID string `json:"actor_id"`
			}](),
			"gate-definitions": reflect.TypeFor[struct {
				DefinitionID string `json:"definition_id"`
				Secret       string `json:"secret"`
			}](),
		}, Element: &InputField{Type: "object", Requirement: "required"}, Variants: map[string][]InputField{
			"authority-profile": {{Name: "actor_id", Type: "string", Requirement: "required"}},
			"gate-definitions":  {{Name: "definition_id", Type: "string", Requirement: "required"}, {Name: "secret", Type: "string", Requirement: "required", Sensitive: true}},
		}},
	}}
}

func TestInputRawArrayVariants(t *testing.T) {
	t.Parallel()
	descriptor := rawArrayDescriptor()
	t.Run("reflection", func(t *testing.T) {
		t.Parallel()
		if err := CheckInputType(descriptor); err != nil {
			t.Fatal(err)
		}
	})
	for _, test := range []struct {
		payload, path string
		value         any
	}{
		{`{"component_type":"authority-profile","content":[{"actor_id":"A"}]}`, "", nil},
		{`{"component_type":"gate-definitions","content":[{"definition_id":"G","secret":"sentinel"}]}`, "", nil},
		{`{"component_type":"authority-profile","content":[{"actor_id":null}]}`, "input.content[0].actor_id", "null"},
		{`{"component_type":"gate-definitions","content":[{"definition_id":"G","secret":null}]}`, "input.content[0].secret", nil},
		{`{"component_type":"gate-definitions","content":[{"definition_id":"G","secret":"sentinel","secret":"sentinel"}]}`, "input.content[0].secret", nil},
	} {
		t.Run(test.payload, func(t *testing.T) {
			t.Parallel()
			err := ValidateInput([]byte(test.payload), descriptor)
			if test.path == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			invalid, ok := err.(*InputError)
			if !ok || invalid.Path != test.path || !reflect.DeepEqual(invalid.Value, test.value) {
				t.Fatalf("error = %#v; want %s value=%#v", err, test.path, test.value)
			}
		})
	}
	t.Run("help and example", func(t *testing.T) {
		t.Parallel()
		help := InputHelp(descriptor)
		for _, text := range []string{"when component_type=authority-profile:", "when component_type=gate-definitions:", "input.content[i].actor_id", "input.content[i].secret"} {
			if !strings.Contains(help, text) {
				t.Errorf("help lacks %s: %s", text, help)
			}
		}
		example, err := InputExample(descriptor)
		if err != nil || string(example) != "{\n  \"component_type\": \"authority-profile\",\n  \"content\": [{\n    \"actor_id\": \"<actor-id>\"\n  }]\n}\n" {
			t.Fatalf("example = %s / %v", example, err)
		}
	})
}

func TestInputValidationOrderAndSensitiveValues(t *testing.T) {
	t.Parallel()
	min := int64(1)
	descriptor := InputDescriptor{Fields: []InputField{
		{Name: "name", Type: "string", Requirement: "required", Format: "ID"},
		{Name: "version", Type: "integer", Requirement: "required", MinInt: &min},
		{Name: "secret", Type: "string", Requirement: "required", Sensitive: true, Members: []string{"one"}},
	}}
	for _, test := range []struct {
		name, payload, path string
		value               any
		allowed             []string
	}{
		{"unknown lexical", `{"name":"N","version":1,"secret":"one","z":1,"a":1}`, "input.a", "a", []string{"name", "secret", "version"}},
		{"duplicate before unknown", `{"name":"N","name":"X","version":1,"secret":"one","a":1}`, "input.name", "name", []string{"name", "secret", "version"}},
		{"member order", `{"version":"bad","secret":"bad"}`, "input.name", nil, []string{"name"}},
		{"null", `{"name":null,"version":1,"secret":"one"}`, "input.name", "null", []string{"string"}},
		{"wrong type", `{"name":3,"version":1,"secret":"one"}`, "input.name", "integer", []string{"string"}},
		{"integer bound", `{"name":"N","version":0,"secret":"one"}`, "input.version", float64(0), []string{"integer range 1..unbounded"}},
		{"sensitive enum", `{"name":"N","version":1,"secret":"PRIVATE"}`, "input.secret", nil, []string{"one"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateInput([]byte(test.payload), descriptor)
			invalid, ok := err.(*InputError)
			if !ok || invalid.Path != test.path || !reflect.DeepEqual(invalid.Value, test.value) || !reflect.DeepEqual(invalid.Allowed, test.allowed) {
				t.Fatalf("error = %#v; want path=%s value=%#v allowed=%v", err, test.path, test.value, test.allowed)
			}
		})
	}
}

func TestInputArrayHelpElementBounds(t *testing.T) {
	t.Parallel()
	minimum, maximum := int64(-255), int64(255)
	descriptor := InputDescriptor{Fields: []InputField{{Name: "codes", Type: "array", Requirement: "required", Element: &InputField{Type: "integer", MinInt: &minimum, MaxInt: &maximum}}}}
	if help := InputHelp(descriptor); !strings.Contains(help, "element integer range -255..255") {
		t.Fatalf("missing element bounds: %s", help)
	}
}

func TestInputRawArrayReflectionRejectsVariantDrift(t *testing.T) {
	t.Parallel()
	for _, variant := range []string{"authority-profile", "gate-definitions"} {
		t.Run(variant, func(t *testing.T) {
			t.Parallel()
			descriptor := rawArrayDescriptor()
			content := &descriptor.Fields[1]
			content.Variants[variant][0].Name = "drifted"
			if err := CheckInputType(descriptor); err == nil || !strings.Contains(err.Error(), variant) {
				t.Fatalf("variant drift accepted: %v", err)
			}
		})
	}
	descriptor := rawArrayDescriptor()
	descriptor.Fields[1].VariantTypes = nil
	if CheckInputType(descriptor) == nil {
		t.Fatal("untyped raw array exemption")
	}
}

func TestInputCodeObjectVariantsRemainSelected(t *testing.T) {
	t.Parallel()
	descriptor := InputDescriptor{Fields: []InputField{
		{Name: "code", Type: "string", Requirement: "required", Members: []string{"one", "two"}},
		{Name: "details", Type: "object", Requirement: "required", Variants: map[string][]InputField{
			"one": {{Name: "id", Type: "string", Requirement: "required"}},
			"two": {{Name: "secret", Type: "string", Requirement: "required", Sensitive: true}},
		}},
	}}
	for _, payload := range []string{`{"code":"one","details":{"id":"I"}}`, `{"code":"two","details":{"secret":"literal"}}`} {
		if err := ValidateInput([]byte(payload), descriptor); err != nil {
			t.Fatal(err)
		}
	}
	err := ValidateInput([]byte(`{"code":"two","details":{"secret":"literal","secret":"literal"}}`), descriptor)
	invalid, ok := err.(*InputError)
	if !ok || invalid.Path != "input.details.secret" || invalid.Value != nil {
		t.Fatalf("duplicate secret: %#v", err)
	}
	if !strings.Contains(InputHelp(descriptor), "when code=two:") {
		t.Fatal("code help lost")
	}
	example, err := InputExample(descriptor)
	if err != nil || string(example) != "{\n  \"code\": \"one\",\n  \"details\": {\n    \"id\": \"<id>\"\n  }\n}\n" {
		t.Fatalf("code example=%s / %v", example, err)
	}
}

func TestInputRawArrayScalarBoundsPreserveOrdinaryCallers(t *testing.T) {
	t.Parallel()
	descriptor := rawArrayDescriptor()
	minimum, maximum := 1, 1
	field := &descriptor.Fields[1].Variants["authority-profile"][0]
	field.MinLen, field.MaxLen = &minimum, &maximum
	if err := ValidateInput([]byte(`{"component_type":"authority-profile","content":[{"actor_id":"é"}]}`), descriptor); err != nil {
		t.Fatalf("raw scalar bound: %v", err)
	}
	if err := ValidateInput([]byte(`{"component_type":"authority-profile","content":[{"actor_id":"éé"}]}`), descriptor); err == nil {
		t.Fatal("raw scalar bound weakened")
	}
	ordinary := InputDescriptor{Fields: []InputField{*field}}
	if err := ValidateInput([]byte(`{"actor_id":"é"}`), ordinary); err == nil {
		t.Fatal("ordinary string length behavior changed")
	}
}
