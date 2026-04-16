package main

import (
	"github.com/codebayu/account-service/cmd/api/handlers"
	"github.com/codebayu/account-service/cmd/api/middlewares"
)

func (app *Application) routes(h handlers.Handler) {
	app.server.GET("/", h.HealthCheck)

	auth := app.server.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	user := app.server.Group("/user", middlewares.AuthMiddleware)
	user.GET("/current", h.GetCurrentUser)
}
