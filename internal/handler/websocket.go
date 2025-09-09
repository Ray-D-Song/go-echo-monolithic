package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/logger"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/response"
	"github.com/ray-d-song/go-echo-monolithic/internal/service"
	"go.uber.org/zap"
)

// WebSocketHandler handles WebSocket HTTP requests
type WebSocketHandler struct {
	wsService *service.WebSocketService
	logger    *logger.Logger
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(wsService *service.WebSocketService, logger *logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		wsService: wsService,
		logger:    logger,
	}
}

// HandleWebSocket upgrades HTTP connection to WebSocket
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	if err := h.wsService.UpgradeConnection(c, userID); err != nil {
		h.logger.Error("Failed to upgrade WebSocket connection", zap.Error(err))
		return response.InternalServerError(c, "Failed to upgrade connection")
	}

	return nil
}

// GetConnectedUsers returns the number of connected users
func (h *WebSocketHandler) GetConnectedUsers(c echo.Context) error {
	count := h.wsService.GetConnectedUsers()

	result := map[string]any{
		"connected_users": count,
	}

	return response.Success(c, result, "Connected users count retrieved")
}

// RegisterRoutes registers WebSocket routes
func (h *WebSocketHandler) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	ws := e.Group("/api/ws")

	// WebSocket upgrade endpoint - requires authentication
	ws.Use(authMiddleware)
	ws.GET("/connect", h.HandleWebSocket)
	ws.GET("/stats", h.GetConnectedUsers)
}
