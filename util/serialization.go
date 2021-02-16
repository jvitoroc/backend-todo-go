package util

import (
	"encoding/json"
	"io"
)

func FromJson(data io.Reader, to interface{}) error {
	if err := json.NewDecoder(data).Decode(&to); err != nil {
		return err
	}

	return nil
}
