package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ray-d-song/go-echo-monolithic/internal/handler/wrapper"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/response"
	"github.com/ray-d-song/go-echo-monolithic/internal/repository"
)

// ConfigHandler handles system configuration HTTP requests
type ConfigHandler struct {
	kvRepo *repository.KVRepository
}

// NewConfigHandler creates a new config handler
func NewConfigHandler(kvRepo *repository.KVRepository) *ConfigHandler {
	return &ConfigHandler{
		kvRepo: kvRepo,
	}
}

// ToggleUserRegistration toggles the user registration setting
// @Summary		Toggle user registration
// @Description	Toggle whether new user registration is allowed
// @Tags			config
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	response.Response	"Registration setting toggled successfully"
// @Failure		401	{object}	response.Response	"Unauthorized"
// @Failure		500	{object}	response.Response	"Internal server error"
// @Router			/api/config/registration/toggle [post]
func (h *ConfigHandler) ToggleUserRegistration(c echo.Context) error {
	key := "system:allow_register"

	// Get current value
	currentValue, err := h.kvRepo.Get(key)
	if err != nil {
		return response.InternalServerError(c, "Failed to get registration setting")
	}

	// Convert to boolean (default to false if empty)
	currentBool := false
	if currentValue != "" {
		currentBool, err = strconv.ParseBool(currentValue)
		if err != nil {
			return response.InternalServerError(c, "Invalid registration setting format")
		}
	}

	// Toggle the value
	newValue := !currentBool
	err = h.kvRepo.Set(key, strconv.FormatBool(newValue))
	if err != nil {
		return response.InternalServerError(c, "Failed to update registration setting")
	}

	return response.Success(c, map[string]interface{}{
		"allow_register": newValue,
		"message":        "Registration setting updated successfully",
	}, "Registration setting toggled successfully")
}

// RegisterRoutes registers config routes
func (h *ConfigHandler) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	config := e.Group("/api/config")

	config.Use(authMiddleware)
	config.POST("/registration/toggle", wrapper.AdminWrapper(h.ToggleUserRegistration))
}
