package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/magomed066/go-task-queue/internal/json"
)

type application struct {
	addr   string
	config config
}

type config struct {
	db             dbConfig
	workerPoolSize int
	taskTimeout    time.Duration
	leaseDuration  time.Duration
	pollInterval   time.Duration
	maxRetries     int
	logLevel       string
	logFormat      string
}

type dbConfig struct {
	dns string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(25 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	r.Get("/health/live", func(w http.ResponseWriter, r *http.Request) {
		json.WriteResponse(w, http.StatusOK, "Live")
	})

	r.Get("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		json.WriteResponse(w, http.StatusOK, "Ready")
	})

	return r
}

func (app *application) run(ctx context.Context, h http.Handler) error {
	server := http.Server{
		Addr:         app.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	fmt.Printf("Server has started on port %s", app.addr)

	errorCh := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()

		errorCh <- err
	}()

	select {
	case err := <-errorCh:
		return err

	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Second*20,
		)
		defer cancel()

		shutdownErr := server.Shutdown(ctx)
		if shutdownErr != nil {
			_ = server.Close()
		}

		serverErr := <-errorCh
		if errors.Is(serverErr, http.ErrServerClosed) {
			serverErr = nil
		}

		return errors.Join(shutdownErr, serverErr)
	}
}
