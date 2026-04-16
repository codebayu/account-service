package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/codebayu/account-service/common/response"
	"github.com/codebayu/account-service/internal/config"
	"github.com/labstack/echo/v5"
)

func SignatureMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Skip validation for swagger and health check
			path := c.Request().URL.Path
			if strings.HasPrefix(path, "/swagger/") || path == "/health" {
				return next(c)
			}

			signature := c.Request().Header.Get("x-signature")
			datetime := c.Request().Header.Get("x-datetime")
			channel := c.Request().Header.Get("x-channel")

			if signature == "" || datetime == "" || channel == "" {
				return c.JSON(http.StatusBadRequest, response.Response{
					Result: response.Result{
						Code:       400000,
						StatusCode: http.StatusBadRequest,
						Message:    "headers x-signature, x-datetime, and x-channel are required",
					},
				})
			}

			// Generate signature: HMAC-SHA256(apiSecret, apiKey + unixTimestamp)
			stringToHash := cfg.APIKey + datetime
			h := hmac.New(sha256.New, []byte(cfg.APISecret))
			h.Write([]byte(stringToHash))
			calculatedSignature := hex.EncodeToString(h.Sum(nil))

			if calculatedSignature != signature {
				return c.JSON(http.StatusUnauthorized, response.Response{
					Result: response.Result{
						Code:       401000,
						StatusCode: http.StatusUnauthorized,
						Message:    "invalid signature",
					},
				})
			}

			if channel != cfg.ChannelID {
				return c.JSON(http.StatusUnauthorized, response.Response{
					Result: response.Result{
						Code:       401001,
						StatusCode: http.StatusUnauthorized,
						Message:    "invalid channel",
					},
				})
			}

			fmt.Println("✅ signature validated")
			return next(c)
		}
	}
}
