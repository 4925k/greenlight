package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

// MarshalJSON is an implementation of JSON marshaller on Runtime
func (r Runtime) MarshalJSON() ([]byte, error) {
	// format as required
	jv := fmt.Sprintf("%d mins", r)

	return []byte(strconv.Quote(jv)), nil // wrap the value in double quotes before sending to marshaller
}
