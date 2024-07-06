package errors

import (
	"fmt"
	"strings"
)

type SError struct {
	Info         string
	Code         *string
	Field        *string
	IsBadRequest bool
}

func (e *SError) Error() string {
	return e.Info
}

func New(code, text string, isBadRequest bool) *SError {
	return &SError{
		Info:         text,
		Code:         &code,
		IsBadRequest: isBadRequest,
	}
}

func (e *SError) WithField(field string) *SError {
	return &SError{
		Info:         e.Info,
		Code:         e.Code,
		Field:        &field,
		IsBadRequest: e.IsBadRequest,
	}
}

func (e *SError) WithInfo(text string) *SError {
	var newInfo string
	if e.Info == "" {
		newInfo = text
	} else {
		if strings.HasPrefix(text, e.Info) {
			newInfo = text
		} else {
			newInfo = fmt.Sprintf("%s : %s", e.Info, text)
		}
	}
	return &SError{
		Info:         newInfo,
		Code:         e.Code,
		Field:        e.Field,
		IsBadRequest: e.IsBadRequest,
	}
}

func (e *SError) WithError(err error) *SError {
	return e.WithInfo(err.Error())
}

func (e *SError) ToResponse() *ErrorResponse {
	return &ErrorResponse{
		Error: e.Info,
		Code:  e.Code,
		Field: e.Field,
	}
}
