package models

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse стандартная модель ошибки API
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Описание ошибки
	// example: something went wrong
	Error string `json:"error"`
}

// SuccessResponse стандартная модель успешного ответа
// swagger:model SuccessResponse
type SuccessResponse struct {
	// Сообщение об успехе
	// example: operation completed successfully
	Message string `json:"message"`
}

// NewErrorResponse создает новый ErrorResponse
func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{Error: msg}
}

// RenderError отправляет ошибку в стандартном формате
func RenderError(w http.ResponseWriter, r *http.Request, msg string, status int) {
	render.Status(r, status)
	render.JSON(w, r, NewErrorResponse(msg))
}
