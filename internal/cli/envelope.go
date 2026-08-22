// Package cli provides the domain-neutral fid-cli.v1 response transport.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
)

const ContractVersion = "fid-cli.v1"

type Event struct {
	ID         string `json:"id"`
	OccurredAt string `json:"occurred_at"`
}

type successEnvelope struct {
	ContractVersion string         `json:"contract_version"`
	OK              bool           `json:"ok"`
	Operation       string         `json:"operation"`
	Result          map[string]any `json:"result"`
	Event           *Event         `json:"event"`
}

type failureEnvelope struct {
	ContractVersion string  `json:"contract_version"`
	OK              bool    `json:"ok"`
	Operation       string  `json:"operation"`
	Error           Failure `json:"error"`
}

func WriteSuccess(writer io.Writer, operation string, result map[string]any, event *Event) error {
	if writer == nil {
		return fmt.Errorf("write success envelope: writer is nil")
	}
	if result == nil {
		return fmt.Errorf("write success envelope: result is nil")
	}

	envelope := successEnvelope{
		ContractVersion: ContractVersion,
		OK:              true,
		Operation:       operation,
		Result:          result,
		Event:           event,
	}
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(envelope); err != nil {
		return fmt.Errorf("write success envelope: %w", err)
	}
	return nil
}

func WriteFailure(writer io.Writer, operation string, failure Failure) error {
	if writer == nil {
		return fmt.Errorf("write failure envelope: writer is nil")
	}
	definition, ok := reasonCodeDefinitions[failure.Code]
	if !ok {
		return fmt.Errorf("write failure envelope: unknown code %q", failure.Code)
	}
	if err := validateFailure(failure, definition); err != nil {
		return fmt.Errorf("write failure envelope: %w", err)
	}

	envelope := failureEnvelope{ContractVersion: ContractVersion, Operation: operation, Error: failure}
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(envelope); err != nil {
		return fmt.Errorf("write failure envelope: %w", err)
	}
	return nil
}
