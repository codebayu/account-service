// create custom middleware
package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

func CustomMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		fmt.Println("Custom middleware executed")
		c.Response().Header().Set("X-Custom-Header", "CustomMiddleware")
		return next(c)
	}
}
