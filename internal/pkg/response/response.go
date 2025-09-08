package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success returns a successful response
func Success(c echo.Context, data interface{}, message ...string) error {
	resp := Response{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		resp.Message = message[0]
	}

	return c.JSON(http.StatusOK, resp)
}

// Created returns a created response
func Created(c echo.Context, data interface{}, message ...string) error {
	resp := Response{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		resp.Message = message[0]
	} else {
		resp.Message = "Resource created successfully"
	}

	return c.JSON(http.StatusCreated, resp)
}

// BadRequest returns a bad request error response
func BadRequest(c echo.Context, message string, details ...interface{}) error {
	resp := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	}

	if len(details) > 0 {
		resp.Error.Details = details[0]
	}

	return c.JSON(http.StatusBadRequest, resp)
}

// Unauthorized returns an unauthorized error response
func Unauthorized(c echo.Context, message ...string) error {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}

	resp := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "UNAUTHORIZED",
			Message: msg,
		},
	}

	return c.JSON(http.StatusUnauthorized, resp)
}

// Forbidden returns a forbidden error response
func Forbidden(c echo.Context, message ...string) error {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}

	resp := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "FORBIDDEN",
			Message: msg,
		},
	}

	return c.JSON(http.StatusForbidden, resp)
}

// NotFound returns a not found error response
func NotFound(c echo.Context, message ...string) error {
	msg := "Resource not found"
	if len(message) > 0 {
		msg = message[0]
	}

	resp := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "NOT_FOUND",
			Message: msg,
		},
	}

	return c.JSON(http.StatusNotFound, resp)
}

// Conflict returns a conflict error response
func Conflict(c echo.Context, message string, details ...interface{}) error {
	resp := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "CONFLICT",
			Message: message,
		},
	}

	if len(details) > 0 {
		resp.Error.Details = details[0]
	}

	return c.JSON(http.StatusConflict, resp)
}

// InternalServerError returns an internal server error response
func InternalServerError(c echo.Context, message ...string) error {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}

	resp := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: msg,
		},
	}

	return c.JSON(http.StatusInternalServerError, resp)
}