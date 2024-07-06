package errors

type ErrorResponse struct {
	Error string  `json:"error"`
	Code  *string `json:"code,omitempty"`
	Field *string `json:"field,omitempty"`
}
