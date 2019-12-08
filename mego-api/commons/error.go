package commons

type ClientError struct {
	message string
}

func (c ClientError) Error() string {
	return c.message
}

func NewClientError(m string) error {
	return &ClientError{message: m}
}
