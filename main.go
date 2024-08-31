package main

import (
	"github.com/askaroe/reservationAPI/internal/handlers"
	"github.com/askaroe/reservationAPI/internal/initializers"
	"github.com/askaroe/reservationAPI/internal/repository"
	"github.com/askaroe/reservationAPI/internal/server"
	"github.com/askaroe/reservationAPI/internal/services"
	"github.com/askaroe/reservationAPI/pkg/jsonlog"
	"github.com/askaroe/reservationAPI/pkg/router"
	"github.com/go-chi/chi/v5"
	"os"
)

var logger *jsonlog.Logger

func main() {
	// init logger
	logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	//init env. variables
	err := initializers.LoadEnvVariables()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	logger.PrintInfo("env. variables loaded successfully", nil)

	// init db
	db, err := initializers.InitDB()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("connection to db was established", nil)

	// init repo
	repo := repository.NewReservationRepository(db)

	// init service
	reservationService := services.NewReservationService(repo)

	// init handler
	reservationHandler := handlers.NewReservationHandler(reservationService)

	// init router
	r := router.NewRouter()
	r.Route("/reservations", func(router chi.Router) {
		r.Post("/", reservationHandler.CreateReservation)
		r.Route("/{roomID}", func(r chi.Router) {
			r.Get("/", reservationHandler.GetReservationsByRoomId)
		})
	})

	srv := server.NewServer(r, logger)
	srv.Start()
	srv.Shutdown()
}
