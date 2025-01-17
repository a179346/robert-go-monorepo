package jsonvalidator

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

func FromReader[T any](reader io.Reader) (T, error) {
	var v T

	if closer, ok := reader.(io.ReadCloser); ok {
		defer closer.Close()
	}

	if err := json.NewDecoder(reader).Decode(&v); err != nil {
		return v, err
	}

	validate := validator.New()
	if err := validate.Struct(v); err != nil {
		return v, err
	}

	return v, nil
}

func FromBytes[T any](data []byte) (T, error) {
	var v T

	if err := json.Unmarshal(data, &v); err != nil {
		return v, err
	}

	validate := validator.New()
	if err := validate.Struct(v); err != nil {
		return v, err
	}

	return v, nil
}
