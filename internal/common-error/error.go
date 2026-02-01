package common_error

import "fmt"

type ErrorType string

const (
	TypeValidation        ErrorType = "validation"
	TypeNotFound          ErrorType = "not_found"
	TypeUnauthorized      ErrorType = "unauthorized"
	TypeInternal          ErrorType = "internal"
	EmailAlreadyConfirmed ErrorType = "email_already_confirmed"
)

type ErrorDefinition struct {
	Type    ErrorType
	Message string
}

func NewError(msg string, _type ErrorType) error {
	return &ErrorDefinition{
		Type:    _type,
		Message: msg,
	}
}

func (e *ErrorDefinition) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
