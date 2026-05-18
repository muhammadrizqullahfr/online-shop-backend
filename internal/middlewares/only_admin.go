package middlewares

import (
	"net/http"

	"github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
	"github.com/labstack/echo/v5"
)

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		user, ok := c.Get("user").(*shared.JwtCustomClaims)
		if !ok || user == nil {
			return c.JSON(http.StatusUnauthorized, response.NewResponseError("Invalid or expired token"))

		}

		if user.Role != constants.RoleAdmin {
			return c.JSON(http.StatusUnauthorized, response.NewResponseError("forbidden: only admin can access this resource"))
		}

		return next(c)
	}
}
