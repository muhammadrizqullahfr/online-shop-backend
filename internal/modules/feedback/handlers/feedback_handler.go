package handlers

import (
	"net/http"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type FeedbackHandler struct {
	fs services.FeedbackService
	v  *validator.Validate
}

func NewFeedbackHandler(fs services.FeedbackService, v *validator.Validate) *FeedbackHandler {
	return &FeedbackHandler{fs: fs, v: v}
}

func (h *FeedbackHandler) Create(c *echo.Context) error {
	req := new(delivery.CreateFeedbackRequest)
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

	res, err := h.fs.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("business_logic", err.Error()),
		))
	}

	return c.JSON(http.StatusCreated, response.NewResponseSuccess(res, "Feedback Submitted"))
}
