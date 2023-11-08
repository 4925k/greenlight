package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Runtime int32

var (
	ErrInvalidRuntimeFormat = errors.New("invalid runtime format")
)

// MarshalJSON is an implementation of JSON marshaller on Runtime
func (r Runtime) MarshalJSON() ([]byte, error) {
	// format as required
	jv := fmt.Sprintf("%d mins", r)

	return []byte(strconv.Quote(jv)), nil // wrap the value in double quotes before sending to marshaller
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	uqJSONvalue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(uqJSONvalue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}
