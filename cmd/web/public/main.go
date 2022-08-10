package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"

	"github.com/marchelbling/webbox/cmd/web/public/app"
)

// RequestProcessingTimeout is the maximum time allowed to process a single request.
const RequestProcessingTimeout = 10 * time.Second

// RequestConnectionTimeout is the maximum time allowed for which the HTTP server will allow to read/write
// an incoming request. It should be strictly higher than RequestProcessingTimeout.
const RequestConnectionTimeout = RequestProcessingTimeout + 5*time.Second

var version string

// Configuration defines env config expected by this command.
type Configuration struct {
	AppHost string `envconfig:"APP_HOST" default:""`
	AppPort string `envconfig:"APP_PORT" default:"80"`
}

func main() {
	var config Configuration
	envconfig.MustProcess("", &config)

	app, err := app.New()
	if err != nil {
		log.Fatalf("app: %v", err)
	}

	addr := net.JoinHostPort(config.AppHost, config.AppPort)
	// check (for timeouts): https://github.com/golang/go/issues/16100
	server := &http.Server{
		Addr:         addr,
		Handler:      app.Routes(),
		ReadTimeout:  RequestConnectionTimeout,
		WriteTimeout: RequestConnectionTimeout,
	}

	done := make(chan struct{})

	go handleInterrupts(server, done)

	log.Printf("Listening on %s...\n", server.Addr)
	err = server.ListenAndServe()
	switch err {
	case http.ErrServerClosed:
		<-done // wait for graceful shutdown
	default: // ListenAndServe always returns a non-nil error
		log.Fatalf("http server: %v", err)
	}
	log.Print("Server stopped")
}

func handleInterrupts(server *http.Server, done chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	sig := <-sigs
	log.Printf("Received signal %v, stopping...", sig)

	timeout := RequestProcessingTimeout + time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Print("Server shutting down...")
	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("shutdown: %v", err)
	}
	done <- struct{}{}
}
