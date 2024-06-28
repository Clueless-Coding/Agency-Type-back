package handlers

import (
	"Agency-Type-back/internal/app/models"
	"Agency-Type-back/internal/app/utils"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GlobalRecordsHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		gamemode := ctx.Param("gamemode")

		if gamemode == "" {
			return utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Mode parameter is required")
		}

		rows, err := db.Query("SELECT * FROM results WHERE game_mode = $1 ORDER BY duration", gamemode)
		if err != nil {
			fmt.Println(err)
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch results")
		}
		defer rows.Close()

		var results []models.Result
		for rows.Next() {
			var result models.Result
			if err := rows.Scan(&result.ID, &result.UserID, &result.GameMode, &result.StartTime, &result.Duration, &result.Misstakes, &result.Accuracy, &result.Words); err != nil {
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

func UserRecordsHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		userID := ctx.QueryParam("user_id")

		if userID == "" {
			return utils.BuildErrorResponse(ctx, http.StatusBadRequest, "UserID parameter is required")
		}

		rows, err := db.Query(
			`WITH record_results AS (
            SELECT *, ROW_NUMBER() OVER (PARTITION BY game_mode ORDER BY duration ASC) AS rn
            FROM results WHERE user_id = $1)
            SELECT id, user_id, game_mode, start_time, duration, misstakes, accuracy, count_words FROM record_results WHERE rn = 1;`,
			userID)

		if err != nil {
			fmt.Println(err)
			return utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch results")
		}
		defer rows.Close()

		var results []models.Result
		for rows.Next() {
			var result models.Result
			if err := rows.Scan(&result.ID, &result.UserID, &result.GameMode, &result.StartTime, &result.Duration, &result.Misstakes, &result.Accuracy, &result.Words); err != nil {
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
