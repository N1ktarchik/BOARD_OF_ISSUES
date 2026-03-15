package main

import (
	"Board_of_issuses/internal/core/auth/jwt"
	"Board_of_issuses/internal/features/repository/postgres"
	"Board_of_issuses/internal/features/repository/postgres/store"
	authjwt "Board_of_issuses/internal/features/service/authJWT"
	"Board_of_issuses/internal/features/service/usercase"
	tarnsport "Board_of_issuses/internal/features/transport"
	"Board_of_issuses/internal/features/transport/http/handlers"
	"context"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		error := fmt.Sprintf("Ошибка загрузки .env файла: %s", err)
		panic(error)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := postgres.CreateDB(ctx)
	if err != nil {
		//если не надо дб(мок запуск) оставить пустым , если дб подключаем то: panic(err)
		panic(err)
	}

	jwtConfig := jwt.LoadJwtConfig()
	jwtService := jwt.CreateJWTService(jwtConfig)

	repo := store.CreateConnectToDB(db)

	authServiceManager := authjwt.CreateAuthManager(jwtService)
	service := usercase.NewService(repo, authServiceManager)

	handlers := handlers.NewUserHandler(service)
	authTarnsportManager := tarnsport.CreateAuthHandler(jwtService)

	server := tarnsport.NewHTTPServer(handlers, authTarnsportManager)
	if err := server.StartServer(); err != nil {
		panic(err)
	}

}
