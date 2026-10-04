package cli

import (
	"reflect"
	"testing"
)

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
