package main

import (
	"fmt"
	"log"

	"github.com/codebayu/account-service/internal/config"
	"github.com/codebayu/account-service/internal/database"
	"github.com/codebayu/account-service/internal/handler"
	"github.com/codebayu/account-service/internal/repository"
	"github.com/codebayu/account-service/internal/service"
	"github.com/codebayu/account-service/internal/middleware"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

type Application struct {
	server       *echo.Echo
	authHandler  *handler.AuthHandler
	userHandler  *handler.UserHandler
	healthHandler *handler.HealthHandler
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Repository
	userRepo := repository.NewUserRepository(db)

	// Service
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	// Handler
	app := Application{
		server:        echo.New(),
		authHandler:   handler.NewAuthHandler(authService),
		userHandler:   handler.NewUserHandler(userService),
		healthHandler: handler.NewHealthHandler(),
	}

	// Middleware
	app.server.Use(echoMiddleware.RequestLogger())
	app.server.Use(middleware.SignatureMiddleware(cfg))

	// Routes
	app.routes()

	appAddress := fmt.Sprintf(":%s", cfg.AppPort)
	if err := app.server.Start(appAddress); err != nil {
		log.Fatal("failed to start server", err)
	}
}
