package utils

import (
	"encoding/json"
	"io"
)

func DecodeJson[T any](body io.ReadCloser) (*T, error) {
	defer body.Close()

	var data T

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}