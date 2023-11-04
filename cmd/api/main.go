package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	version = "1.0.0"
)

// config will hold all the configuration settings for out application
type config struct {
	port int
	env  string
}

// application will hold all the dependencies for out HTTP handlers, helpers and middleware
type application struct {
	config config
	logger *log.Logger
}

func main() {
	// instant of the config struct
	var cfg config

	// use command line arguments to load config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// set up logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// instance of the application struct
	app := &application{
		config: cfg,
		logger: logger,
	}

	// create a HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// start server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
