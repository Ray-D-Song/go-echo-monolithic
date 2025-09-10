package wrapper

import (
	"github.com/labstack/echo/v4"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/response"
)

// The handler wrapped by AdminWrapper is accessible only by admin users
func AdminWrapper(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("role")
		if role != "admin" {
			return response.Forbidden(c, "Accessible only by administrators")
		}
		return h(c)
	}
}
