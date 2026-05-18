package handlers

import (
	"net/http"
	"strconv"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type OrderHandler struct {
	os services.OrderService
	v  *validator.Validate
}

func NewOrderHandler(os services.OrderService, v *validator.Validate) *OrderHandler {
	return &OrderHandler{os: os, v: v}
}

func (h *OrderHandler) CreateOrder(c *echo.Context) error {
	req := new(delivery.CreateOrderRequest)
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

	var userID *string
	if claims, ok := c.Get("user").(*shared.JwtCustomClaims); ok && claims != nil {
		uid := strconv.FormatUint(uint64(claims.UserID), 10)
		userID = &uid
	}

	res, err := h.os.CreateOrder(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusCreated, response.NewResponseSuccess(res, "Order Created"))
}

func (h *OrderHandler) GetOrderByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Order ID",
			*customs.NewErrorValue("validation", "id is required"),
		))
	}

	res, err := h.os.GetOrderByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Order Found"))
}

func (h *OrderHandler) ListOrders(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	status := c.QueryParam("status")
	userID := c.QueryParam("user_id")

	res, err := h.os.ListOrders(c.Request().Context(), page, limit, search, status, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Orders Retrieved"))
}

func (h *OrderHandler) GetMyOrders(c *echo.Context) error {
	claims, ok := c.Get("user").(*shared.JwtCustomClaims)
	if !ok || claims == nil {
		return c.JSON(http.StatusUnauthorized, response.NewResponseError(
			"Unauthorized",
			*customs.NewErrorValue("auth", "missing user claims"),
		))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	status := c.QueryParam("status")
	userID := strconv.FormatUint(uint64(claims.UserID), 10)

	res, err := h.os.ListOrders(c.Request().Context(), page, limit, search, status, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "My Orders Retrieved"))
}

func (h *OrderHandler) GetOrdersByUserID(c *echo.Context) error {
	userID := c.Param("userId")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid User ID",
			*customs.NewErrorValue("validation", "userId is required"),
		))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	status := c.QueryParam("status")

	res, err := h.os.ListOrders(c.Request().Context(), page, limit, search, status, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Orders Retrieved"))
}

func (h *OrderHandler) UpdateOrder(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Order ID",
			*customs.NewErrorValue("validation", "id is required"),
		))
	}

	req := new(delivery.UpdateOrderRequest)
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

	res, err := h.os.UpdateOrder(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Order Updated"))
}

func (h *OrderHandler) DeleteOrder(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Order ID",
			*customs.NewErrorValue("validation", "id is required"),
		))
	}

	if err := h.os.DeleteOrder(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess[*bool](nil, "Order Deleted"))
}

func (h *OrderHandler) AddItem(c *echo.Context) error {
	orderID := c.Param("id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Order ID",
			*customs.NewErrorValue("validation", "id is required"),
		))
	}

	req := new(delivery.CreateOrderItemRequest)
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

	res, err := h.os.AddItem(c.Request().Context(), orderID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusCreated, response.NewResponseSuccess(res, "Item Added"))
}

func (h *OrderHandler) UpdateItem(c *echo.Context) error {
	orderID := c.Param("id")
	itemID := c.Param("itemId")
	if orderID == "" || itemID == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Path",
			*customs.NewErrorValue("validation", "order id and item id are required"),
		))
	}

	req := new(delivery.UpdateOrderItemRequest)
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

	res, err := h.os.UpdateItem(c.Request().Context(), orderID, itemID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Item Updated"))
}

func (h *OrderHandler) DeleteItem(c *echo.Context) error {
	orderID := c.Param("id")
	itemID := c.Param("itemId")
	if orderID == "" || itemID == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Path",
			*customs.NewErrorValue("validation", "order id and item id are required"),
		))
	}

	if err := h.os.DeleteItem(c.Request().Context(), orderID, itemID); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess[*bool](nil, "Item Deleted"))
}

func (h *OrderHandler) ListItems(c *echo.Context) error {
	orderID := c.Param("id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Order ID",
			*customs.NewErrorValue("validation", "id is required"),
		))
	}

	res, err := h.os.ListItems(c.Request().Context(), orderID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	if res == nil {
		res = []delivery.OrderItemResponse{}
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Items Retrieved"))
}
