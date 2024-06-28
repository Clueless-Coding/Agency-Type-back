package utils

import (
	"github.com/labstack/echo/v4"
)

func BuildErrorResponse(ctx echo.Context, status int, message string) error {
	return ctx.JSON(status, map[string]string{"error": message})
}
