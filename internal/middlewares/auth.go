package middlewares

import (
	"net/http"
	"strings"

	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func AuthMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, response.NewResponseError("Missing authorization header", *customs.NewErrorValue("auth", "Missing authorization header")))
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, response.NewResponseError("Invalid token format"))
			}

			tokenString := parts[1]

			claims := &shared.JwtCustomClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, response.NewResponseError("Invalid or expired token"))
			}

			claims, ok := token.Claims.(*shared.JwtCustomClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, response.NewResponseError("Invalid token claims"))
			}

			c.Set("user", claims)

			return next(c)
		}
	}
}
