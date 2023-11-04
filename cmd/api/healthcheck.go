package main

import (
	"fmt"
	"net/http"
)

// healthcheckHandler returns information about the server
// curl -i localhost:4000/v1/healthcheck
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
