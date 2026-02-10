package exceptions

type WrongPasswordError struct {
	Message string
}

func (e *WrongPasswordError) Error() string {
	return e.Message
}

func NewWrongPasswordError(message string) *WrongPasswordError {
	return &WrongPasswordError{Message: message}
}
