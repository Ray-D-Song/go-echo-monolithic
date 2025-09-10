package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/response"
	"github.com/ray-d-song/go-echo-monolithic/internal/service"
	"github.com/ray-d-song/go-echo-monolithic/internal/types"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile retrieves current user profile
// @Summary		Get current user profile
// @Description	Get the profile of the currently authenticated user
// @Tags			users
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	response.Response	"Profile retrieved successfully"
// @Failure		401	{object}	response.Response	"Unauthorized"
// @Failure		404	{object}	response.Response	"User not found"
// @Failure		500	{object}	response.Response	"Internal server error"
// @Router			/users/profile [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	user, err := h.userService.GetByID(userID)
	if err != nil {
		if err == types.ErrUserNotFound {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to get user profile")
	}

	userResponse := h.userService.ToResponse(user)
	return response.Success(c, userResponse, "Profile retrieved successfully")
}

// GetUserByID retrieves user by ID
// @Summary		Get user by ID
// @Description	Get user information by user ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	response.Response	"User retrieved successfully"
// @Failure		400	{object}	response.Response	"Invalid user ID"
// @Failure		401	{object}	response.Response	"Unauthorized"
// @Failure		404	{object}	response.Response	"User not found"
// @Failure		500	{object}	response.Response	"Internal server error"
// @Router			/users/{id} [get]
func (h *UserHandler) GetUserByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		if err == types.ErrUserNotFound {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to get user")
	}

	userResponse := h.userService.ToResponse(user)
	return response.Success(c, userResponse, "User retrieved successfully")
}

// GetUserByUsername retrieves user by username
// @Summary		Get user by username
// @Description	Get user information by username
// @Tags			users
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			username	path		string	true	"Username"
// @Success		200			{object}	response.Response	"User retrieved successfully"
// @Failure		400			{object}	response.Response	"Username is required"
// @Failure		401			{object}	response.Response	"Unauthorized"
// @Failure		404			{object}	response.Response	"User not found"
// @Failure		500			{object}	response.Response	"Internal server error"
// @Router			/users/username/{username} [get]
func (h *UserHandler) GetUserByUsername(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return response.BadRequest(c, "Username is required")
	}

	user, err := h.userService.GetByUsername(username)
	if err != nil {
		if err == types.ErrUserNotFound {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to get user")
	}

	userResponse := h.userService.ToResponse(user)
	return response.Success(c, userResponse, "User retrieved successfully")
}

// UpdateProfile updates current user profile
// @Summary		Update user profile
// @Description	Update the profile of the currently authenticated user
// @Tags			users
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		types.UpdateUserRequest	true	"Update user request"
// @Success		200		{object}	response.Response			"Profile updated successfully"
// @Failure		400		{object}	response.Response			"Bad request"
// @Failure		401		{object}	response.Response			"Unauthorized"
// @Failure		404		{object}	response.Response			"User not found"
// @Failure		409		{object}	response.Response			"Email already in use"
// @Failure		500		{object}	response.Response			"Internal server error"
// @Router			/users/profile [put]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req types.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request data")
	}

	user, err := h.userService.Update(userID, &req)
	if err != nil {
		switch err {
		case types.ErrUserNotFound:
			return response.NotFound(c, "User not found")
		case types.ErrUserAlreadyExists:
			return response.Conflict(c, "Email already in use")
		default:
			return response.InternalServerError(c, "Failed to update profile")
		}
	}

	userResponse := h.userService.ToResponse(user)
	return response.Success(c, userResponse, "Profile updated successfully")
}

// ListUsers retrieves users with pagination
// @Summary		List users
// @Description	Retrieve a paginated list of users
// @Tags			users
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page	query		int	false	"Page number (default: 1)"
// @Param			limit	query		int	false	"Items per page (default: 20, max: 100)"
// @Success		200	{object}	response.Response	"Users retrieved successfully"
// @Failure		401	{object}	response.Response	"Unauthorized"
// @Failure		500	{object}	response.Response	"Internal server error"
// @Router			/users [get]
func (h *UserHandler) ListUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	users, err := h.userService.List(offset, limit)
	if err != nil {
		return response.InternalServerError(c, "Failed to list users")
	}

	total, err := h.userService.Count()
	if err != nil {
		return response.InternalServerError(c, "Failed to get user count")
	}

	userResponses := make([]*types.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = h.userService.ToResponse(user)
	}

	paginatedResponse := types.PaginatedResponse{
		Data:    userResponses,
		Total:   total,
		Page:    page,
		Limit:   limit,
		Pages:   (total + int64(limit) - 1) / int64(limit),
	}

	return response.Success(c, paginatedResponse, "Users retrieved successfully")
}

// DeleteUser soft deletes a user
// @Summary		Delete user
// @Description	Soft delete a user by ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	response.Response	"User deleted successfully"
// @Failure		400	{object}	response.Response	"Invalid user ID"
// @Failure		401	{object}	response.Response	"Unauthorized"
// @Failure		404	{object}	response.Response	"User not found"
// @Failure		500	{object}	response.Response	"Internal server error"
// @Router			/users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		if err == types.ErrUserNotFound {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to delete user")
	}

	return response.Success(c, nil, "User deleted successfully")
}

// RegisterRoutes registers user routes
func (h *UserHandler) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	users := e.Group("/api/users")
	
	// Protected routes
	users.Use(authMiddleware)
	users.GET("/profile", h.GetProfile)
	users.PUT("/profile", h.UpdateProfile)
	users.GET("/:id", h.GetUserByID)
	users.GET("/username/:username", h.GetUserByUsername)
	users.GET("", h.ListUsers)
	users.DELETE("/:id", h.DeleteUser)
}