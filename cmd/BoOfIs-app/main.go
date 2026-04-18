package main

import (
	auth_service "Board_of_issuses/internal/features/auth"

	users_repository "Board_of_issuses/internal/features/users/repository"
	users_service "Board_of_issuses/internal/features/users/service"
	users_transport_http "Board_of_issuses/internal/features/users/transport/http"

	desks_repository "Board_of_issuses/internal/features/desks/repository"
	desks_service "Board_of_issuses/internal/features/desks/service"
	desks_transport_http "Board_of_issuses/internal/features/desks/transport/http"

	tasks_repository "Board_of_issuses/internal/features/tasks/repository"
	tasks_service "Board_of_issuses/internal/features/tasks/service"
	tasks_transport_http "Board_of_issuses/internal/features/tasks/transport/http"

	"Board_of_issuses/internal/core/logger"
	"Board_of_issuses/internal/core/repository/postgres"
	"Board_of_issuses/internal/core/transport/server"
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// @title 			Board Of Issuses
// @version 		2.0
// @description 	Golang Board Of Issuses API

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Enter the token in the format: Bearer <JWT_TOKEN>

// @host 			localhost:8080
// @BasePath 		/
func main() {
	log := logger.Setup()
	slog.SetDefault(log)

	log.Info("starting Board Of Issuses app")

	if err := godotenv.Load(); err != nil {
		log.Error("godotenv: .env file not found")
		panic(".env file not found")
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Error("SECRET_KEY not set, using default value")
		panic("SECRET_KEY is not set")
	}

	connStr := postgres.GetPostgresValues()
	config := postgres.NewDatabaseConfig(connStr, 25, 5, 30*time.Minute, 5*time.Minute, 1*time.Minute)
	pool, err := postgres.CreatePool(context.Background(), config, log)
	if err != nil {
		log.Error("failed to initialize database pool", slog.Any("err", err))
		panic(err)
	}

	log.Info("database connection pool established")

	jwt_liveTime := os.Getenv("JWT_LIVE_TIME")
	jwtLiveTimeMinutes := 15

	if jwt_liveTime != "" {
		if val, err := strconv.Atoi(jwt_liveTime); err == nil {
			jwtLiveTimeMinutes = val
		} else {
			log.Warn("jwt live time parsing error, using default", slog.String("val", jwt_liveTime))
		}
	} else {
		log.Warn("jwt live time is not set, using default")
	}

	authService := auth_service.CreateJWTService(secret, log, jwtLiveTimeMinutes)

	usersRepository := users_repository.NewUsersRepository(pool, log)
	usersService := users_service.NewUsersService(usersRepository, authService, log)
	usersTransportHttp := users_transport_http.NewUsersHandler(usersService, log)

	desksRepository := desks_repository.NewDesksRepository(pool, log)
	desksService := desks_service.NewDesksService(desksRepository, log)
	desksTransportHttp := desks_transport_http.NewDesksHandler(desksService, log)

	tasksRepository := tasks_repository.NewTasksRepository(pool, log)
	tasksService := tasks_service.NewTasksService(tasksRepository, log)
	tasksTransportHttp := tasks_transport_http.NewTasksHandler(tasksService, log)

	mw := server.NewMiddleWare(authService, log)
	srv := server.NewServer(log)
	srv.RegisterSwagger()
	r := srv.Router

	r.HandleFunc("/register", usersTransportHttp.RegisterUser).Methods("POST")
	r.HandleFunc("/login", usersTransportHttp.LoginUser).Methods("POST")

	api := r.PathPrefix("/").Subrouter()
	api.Use(mw.AuthMiddleware)

	api.HandleFunc("/users/update", usersTransportHttp.ChangesUserData).Methods("PATCH")

	api.HandleFunc("/desks/create", desksTransportHttp.CreateDesk).Methods("POST")
	api.HandleFunc("/desks/{id}", desksTransportHttp.DeleteDesk).Methods("DELETE")
	api.HandleFunc("/desks/my", desksTransportHttp.GetUsersDesks).Methods("GET")
	api.HandleFunc("/desks/connect", desksTransportHttp.ConnectUserToDesk).Methods("POST")
	api.HandleFunc("/desks/update", tasksTransportHttp.ChangeTaskData).Methods("PATCH")

	api.HandleFunc("/tasks/create", tasksTransportHttp.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}/complete", tasksTransportHttp.CompleteTask).Methods("PATCH")
	api.HandleFunc("/tasks/update", tasksTransportHttp.ChangeTaskData).Methods("PATCH")
	api.HandleFunc("/tasks/{id}", tasksTransportHttp.DeleteTask).Methods("DELETE")
	api.HandleFunc("/tasks/all/{deskId}", tasksTransportHttp.GetTasksFromOneDesk).Methods("GET")
	api.HandleFunc("/tasks/{taskId}", tasksTransportHttp.GetTaskByID).Methods("GET")

	log.Info("all services initialized, transport starting", slog.String("addr", ":8080"))
	if err := srv.Run(":8080"); err != nil {
		log.Error("server crashed", slog.Any("err", err))
	}
}
