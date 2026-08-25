package issues

import "fmt"

// CapabilityUnavailableError marks a deferred judgment mutation with no writer.
type CapabilityUnavailableError struct {
	CapabilityID string
	OwnerPhase   int
}

func (err *CapabilityUnavailableError) Error() string {
	return fmt.Sprintf("capability %q is unavailable until phase %d", err.CapabilityID, err.OwnerPhase)
}

func unavailableCapability(capabilityID string, ownerPhase int) error {
	return &CapabilityUnavailableError{CapabilityID: capabilityID, OwnerPhase: ownerPhase}
}
