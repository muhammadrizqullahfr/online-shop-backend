package handlers

import (
	"net/http"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/labstack/echo/v5"
)

type OrderCallbackHandler struct {
	cs services.OrderCallbackService
}

func NewOrderCallbackHandler(cs services.OrderCallbackService) *OrderCallbackHandler {
	return &OrderCallbackHandler{cs: cs}
}

func (h *OrderCallbackHandler) Callback(c *echo.Context) error {
	req := new(delivery.XenditCallbackRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid payload",
			customs.HandleBindError(err)...,
		))
	}

	if err := h.cs.HandleCallback(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess[any](nil, "Callback processed"))
}
