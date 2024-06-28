package middleware

import (
	"Agency-Type-back/internal/app/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		tokenString := ctx.Request().Header.Get("token")

		if tokenString == "" {
			return utils.BuildErrorResponse(ctx, http.StatusUnauthorized, "Missing token")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что метод подписи тот, что мы ожидаем
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_TOKEN")), nil
		})

		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusUnauthorized, "Invalid token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Set("user_id", int(claims["user_id"].(float64)))
			return next(ctx)
		}

		return utils.BuildErrorResponse(ctx, http.StatusUnauthorized, "Invalid token")
	}
}
