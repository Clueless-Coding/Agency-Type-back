package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Result struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	GameMode  string    `json:"gamemode"`
	StartTime time.Time `json:"start_time"`
	Duration  time.Time `json:"duration"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	//cors
	app.Use(middleware.CORS())
	//
	app.POST("/register", registerHandler(db))
	app.POST("/login", loginHandler(db))
	app.POST("/results", newResultHandler(db))              //new Result
	app.GET("/results", userResultsHandler(db))             //all Results
	app.GET("/results/:id", ResultHandler(db))              //about one Result
	app.GET("/records/:gamemode", globalRecordsHandler(db)) //global records per mode
	app.GET("/records", userRecordsHandler(db))             //users records

	app.Logger.Fatal(app.Start(":8080"))

}

func builderErrorRersponse(ctx echo.Context, status int, str string) error {
	return ctx.JSON(status, map[string]string{"error": str})
}

func newResultHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		result := new(Result)
		if err := ctx.Bind(result); err != nil {
			return builderErrorRersponse(ctx, http.StatusBadRequest, "Invalid request payload")
		}

		_, err := db.Exec("INSERT INTO results (user_id, duration) VALUES ($1, $2)", result.UserID, result.Duration)
		if err != nil {
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to create result")
		}
		return ctx.JSON(http.StatusCreated, map[string]string{"message": "Result created successfully"})
	}
}

func globalRecordsHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		gamemode := ctx.Param("gamemode")

		if gamemode == "" {
			return builderErrorRersponse(ctx, http.StatusBadRequest, "Mode parameter is required")
		}

		rows, err := db.Query("SELECT * FROM results WHERE game_mode = $1 ORDER BY duration", gamemode)
		if err != nil {
			fmt.Println(err)
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to fetch results")
		}
		defer rows.Close()

		var results []Result
		for rows.Next() {
			var result Result
			if err := rows.Scan(&result.ID, &result.UserID, &result.GameMode, &result.StartTime, &result.Duration); err != nil {
				fmt.Println(err)
				return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to scan results")
			}
			results = append(results, result)

		}

		if err := rows.Err(); err != nil {
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Error while iterating over results")
		}

		return ctx.JSON(http.StatusOK, results)
	}
}

func userRecordsHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		userID := ctx.QueryParam("user_id")

		if userID == "" {
			return builderErrorRersponse(ctx, http.StatusBadRequest, "UserID parameter is required")
		}

		rows, err := db.Query(
			`WITH record_results AS (
            SELECT *, ROW_NUMBER() OVER (PARTITION BY game_mode ORDER BY duration ASC) AS rn
            FROM results WHERE user_id = $1)
            SELECT id, user_id, game_mode, start_time, duration FROM record_results WHERE rn = 1;`,
			userID)

		if err != nil {
			fmt.Println(err)
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to fetch results")
		}
		defer rows.Close()

		var results []Result
		for rows.Next() {
			var result Result
			if err := rows.Scan(&result.ID, &result.UserID, &result.GameMode, &result.StartTime, &result.Duration); err != nil {
				fmt.Println(err)
				return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to scan results")
			}
			results = append(results, result)

		}

		if err := rows.Err(); err != nil {
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Error while iterating over results")
		}

		return ctx.JSON(http.StatusOK, results)
	}
}

func userResultsHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		userID := ctx.QueryParam("user_id")

		if userID == "" {
			return builderErrorRersponse(ctx, http.StatusBadRequest, "UserID parameter is required")
		}

		rows, err := db.Query("SELECT * FROM results WHERE user_id = $1", userID)
		if err != nil {
			fmt.Println(err)
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to fetch results")
		}
		defer rows.Close()

		var results []Result
		for rows.Next() {
			var result Result
			if err := rows.Scan(&result.ID, &result.UserID, &result.GameMode, &result.StartTime, &result.Duration); err != nil {
				fmt.Println(err)
				return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to scan results")
			}
			results = append(results, result)

		}

		if err := rows.Err(); err != nil {
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Error while iterating over results")
		}

		return ctx.JSON(http.StatusOK, results)
	}
}

func ResultHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		resultID := ctx.Param("id")

		if resultID == "" {
			return builderErrorRersponse(ctx, http.StatusBadRequest, "ResultID parameter is required")
		}

		row := db.QueryRow("SELECT * FROM results WHERE id = $1", resultID)

		var result Result
		if err := row.Scan(&result.ID, &result.UserID, &result.GameMode, &result.StartTime, &result.Duration); err != nil {
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to scan result")
		}

		return ctx.JSON(http.StatusOK, result)
	}
}

func registerHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		user := new(User)

		if err := ctx.Bind(user); err != nil {
			return builderErrorRersponse(ctx, http.StatusBadRequest, "Invalid request payload")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to hash password")
		}

		_, err = db.Exec("INSERT INTO users (login, password_hash) VALUES ($1, $2)", user.Login, hashedPassword)
		if err != nil {
			return builderErrorRersponse(ctx, http.StatusInternalServerError, "Failed to register user")
		}

		return ctx.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
	}
}

func loginHandler(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		user := new(User)
		if err := ctx.Bind(user); err != nil {
			return builderErrorRersponse(ctx, http.StatusBadRequest, "Invalid request payload")
		}

		var dbUser User
		err := db.QueryRow("SELECT login, password_hash FROM users WHERE login = $1", user.Login).Scan(&dbUser.Login, &dbUser.Password)
		if err != nil {
			return builderErrorRersponse(ctx, http.StatusUnauthorized, "User not found")
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if err != nil {
			return builderErrorRersponse(ctx, http.StatusUnauthorized, "Invalid login credentials")
		}

		return ctx.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
	}
}
