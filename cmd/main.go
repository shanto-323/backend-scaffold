package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/shanto-323/backend-scaffold/config"
	"github.com/shanto-323/backend-scaffold/internal/server"
	"github.com/shanto-323/backend-scaffold/internal/server/handler"
	"github.com/shanto-323/backend-scaffold/internal/server/router"
)

const CleaningTime time.Duration = 1 * time.Second

func main() {
	logger := zerolog.New(os.Stdout)

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config %w", err)
	}

	s, err := server.NewServer(&logger, config)
	if err != nil {
		log.Fatal("Error creating new server %w", err)
	}

	h := handler.NewHandlers(s)

	r := router.NewRouter(s, h)

	stopChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		s.SetUpHTTPServer(r)
		if err := s.Run(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-stopChan:
		log.Printf("Stopping server in %d sec \n", int(CleaningTime.Seconds()))
		ctx, cancel := context.WithTimeout(context.Background(), CleaningTime)
		defer cancel()

		if err := s.Stop(); err != nil {
			errChan <- err
		}

		select {
		case <-ctx.Done():
			log.Fatal("Error stopping server")
		case err := <-errChan:
			log.Fatal("Error stopping server %w", err)
		}
	case err := <-errChan:
		log.Fatal("Error running server %w", err)
	}
}
