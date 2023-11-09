package main

import (
	"fmt"
	"github.com/4925k/greenlight/internal/data"
	"github.com/4925k/greenlight/internal/validator"
	"net/http"
	"time"
)

// createMovieHandler will create a new movie entry
// curl -X POST localhost:4000/v1/movies
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	/*
		The problem with decoding directly into a Movie struct is that a client could provide the
		keys id and version in their JSON request, and the corresponding values would be
		decoded without any error into the ID and Version fields of the Movie struct — even though
		we don’t want them to be. We could check the necessary fields in the Movie struct after the
		event to make sure that they are empty, but that feels a bit hacky, and decoding into an
		intermediary struct (like we are in our handler) is a cleaner, simpler, and more robust
		approach — albeit a little bit verbose
	*/
	var input struct {
		Title   string       `json:"title,omitempty"`
		Year    int32        `json:"year,omitempty"`
		Runtime data.Runtime `json:"runtime,omitempty"`
		Genres  []string     `json:"genres,omitempty"`
	}

	// json.Unmarshal() requires about 80% more memory than json.Decoder, as well as being a tiny bit slower
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// Validation process
	v := validator.New()
	data.ValidateMovie(v, movie)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

// showMovieHandler will return details about the given movie id
// curl localhost:4000/v1/movies/123
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	// get id from URL
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Year:      1998,
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	// return movie details
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
