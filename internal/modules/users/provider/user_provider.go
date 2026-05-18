package provider

import (
	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/handlers"

	userService "github.com/RakaMurdiarta/online-shop-system/internal/modules/users/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func UserProvider(
	privateRoute *echo.Group,
	tx *database.TransactionManagerImpl,
	conf *config.Config,
	userService userService.UserService,

) {
	v := validator.New()

	userHandler := handlers.NewUserHandler(userService, v)

	users := privateRoute.Group("/users")
	users.POST("", userHandler.CreateUser)
	users.GET("", userHandler.ListUsers)
	users.GET("/:id", userHandler.GetUserByID)
	users.DELETE("/:id", userHandler.DeleteUser)
}
