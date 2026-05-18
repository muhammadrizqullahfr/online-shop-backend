package handlers

import (
	"net/http"

	"github.com/RakaMurdiarta/online-shop-system/internal/modules/upload/services"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/labstack/echo/v5"
)

type UploadHandler struct {
	us services.UploadService
}

func NewUploadHandler(us services.UploadService) *UploadHandler {
	return &UploadHandler{us: us}
}

func (h *UploadHandler) UploadImage(c *echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Invalid Request",
			*customs.NewErrorValue("body", "file field is required"),
		))
	}

	res, err := h.us.UploadImage(c.Request().Context(), file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			err.Error(),
			*customs.NewErrorValue("upload", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, response.NewResponseSuccess(res, "File uploaded successfully"))
}
