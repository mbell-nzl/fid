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

func WriteSuccess(writer io.Writer, operation string, result map[string]any, event *Event) error {
	if writer == nil {
		return fmt.Errorf("write success envelope: writer is nil")
	}
	if result == nil {
		return fmt.Errorf("write success envelope: result is nil")
	}

	envelope := map[string]any{"contract_version": ContractVersion, "event": event, "ok": true, "operation": operation, "result": result}
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
	definition, ok := failureCodeDefinitions[failure.Code]
	if !ok {
		return fmt.Errorf("write failure envelope: unknown code %q", failure.Code)
	}
	if err := validateFailure(failure, definition); err != nil {
		return fmt.Errorf("write failure envelope: %w", err)
	}

	subject := any(nil)
	if failure.Subject != nil {
		subject = map[string]any{"id": failure.Subject.ID, "type": failure.Subject.Type}
	}
	envelope := map[string]any{
		"contract_version": ContractVersion,
		"error": map[string]any{
			"category": failure.Category, "code": failure.Code, "details": failure.Details,
			"message": failure.Message, "retryable": failure.Retryable, "subject": subject,
		},
		"ok": false, "operation": operation,
	}
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(envelope); err != nil {
		return fmt.Errorf("write failure envelope: %w", err)
	}
	return nil
}
