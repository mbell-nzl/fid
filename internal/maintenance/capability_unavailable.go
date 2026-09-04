package maintenance

import "fmt"

// CapabilityUnavailableError marks a deferred maintenance mutation with no writer.
type CapabilityUnavailableError struct {
	CapabilityID string
	OwnerPhase   int
}

func (err *CapabilityUnavailableError) Error() string {
	return fmt.Sprintf("capability %q is unavailable until phase %d", err.CapabilityID, err.OwnerPhase)
}
