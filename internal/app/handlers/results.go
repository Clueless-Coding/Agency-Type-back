package handlers

import (
	"Agency-Type-back/internal/app/models"
	"Agency-Type-back/internal/app/utils"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewResultHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get("user_id").(int)

		result := new(models.Result)
		if err := ctx.Bind(result); err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
		}

		_, err := db.Exec("INSERT INTO results (user_id,game_mode,duration,mistakes,accuracy,count_words,wpm,cpm) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", userID, result.GameMode, result.Duration, result.Mistakes, result.Accuracy, result.Words, result.WPM, result.CPM)
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to create result")
		}
		return ctx.JSON(http.StatusCreated, map[string]string{"message": "Result created successfully"})
	}
}

func UserResultsHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		userID := ctx.QueryParam("user_id")

		if userID == "" {
			return utils.BuildErrorResponse(ctx, http.StatusBadRequest, "UserID parameter is required")
		}

		rows, err := db.Query("SELECT * FROM results WHERE user_id = $1", userID)
		if err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch results")
		}
		defer rows.Close()

		var results []models.Result
		for rows.Next() {
			var result models.Result
			if err := rows.Scan(&result.ID, &result.UserID, &result.GameMode, &result.StartTime, &result.Duration, &result.Mistakes, &result.Accuracy, &result.Words, &result.WPM, &result.CPM); err != nil {
				fmt.Println(err)
				return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to scan results")
			}
			results = append(results, result)

		}

		if err := rows.Err(); err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Error while iterating over results")
		}

		return ctx.JSON(http.StatusOK, results)
	}
}

func ResultHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		resultID := ctx.Param("id")

		if resultID == "" {
			return utils.BuildErrorResponse(ctx, http.StatusBadRequest, "ResultID parameter is required")
		}

		row := db.QueryRow("SELECT * FROM results WHERE id = $1", resultID)

		var result models.Result
		if err := row.Scan(&result.ID, &result.UserID, &result.GameMode, &result.StartTime, &result.Duration, &result.Mistakes, &result.Accuracy, &result.Words, &result.WPM, &result.CPM); err != nil {
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to scan result")
		}

		return ctx.JSON(http.StatusOK, result)
	}
}
