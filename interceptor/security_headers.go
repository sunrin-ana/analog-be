package interceptor

import (
	"github.com/labstack/echo/v4"
)

func NewSecurityHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			c.Response().Header().Set("X-Frame-Options", "DENY")

			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")

			c.Response().Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

			if c.Request().Header.Get("X-Forwarded-Proto") == "https" {
				c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}

			return next(c)
		}
	}
}
