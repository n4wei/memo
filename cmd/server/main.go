package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/n4wei/memo/api"
	"github.com/n4wei/memo/db/memo"
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

	server := &http.Server{
		Addr:    ":" + portStr,
		Handler: api.NewController(memo.New(logger), logger),
	}

	idleConnsClosed := make(chan struct{})
	go handleInterrupt(server, logger, idleConnsClosed)

	logger.Printf("starting server on %s...", portStr)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Errorf("server error: %v", err)
		os.Exit(1)
	}

	<-idleConnsClosed
	logger.Println("server shutdown successful. byebye")
}

func handleInterrupt(server *http.Server, logger logger.Logger, idleConnsClosed chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	logger.Println("caught interrupt, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), defaultServerShutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		// error from closing listeners, or context timeout:
		logger.Errorf("server shutdown error: %v", err)
	}

	close(idleConnsClosed)
}
