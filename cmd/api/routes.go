package main

import (
	"github.com/codebayu/account-service/internal/middleware"
)

func (app *Application) routes() {
	app.server.GET("/", app.healthHandler.HealthCheck)

	auth := app.server.Group("/auth")
	auth.POST("/register", app.authHandler.Register)
	auth.POST("/login", app.authHandler.Login)
	auth.POST("/refresh-token", app.authHandler.RefreshToken)
	auth.POST("/logout", app.authHandler.Logout, middleware.AuthMiddleware(app.tokenRepo))

	user := app.server.Group("/user", middleware.AuthMiddleware(app.tokenRepo))
	user.GET("/current", app.userHandler.GetCurrentUser)
}
