package exceptions

type InvalideRequestError struct {
	msg string
	error
}

func (e InvalideRequestError) Error() string {
	return e.msg
}

func NewInvalideRequestError(message string) InvalideRequestError {
	return InvalideRequestError{
		msg: message,
	}
}
