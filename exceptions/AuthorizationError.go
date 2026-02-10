package exceptions

type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

func NewAuthorizationError(message string) *AuthorizationError {
	return &AuthorizationError{Message: message}
}
