package handlers

import (
	"net/http"
	"strconv"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type ProductHandler struct {
	ps services.ProductService
	v  *validator.Validate
}

func NewProductHandler(ps services.ProductService, v *validator.Validate) *ProductHandler {
	return &ProductHandler{ps: ps, v: v}
}

func (h *ProductHandler) Create(c *echo.Context) error {
	req := new(delivery.CreateNewProductRequest)
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

	res, err := h.ps.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusCreated, response.NewResponseSuccess(res, "Product Created"))
}

func (h *ProductHandler) Update(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Product ID",
			*customs.NewErrorValue("validation", "id must be a valid UUID"),
		))
	}

	req := new(delivery.UpdateNewProductRequest)
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

	res, err := h.ps.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Product Updated"))
}

func (h *ProductHandler) GetByID(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Product ID",
			*customs.NewErrorValue("validation", "id must be a valid UUID"),
		))
	}

	res, err := h.ps.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Product Found"))
}

func (h *ProductHandler) GetAll(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	category := c.QueryParam("category")

	res, err := h.ps.GetAll(c.Request().Context(), page, limit, search, category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Products Retrieved"))
}

func (h *ProductHandler) Delete(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Product ID",
			*customs.NewErrorValue("validation", "id must be a valid UUID"),
		))
	}

	if err := h.ps.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess[*bool](nil, "Product Deleted"))
}
