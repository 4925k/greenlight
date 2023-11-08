package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

type envelope map[string]interface{}

// readIDParam will retrieve the ID URL parameter from the request context and return it as an in64
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return -1, errors.New("invalid id parameter")
	}

	return id, nil
}

/*
	writeJSON will help to send json responses for the api

It’s also possible to use Go’s json.Encoder type to perform the encoding.
This allows you to encode an object to JSON and write that JSON
to an output stream in a single step

When we call json.NewEncoder(w).Encode(data) the JSON is created and written to the http.ResponseWriter in a single step,
which means there’s no opportunity to set HTTP response headers
conditionally based on whether the Encode() method returns an error or not.
*/
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	// marshal data
	// While using json.MarshalIndent() is positive from a readability and user-experience point of view
	// the responses are now slightly larger in terms of total bytes,
	// the extra work that Go does to add the whitespace has a notable impact on performance
	// json.MarshalIndent() takes 65% longer to run and uses around 30% more memory than json.Marshal(),
	// as well as making two more heap allocations.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// readJSON helps to triage the errors and make custom messages when necessary
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// limit the size of request body
	maxBytes := 1_48_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// disallow unknown fields in the json request
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarhsalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON at %d", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}

			return fmt.Errorf("body contains incorrect JSON type at %s", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")

		case errors.As(err, &invalidUnmarhsalError):
			panic(err)

		default:
			return err
		}
	}

	// call decode again to check for additional data
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contains a single JSON value")
	}

	return nil
}
