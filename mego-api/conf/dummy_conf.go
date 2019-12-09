package conf

import (
	"io"
	"strings"
)

type DummySource struct {
}
type dummyReaderCloser struct {
	io.Reader
}

func (d dummyReaderCloser) Close() error {
	return nil
}

func (d DummySource) Read() (io.ReadCloser, error) {
	return &dummyReaderCloser{strings.NewReader(``)}, nil
}
