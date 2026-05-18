package handlers

import (
	"net/http"
	"strconv"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type ArticleHandler struct {
	as services.ArticleService
	v  *validator.Validate
}

func NewArticleHandler(as services.ArticleService, v *validator.Validate) *ArticleHandler {
	return &ArticleHandler{as: as, v: v}
}

func (h *ArticleHandler) Create(c *echo.Context) error {
	req := new(delivery.CreateArticleRequest)
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

	res, err := h.as.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusCreated, response.NewResponseSuccess(res, "Article Created"))
}

func (h *ArticleHandler) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Article ID",
			*customs.NewErrorValue("validation", "id is required"),
		))
	}

	req := new(delivery.UpdateArticleRequest)
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

	res, err := h.as.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Article Updated"))
}

func (h *ArticleHandler) GetByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Article ID",
			*customs.NewErrorValue("validation", "id is required"),
		))
	}

	res, err := h.as.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Article Found"))
}

func (h *ArticleHandler) GetAll(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	category := c.QueryParam("category")

	res, err := h.as.GetAll(c.Request().Context(), page, limit, search, category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "Articles Retrieved"))
}

func (h *ArticleHandler) Delete(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Article ID",
			*customs.NewErrorValue("validation", "id is required"),
		))
	}

	if err := h.as.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess[*bool](nil, "Article Deleted"))
}
