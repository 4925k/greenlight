package main

import (
	"context"
	"github.com/4925k/greenlight/internal/data"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

// contextSetUser sets the user struct into the context
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser fetches the user struct from the context
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
