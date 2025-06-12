package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StandardizedErrorResponse struct {
	Status    string            `json:"status"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	RequestID string            `json:"request_id"`
	Errors    []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
	Status    string      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
}

type AppError struct {
	HTTPCode         int
	Code             string
	Message          string
	Err              error
	ValidationErrors []ValidationError
}

func SendCreatedResponse(ctx *gin.Context, message string, requestID string, data interface{}) {
	sendSuccessResponse(ctx, http.StatusCreated, message, requestID, data)
}

func SendOKResponse(ctx *gin.Context, message string, requestID string, data interface{}) {
	sendSuccessResponse(ctx, http.StatusOK, message, requestID, data)
}

func sendSuccessResponse(ctx *gin.Context, statusCode int, message string, requestID string, data interface{}) {
	ctx.JSON(statusCode, SuccessResponse{
		Message:   message,
		RequestID: requestID,
		Status:    "SUCCESS",
		Data:      data,
	})
}

func (ae AppError) Error() string {
	if ae.Message != "" {
		return ae.Message
	}
	if ae.Err != nil {
		return ae.Err.Error()
	}
	return ""
}

func NewNotFoundError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusNotFound,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func NewBadRequestError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusBadRequest,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func NewInternalServerError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusInternalServerError,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func NewUnauthorizedError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusUnauthorized,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func NewForbiddenError(code string, message string, err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusForbidden,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

func HandleErrorResponse(ctx *gin.Context, err error, requestID string) {
	var response StandardizedErrorResponse
	response.RequestID = requestID
	response.Status = "ERROR"

	appErr, ok := err.(*AppError)
	if !ok {
		appErr = NewInternalServerError("INTERNAL_ERROR", "An unexpected error occurred", err)
	}

	response.Code = appErr.Code
	response.Message = appErr.Message

	if len(appErr.ValidationErrors) > 0 {
		response.Errors = appErr.ValidationErrors
	}

	ctx.JSON(appErr.HTTPCode, response)
}

func GetRequestID(ctx *gin.Context) string {
	requestID := ctx.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}
