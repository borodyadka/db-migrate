package drivers

import "net/url"

type ExplicitError interface {
	error
	ExplicitError() string
}

type ErrInvalidConnectionString struct {
	url url.URL
	// TODO: error kind: protocol, host, schema/table, etc.
}

func (e *ErrInvalidConnectionString) Error() string {
	return ""
}

func (e *ErrInvalidConnectionString) ExplicitError() string {
	return ""
}

func NewErrInvalidConnectionString(url url.URL) *ErrInvalidConnectionString {
	return &ErrInvalidConnectionString{
		url: url,
	}
}

type InvalidTableNameReason int

var (
	InvalidTableNameReasonEmpty InvalidTableNameReason = 1
)

type ErrInvalidTableName struct {
	t string
	r InvalidTableNameReason
}

func (e *ErrInvalidTableName) Error() string {
	switch e.r {
	case InvalidTableNameReasonEmpty:
		return "name is empty"
	}
	return ""
}

func (e *ErrInvalidTableName) ExplicitError() string {
	return ""
}

func NewErrInvalidTableName(t string, r InvalidTableNameReason) *ErrInvalidTableName {
	return &ErrInvalidTableName{
		t: t,
		r: r,
	}
}

type ErrNotConnected struct{}
