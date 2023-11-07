package main

import (
	"fmt"
	"github.com/4925k/greenlight/internal/data"
	"net/http"
	"time"
)

// createMovieHandler will create a new movie entry
// curl -X POST localhost:4000/v1/movies
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a movie")
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
