package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Define a custom key to avoid context key collisions
type contextKey string

const pubKeyCtxKey = contextKey("public_key")

// AuthMiddleware validates the JWT and adds token data to the context
func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing or invalid Authorization header",
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate the JWT token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			// Extract public key from token claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token claims",
				})
			}

			publicKey, _ := claims["public_key"].(string)

			log.Println("Public key:", publicKey)

			// Add the public key to the request context
			ctx := context.WithValue(c.Request().Context(), "public_key", publicKey)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
