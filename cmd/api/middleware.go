package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Using this pattern for rate-limiting will only work if your API application is running on a
// single-machine. If your infrastructure is distributed, with your application running on
// multiple servers behind a load balancer, then you’ll need to use an alternative approach
func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, user := range clients {
			if time.Since(user.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			mu.Lock()

			if _, ok := clients[ip]; !ok {
				clients[ip] = &client{rate.NewLimiter(rate.Limit(app.config.limiter.rps), int(app.config.limiter.burst)), time.Now()}
			}

			if !clients[ip].limiter.Allow() {
				app.rateLimitExceededResponse(w, r)
				mu.Unlock()
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
