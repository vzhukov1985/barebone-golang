// Schema generated model - DO NOT EDIT!!!

package errors

var (
	LockError    = New("lockError", "lock error", false)
	ParseRequest = New("parseRequest", "request parse error", true)
	Undefined    = New("undefined", "undefined", false)
)

var Slice = []*SError{
	LockError,
	ParseRequest,
	Undefined,
}
