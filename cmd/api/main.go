package main

import (
	"fmt"
	"os"

	"github.com/codebayu/account-service/cmd/api/handlers"
	"github.com/codebayu/account-service/cmd/api/middlewares"
	"github.com/codebayu/account-service/cmd/api/services"
	"github.com/codebayu/account-service/internal/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Application struct {
	server  *echo.Echo
	handler handlers.Handler
}

func main() {
	err := godotenv.Load(".env")
	e := echo.New()
	if err != nil {
		e.Logger.Error("Error loading .env file")
	}

	db, err := database.NewPostgres()
	if err != nil {
		e.Logger.Error(err.Error())
	}

	authService := services.AuthService{
		DB: db,
	}

	h := handlers.Handler{
		DB:          db,
		AuthService: authService,
	}
	app := Application{
		server:  e,
		handler: h,
	}

	e.Use(middlewares.CustomMiddleware, middleware.RequestLogger())
	app.routes(h)
	fmt.Println(app)
	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)

	if err := e.Start(appAddress); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
