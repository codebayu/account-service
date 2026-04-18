package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/codebayu/account-service/common/response"
	"github.com/codebayu/account-service/internal/config"
	"github.com/labstack/echo/v5"
)

const signatureWindow = 5 * time.Minute

func SignatureMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
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

			// 1. Validate timestamp — prevent replay attacks
			ts, err := strconv.ParseInt(datetime, 10, 64)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, response.Response{
					Result: response.Result{
						Code:       401002,
						StatusCode: http.StatusUnauthorized,
						Message:    "invalid x-datetime format, must be unix timestamp",
					},
				})
			}

			requestTime := time.Unix(ts, 0)
			diff := time.Since(requestTime)
			if diff > signatureWindow || diff < -signatureWindow {
				return c.JSON(http.StatusUnauthorized, response.Response{
					Result: response.Result{
						Code:       401003,
						StatusCode: http.StatusUnauthorized,
						Message:    "request expired, timestamp out of allowed window",
					},
				})
			}

			// 2. Validate HMAC — constant-time compare (prevent timing attacks)
			stringToHash := cfg.APIKey + datetime
			h := hmac.New(sha256.New, []byte(cfg.APISecret))
			h.Write([]byte(stringToHash))
			calculatedSignature := hex.EncodeToString(h.Sum(nil))

			if !hmac.Equal([]byte(calculatedSignature), []byte(signature)) {
				return c.JSON(http.StatusUnauthorized, response.Response{
					Result: response.Result{
						Code:       401000,
						StatusCode: http.StatusUnauthorized,
						Message:    "invalid signature",
					},
				})
			}

			// 3. Validate channel
			if channel != cfg.ChannelID {
				return c.JSON(http.StatusUnauthorized, response.Response{
					Result: response.Result{
						Code:       401001,
						StatusCode: http.StatusUnauthorized,
						Message:    "invalid channel",
					},
				})
			}

			return next(c)
		}
	}
}
