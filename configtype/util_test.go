package configtype_test

import "encoding/json"

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func marshalString[T any](v T) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b)[1 : len(b)-1], nil
}
