package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// createMovieHandler will create a new movie entry
// curl -X POST localhost:4000/v1/movies
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a movie")
}

// showMovieHandler will return details about the given movie id
// curl localhost:4000/v1/movies/123
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// using http router to parse a request, any interpolated URL params will be stored in the request context
	params := httprouter.ParamsFromContext(r.Context())

	// fetch the movie id
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// show the movie details
	fmt.Fprintf(w, "showing details about movie %d\n", id)
}
