package exceptions

type EventNotFoundError struct {
	Message string
}

func (e *EventNotFoundError) Error() string {
	return e.Message
}
