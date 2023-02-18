package main

import (
	"net/http"
)

// healthcheckHandler gives information on applications status, environment and version
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.wrtieJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
