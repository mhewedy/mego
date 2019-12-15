package commons

import (
	"encoding/json"
	"io"
)

func JSONDecode(r io.Reader, i interface{}) error {
	err := json.NewDecoder(r).Decode(&i)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
