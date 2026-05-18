package handlers

import (
	"net/http"
	"strconv"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/users/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type userHandler struct {
	us services.UserService
	v  *validator.Validate
}

func NewUserHandler(s services.UserService, v *validator.Validate,
) *userHandler {
	return &userHandler{
		us: s,
		v:  v,
	}
}

func (h *userHandler) CreateUser(c *echo.Context) error {
	req := new(delivery.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Request",
			customs.HandleBindError(err)...,
		))
	}

	if err := h.v.StructCtx(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Validation Failed",
			*customs.NewErrorValue("validation", err.Error()),
		))
	}

	res, err := h.us.CreateUser(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusCreated, response.NewResponseSuccess(res, "User Created"))
}

func (h *userHandler) GetUserByID(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid User ID",
			*customs.NewErrorValue("validation", "id must be a positive integer"),
		))
	}

	res, err := h.us.GetUserByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "User Found"))
}

func (h *userHandler) ListUsers(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	role := c.QueryParam("role")

	res, err := h.us.ListUsers(c.Request().Context(), page, limit, search, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Users Retrieved"))
}

func (h *userHandler) DeleteUser(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid User ID",
			*customs.NewErrorValue("validation", "id must be a positive integer"),
		))
	}

	if err := h.us.DeleteUser(c.Request().Context(), uint(id)); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess[*bool](nil, "User Deleted"))
}
