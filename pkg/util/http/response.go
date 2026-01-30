package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// BaseResponse represents base http response
type BaseResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteOkResponse writes 200 response using echo.
func WriteOkResponse(ctx echo.Context, data interface{}, message string) error {
	resp := BaseResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
	return ctx.JSON(http.StatusOK, resp)
}

// WriteErrorResponse writes error response
func WriteErrorResponse(ctx echo.Context, statusCode int, message string) error {
	resp := BaseResponse{
		Status:  statusCode,
		Message: message,
	}
	return ctx.JSON(statusCode, resp)
}
