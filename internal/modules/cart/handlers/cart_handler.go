package handlers

import (
	"net/http"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type CartHandler struct {
	cs services.CartService
	v  *validator.Validate
}

func NewCartHandler(cs services.CartService, v *validator.Validate) *CartHandler {
	return &CartHandler{cs: cs, v: v}
}

func (h *CartHandler) GetCart(c *echo.Context) error {
	user := c.Get("user").(*shared.JwtCustomClaims)

	res, err := h.cs.GetCart(c.Request().Context(), user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Cart Retrieved"))
}

func (h *CartHandler) AddItem(c *echo.Context) error {
	req := new(delivery.AddItemRequest)
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

	user := c.Get("user").(*shared.JwtCustomClaims)

	res, err := h.cs.AddItem(c.Request().Context(), user.UserID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusCreated, response.NewResponseSuccess(res, "Item Added To Cart"))
}

func (h *CartHandler) UpdateItemQuantity(c *echo.Context) error {
	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Item ID",
			*customs.NewErrorValue("validation", "itemId must be a valid UUID"),
		))
	}

	req := new(delivery.UpdateItemQuantityRequest)
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

	user := c.Get("user").(*shared.JwtCustomClaims)

	res, err := h.cs.UpdateItemQuantity(c.Request().Context(), user.UserID, itemID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Cart Item Updated"))
}

func (h *CartHandler) DeleteItem(c *echo.Context) error {
	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Item ID",
			*customs.NewErrorValue("validation", "itemId must be a valid UUID"),
		))
	}

	user := c.Get("user").(*shared.JwtCustomClaims)

	if err := h.cs.DeleteItem(c.Request().Context(), user.UserID, itemID); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess[interface{}](nil, "Cart Item Deleted"))
}

func (h *CartHandler) ClearCart(c *echo.Context) error {
	user := c.Get("user").(*shared.JwtCustomClaims)

	if err := h.cs.ClearCart(c.Request().Context(), user.UserID); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess[interface{}](nil, "Cart Cleared"))
}
