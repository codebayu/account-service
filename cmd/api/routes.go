package main

import (
	"github.com/codebayu/account-service/internal/middleware"
)

func (app *Application) routes() {
	app.server.GET("/", app.healthHandler.HealthCheck)

	auth := app.server.Group("/auth")
	auth.POST("/register", app.authHandler.Register)
	auth.POST("/login", app.authHandler.Login)

	user := app.server.Group("/user", middleware.AuthMiddleware)
	user.GET("/current", app.userHandler.GetCurrentUser)
}
