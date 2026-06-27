package links

import "fmt"

type UnknownRelationError struct {
	Relation string
}

func (err *UnknownRelationError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("link add: unknown relation %q", err.Relation)
}

type SelfLinkError struct {
	FromID   string
	Relation string
	ToID     string
}

func (err *SelfLinkError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("link add: self link is not allowed: %s %s %s", err.FromID, err.Relation, err.ToID)
}

type EndpointNotFoundError struct {
	EndpointType string
	EndpointID   string
	Relation     string
	OtherType    string
	OtherID      string
}

func (err *EndpointNotFoundError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"link add: %s %q does not resolve for relation %q between %s %q and %s %q",
		err.EndpointType,
		err.EndpointID,
		err.Relation,
		err.EndpointType,
		err.EndpointID,
		err.OtherType,
		err.OtherID,
	)
}

type IncompatibleRelationError struct {
	Relation string
	FromType string
	FromID   string
	ToType   string
	ToID     string
}

func (err *IncompatibleRelationError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"link add: relation %q is not valid between %s %q and %s %q",
		err.Relation,
		err.FromType,
		err.FromID,
		err.ToType,
		err.ToID,
	)
}

type DuplicateLinkError struct {
	FromID   string
	Relation string
	ToID     string
}

func (err *DuplicateLinkError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("link add: link already exists: %s %s %s", err.FromID, err.Relation, err.ToID)
}
