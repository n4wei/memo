package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/n4wei/memo/api"
	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/db/cache"
	"github.com/n4wei/memo/lib/logger"
)

const (
	defaultServerShutdownTimeout = 5 * time.Second
)

func main() {
	var portStr string
	flag.StringVar(&portStr, "port", "8080", "set the port for the server to listen on")
	flag.Parse()

	logger := logger.New()

	dbClient, err := cache.New(logger)
	if err != nil {
		logger.Errorf("error starting db client: %v", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    ":" + portStr,
		Handler: api.New(dbClient, logger),
	}

	idleConnsClosed := make(chan struct{})
	go handleInterrupt(server, dbClient, logger, idleConnsClosed)

	logger.Printf("starting server on port %s...", portStr)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Errorf("error starting server: %v", err)
		os.Exit(1)
	}

	<-idleConnsClosed
	logger.Println("server shutdown successful. byebye")
}

func handleInterrupt(server *http.Server, dbClient db.Client, logger logger.Logger, idleConnsClosed chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	logger.Println("caught interrupt, shutting down server...")

	err := dbClient.Close()
	if err != nil {
		logger.Errorf("error closing db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultServerShutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		// error from closing listeners, or context timeout:
		logger.Errorf("error shutting down server: %v", err)
	}

	close(idleConnsClosed)
}
