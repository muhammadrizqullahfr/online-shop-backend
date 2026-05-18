package provider

import (
	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	handler "github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/handlers"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/services"
	ImplService "github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/services/impl"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/repository"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

func AuthProvide(
	privateRoute *echo.Group,
	publicRoute *echo.Group,
	conf *config.Config,
	userRepo repository.UserRepository,
	authService services.AuthService,
) {

	v := validator.New()

	google := ImplService.NewGoogleProvider(userRepo, conf)
	github := ImplService.NewGithubProvider(userRepo, conf)

	authHandler := handler.NewAuthHandler(authService, v, google, github, conf)

	auth := publicRoute.Group("/auth")
	auth.POST("/signup", authHandler.SignUpLocal)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)

	authProtected := privateRoute.Group("/auth")
	authProtected.GET("/me", authHandler.Me)

	oauth := publicRoute.Group("/oauth")
	oauth.GET("/google/login", authHandler.LoginGoogle)
	oauth.GET("/google/cb", authHandler.GoogleCallback)

	oauth.GET("/github/login", authHandler.LoginGithub)
	oauth.GET("/github/cb", authHandler.GithubCallback)

}
