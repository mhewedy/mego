package commons

type clientError struct {
	message string
}

func (c clientError) Error() string {
	return c.message
}

func NewClientError(m string) error {
	return &clientError{message: m}
}

func IsClientError(err error) bool {
	_, ok := err.(*clientError)
	return ok
}
