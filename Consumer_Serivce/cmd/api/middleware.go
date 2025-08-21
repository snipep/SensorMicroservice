package main

import (
	"os"

	"github.com/labstack/echo/v4"
)

// RegisterMiddlewares attaches global middlewares, including JWT auth for all routes
// except the public auth endpoints.
func RegisterMiddlewares(e *echo.Echo) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change"
	}
	appJWTSecret = jwtSecret

	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			// Allow public auth endpoints
			if path == "/api/v1/signup" || path == "/api/v1/signin" {
				return next(c)
			}

			auth := c.Request().Header.Get("Authorization")
			if !validateJWT(auth, jwtSecret) {
				return writeJSONError(c, 401, "invalid or missing token")
			}
			return next(c)
		}
	})
}
