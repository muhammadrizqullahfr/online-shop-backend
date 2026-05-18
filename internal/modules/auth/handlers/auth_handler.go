package handler

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/constants"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

const (
	oauthStateCookieGoogle = "oauth_state_google"
	oauthStateCookieGithub = "oauth_state_github"
	oauthStateMaxAgeSec    = 600
)

type authHandler struct {
	as     services.AuthService
	v      *validator.Validate
	google services.OauthService
	github services.OauthService
	conf   *config.Config
}

func NewAuthHandler(s services.AuthService, v *validator.Validate, google services.OauthService, github services.OauthService, conf *config.Config) *authHandler {
	return &authHandler{
		as:     s,
		v:      v,
		google: google,
		conf:   conf,
		github: github,
	}
}

func (a *authHandler) SignUpLocal(c *echo.Context) error {
	req := new(delivery.SingUpRequest)

	//bind
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError("Validation Failed", customs.HandleBindError(err)...))
	}

	//validate
	if err := a.v.StructCtx(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError("Validation Failed", *customs.NewErrorValue("validation ", err.Error())))

	}

	res, err := a.as.RegisterLocal(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(err.Error(), *customs.NewErrorValue("bussines_logic", err.Error())))
	}

	return c.JSON(http.StatusCreated, response.NewResponseSuccess(res, "Account Created"))
}

func (a *authHandler) Login(c *echo.Context) error {

	req := new(delivery.LoginRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError("Validation Failed", customs.HandleBindError(err)...))
	}

	//validate
	if err := a.v.StructCtx(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError("Validation Failed", *customs.NewErrorValue("validation ", err.Error())))

	}

	res, err := a.as.LoginLocal(c.Request().Context(), req)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(err.Error(), *customs.NewErrorValue("bussines_logic", err.Error())))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Login Successfully"))

}

func (a *authHandler) LoginGoogle(c *echo.Context) error {
	return a.startOAuthFlow(c, oauthStateCookieGoogle, a.google)
}

func (a *authHandler) LoginGithub(c *echo.Context) error {
	return a.startOAuthFlow(c, oauthStateCookieGithub, a.github)
}

func (a *authHandler) GoogleCallback(c *echo.Context) error {
	return a.handleOAuthCallback(c, oauthStateCookieGoogle, a.google, constants.ProviderGoogle)
}

func (a *authHandler) GithubCallback(c *echo.Context) error {
	return a.handleOAuthCallback(c, oauthStateCookieGithub, a.github, constants.ProviderGithub)
}

// startOAuthFlow issues a fresh CSRF state, persists it in an HTTP-only cookie,
// and 302-redirects the user to the upstream provider's authorize endpoint
// with that same state embedded in the URL.
func (a *authHandler) startOAuthFlow(c *echo.Context, cookieName string, provider services.OauthService) error {
	state, err := generateOAuthState()
	if err != nil {
		return a.redirectWithError(c, "failed to start oauth flow")
	}
	a.setStateCookie(c, cookieName, state)
	return c.Redirect(http.StatusFound, provider.LoginRedirectUrl(state))
}

// handleOAuthCallback validates the CSRF state, exchanges the code for the
// provider's user profile, upserts the local user, issues access + refresh
// JWTs, and 302-redirects to FrontendCallbackURL with token/refresh_token/role
// in the query string. Every error path also 302-redirects (with ?error=...)
// so the browser stays on the frontend.
func (a *authHandler) handleOAuthCallback(
	c *echo.Context,
	cookieName string,
	provider services.OauthService,
	providerKind constants.Provider,
) error {
	if errParam := c.QueryParam("error"); errParam != "" {
		a.clearStateCookie(c, cookieName)
		return a.redirectWithError(c, errParam)
	}

	queryState := c.QueryParam("state")
	cookie, err := c.Cookie(cookieName)
	a.clearStateCookie(c, cookieName)
	if err != nil || cookie.Value == "" || queryState == "" {
		return a.redirectWithError(c, "missing oauth state")
	}
	if subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(queryState)) != 1 {
		return a.redirectWithError(c, "invalid oauth state")
	}

	code := c.QueryParam("code")
	if code == "" {
		return a.redirectWithError(c, "missing authorization code")
	}

	providerResponse, err := provider.Callback(c.Request().Context(), code)
	if err != nil {
		return a.redirectWithError(c, "oauth exchange failed")
	}

	user, err := a.as.HandleOAuthCallback(c.Request().Context(), &delivery.OAuthUserRequest{
		AvatarURL:  providerResponse.Picture,
		Email:      providerResponse.Email,
		Provider:   providerKind,
		ProviderID: providerResponse.Sub,
		FullName:   providerResponse.Name,
	})
	if err != nil {
		return a.redirectWithError(c, "failed to resolve user")
	}

	token, err := shared.GenerateToken(uint(user.ID), user.Role, a.conf.JwtSecretKey, 24)
	if err != nil {
		return a.redirectWithError(c, "failed to generate token")
	}
	refreshToken, err := shared.GenerateToken(uint(user.ID), user.Role, a.conf.JwtSecretKey, 168)
	if err != nil {
		return a.redirectWithError(c, "failed to generate refresh token")
	}

	return a.redirectWithSuccess(c, token, refreshToken, string(user.Role))
}

func (a *authHandler) redirectWithSuccess(c *echo.Context, token, refreshToken, role string) error {
	u, err := url.Parse(a.conf.FrontendCallbackURL)
	if err != nil {
		return c.String(http.StatusInternalServerError, "invalid frontend callback url")
	}
	q := u.Query()
	q.Set("token", token)
	q.Set("refresh_token", refreshToken)
	q.Set("role", role)
	u.RawQuery = q.Encode()
	return c.Redirect(http.StatusFound, u.String())
}

func (a *authHandler) redirectWithError(c *echo.Context, message string) error {
	u, err := url.Parse(a.conf.FrontendCallbackURL)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("oauth error: %s", message))
	}
	q := u.Query()
	q.Set("error", message)
	u.RawQuery = q.Encode()
	return c.Redirect(http.StatusFound, u.String())
}

func (a *authHandler) setStateCookie(c *echo.Context, name, value string) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   oauthStateMaxAgeSec,
		HttpOnly: true,
		Secure:   a.conf.AppEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
}

func (a *authHandler) clearStateCookie(c *echo.Context, name string) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   a.conf.AppEnv == "production",
		SameSite: http.SameSiteLaxMode,
	})
}

func generateOAuthState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (a *authHandler) Me(c *echo.Context) error {
	claims := c.Get("user").(*shared.JwtCustomClaims)

	res, err := a.as.Me(c.Request().Context(), claims.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Profile Retrieved"))
}

func (a *authHandler) RefreshToken(c *echo.Context) error {
	req := new(delivery.RefreshTokenRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError("Validation Failed", customs.HandleBindError(err)...))
	}

	if err := a.v.StructCtx(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError("Validation Failed", *customs.NewErrorValue("validation ", err.Error())))

	}

	res, err := a.as.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(err.Error(), *customs.NewErrorValue("business_logic", err.Error())))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Refresh Token Fetched"))
}
