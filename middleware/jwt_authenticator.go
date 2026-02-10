package middleware

import (
	"bookingservice/services"
	"errors"
	"strings"

	"github.com/gobuffalo/buffalo"
)

func JWTAuthenticator() buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			// Skip authentication for login and health check routes
			path := c.Request().URL.Path
			if path == "/users/login" || path == "/health" || path == "/users/login/" || path == "/health/" {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.Error(401, errors.New("missing authorization header"))
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			userService := services.NewUserService()
			user, err := userService.ValidateAuthenticationToken(tokenString)
			if err != nil {
				panic(err)
			}
			c.Set("user_info", user)
			return next(c)
		}
	}
}
