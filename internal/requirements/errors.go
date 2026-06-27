package requirements

import "fmt"

// UnknownTypeError reports an unsupported --type filter.
type UnknownTypeError struct {
	Type string
}

func (err *UnknownTypeError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("requirement sync: --type must be PRD or TRD, got %q", err.Type)
}

// NoParseableIDsError reports that the source document had no canonical
// requirement rows after applying the requested filters.
type NoParseableIDsError struct {
	SourceDoc string
}

func (err *NoParseableIDsError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("requirement sync: no parseable requirement ids found in %q", err.SourceDoc)
}

// DuplicateIDError reports a repeated canonical requirement ID within one
// source document.
type DuplicateIDError struct {
	ID string
}

func (err *DuplicateIDError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("requirement sync: duplicate requirement id %q in source document", err.ID)
}

// RequirementNotFoundInSourceError reports a requested --id that is absent from
// the source document.
type RequirementNotFoundInSourceError struct {
	ID        string
	SourceDoc string
}

func (err *RequirementNotFoundInSourceError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("requirement sync: requirement id %q not found in %q", err.ID, err.SourceDoc)
}
