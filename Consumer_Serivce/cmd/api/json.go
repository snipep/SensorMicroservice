package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func writeJSON(c echo.Context, status int, data any) error {

	return c.JSON(status, data)
}

func readJSON(c echo.Context, data any) error {
	// Limit body size (1MB max)
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 1_048_576)

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJSONError(c echo.Context, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return writeJSON(c, status, &envelope{Error: message})
}

func (app *Application) jsonResponse(c echo.Context, status int, data any) error {
	type envelop struct {
		Data any `json:"data"`
	}
	return writeJSON(c, status, &envelop{Data: data})
}

var appJWTSecret string

func (app *Application) generateJWT(userID int64, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "consumer-service",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(appJWTSecret))
}

func validateJWT(authHeader string, secret string) bool {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return false
	}
	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.ErrUnauthorized
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	return true
}
