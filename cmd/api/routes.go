package main

import (
	"github.com/codebayu/account-service/cmd/api/handlers"
)

func (app *Application) routes(h handlers.Handler) {
	app.server.GET("/", h.HealthCheck)
}
