package app

import (
	"context"
	"errors"
	"github.com/askaroe/reservationAPI/internal/handlers"
	"github.com/askaroe/reservationAPI/internal/repository"
	"github.com/askaroe/reservationAPI/internal/services"
	"github.com/askaroe/reservationAPI/pkg/jsonlog"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var logger *jsonlog.Logger

func Run() {
	logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := initDB()

	logger.PrintFatal(err, nil)
	defer db.Close()

	repo := repository.NewReservationRepository(db)

	reservationService := services.NewReservationService(repo)

	reservationHandler := handlers.NewReservationHandler(reservationService)

	r := chi.NewRouter()

	r.Route("/reservations", func(router chi.Router) {
		r.Post("/", reservationHandler.CreateReservation)
		r.Route("/{roomID}", func(r chi.Router) {
			r.Get("/", reservationHandler.GetReservationsByRoomId)
		})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		logger.PrintInfo("starting server on :8080", nil)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.PrintFatal(err, nil)
		}
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	// Create a channel to listen for OS interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Block main thread until interrupt signal is received
	<-quit

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	logger.PrintInfo("shutting down server", nil)
	if err := server.Shutdown(ctx); err != nil {
		logger.PrintFatal(err, nil)
	}

	logger.PrintInfo("server stopped gracefully", nil)
}
