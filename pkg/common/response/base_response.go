package response

import "github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"

type ApiResponse[Data any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`

	Error any `json:"error,omitempty"`
}

func NewResponseSuccess[Data any](data Data, message string) *ApiResponse[Data] {
	return &ApiResponse[Data]{
		Code:    "SUCCESS",
		Message: message,
		Data:    data,
	}
}

func NewResponseError(message string, err ...customs.ErrorValue) *ApiResponse[*any] {
	return &ApiResponse[*any]{
		Code:    "ERROR",
		Message: message,
		Error:   customs.NewErrorValues(err...),
	}
}
