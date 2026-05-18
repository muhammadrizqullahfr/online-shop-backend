package handlers

import (
	"net/http"
	"strconv"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type CategoryHandler struct {
	cs services.CategoryService
	v  *validator.Validate
}

func NewCategoryHandler(categoryService services.CategoryService, v *validator.Validate,
) *CategoryHandler {
	return &CategoryHandler{
		cs: categoryService,
		v:  v,
	}
}

func (h *CategoryHandler) CreateCategory(c *echo.Context) error {
	req := new(delivery.CreateCategoryRequest)

	// bind
	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				"Validation Failed",
				customs.HandleBindError(err)...,
			),
		)
	}

	// validate
	if err := h.v.StructCtx(c.Request().Context(), req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				"Validation Failed",
				*customs.NewErrorValue("validation", err.Error()),
			),
		)
	}

	res, err := h.cs.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				err.Error(),
				*customs.NewErrorValue("business_logic", err.Error()),
			),
		)
	}

	return c.JSON(
		http.StatusCreated,
		response.NewResponseSuccess(res, "Category Created"),
	)
}

func (h *CategoryHandler) GetCategoryByID(c *echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				"Invalid Category ID",
				*customs.NewErrorValue("validation", "id must be number"),
			),
		)
	}

	res, err := h.cs.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			response.NewResponseError(
				err.Error(),
				*customs.NewErrorValue("business_logic", err.Error()),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.NewResponseSuccess(res, "Category Found"),
	)
}

func (h *CategoryHandler) GetCategoryBySlug(c *echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				"Invalid Category Slug",
				*customs.NewErrorValue("validation", "slug is required"),
			),
		)
	}

	res, err := h.cs.GetBySlug(c.Request().Context(), slug)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			response.NewResponseError(
				err.Error(),
				*customs.NewErrorValue("business_logic", err.Error()),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.NewResponseSuccess(res, "Category Found"),
	)
}

func (h *CategoryHandler) GetAllCategories(c *echo.Context) error {

	res, err := h.cs.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.NewResponseError(
				err.Error(),
				*customs.NewErrorValue("business_logic", err.Error()),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.NewResponseSuccess(res, "Categories Retrieved"),
	)
}

func (h *CategoryHandler) GetCategoryTree(c *echo.Context) error {

	res, err := h.cs.GetTree(c.Request().Context())
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			response.NewResponseError(
				err.Error(),
				*customs.NewErrorValue("business_logic", err.Error()),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.NewResponseSuccess(res, "Category Tree Retrieved"),
	)
}

func (h *CategoryHandler) UpdateCategory(c *echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				"Invalid Category ID",
				*customs.NewErrorValue("validation", "id must be number"),
			),
		)
	}

	req := new(delivery.UpdateCategoryRequest)

	// bind
	if err := c.Bind(req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				"Validation Failed",
				customs.HandleBindError(err)...,
			),
		)
	}

	// validate
	if err := h.v.StructCtx(c.Request().Context(), req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				"Validation Failed",
				*customs.NewErrorValue("validation", err.Error()),
			),
		)
	}

	res, err := h.cs.Update(c.Request().Context(), uint(id), req)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				err.Error(),
				*customs.NewErrorValue("business_logic", err.Error()),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.NewResponseSuccess(res, "Category Updated"),
	)
}

func (h *CategoryHandler) DeleteCategory(c *echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				"Invalid Category ID",
				*customs.NewErrorValue("validation", "id must be number"),
			),
		)
	}

	if err := h.cs.Delete(c.Request().Context(), uint(id)); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			response.NewResponseError(
				err.Error(),
				*customs.NewErrorValue("business_logic", err.Error()),
			),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.NewResponseSuccess[*bool](nil, "Category Deleted"),
	)
}
