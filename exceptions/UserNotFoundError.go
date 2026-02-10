package exceptions

type UserNotFoundError struct {
	Message string
}

func (e *UserNotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *UserNotFoundError {
	return &UserNotFoundError{Message: message}
}
