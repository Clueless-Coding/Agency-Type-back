package handlers

import (
	"Agency-Type-back/internal/app/models"
	"Agency-Type-back/internal/app/utils"
	"Agency-Type-back/pkg/auth"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := new(models.User)
		if err := ctx.Bind(user); err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to hash password")
		}

		var userID int
		err = db.QueryRow("INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id", user.Login, hashedPassword).Scan(&userID)
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to register user")
		}

		exitID := userID

		token, err := auth.GenerateToken(userID)
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate token")
		}

		_, err = db.Exec("UPDATE users SET token = $1 WHERE id = $2", token, userID)
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to update user token")
		}

		return ctx.JSON(http.StatusCreated, map[string]interface{}{"message": "User registered successfully", "token": token, "user_id": exitID})
	}
}

func LoginHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := new(models.User)
		if err := ctx.Bind(user); err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
		}

		var dbUser models.User
		var userID int
		err := db.QueryRow("SELECT id, login, password_hash FROM users WHERE login = $1", user.Login).Scan(&userID, &dbUser.Login, &dbUser.Password)
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusUnauthorized, "User not found")
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusUnauthorized, "Invalid login credentials")
		}

		token, err := auth.GenerateToken(userID)
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate token")
		}

		_, err = db.Exec("UPDATE users SET token = $1 WHERE id = $2", token, userID)
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to update user token")
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Login successful", "token": token, "user_id": userID})
	}
}
