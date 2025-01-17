package main

import (
	"Agency-Type-back/internal/app/handlers"
	"Agency-Type-back/internal/app/middleware"
	"Agency-Type-back/internal/database"
	"log"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	app := echo.New()

	app.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: "time = ${time_rfc3339}\nmethod = ${method}, uri = ${uri}, status = ${status}\nheader = ${header:token}, query = ${query:user_id}\nremote_ip = ${remote_ip}, host = ${host}\nuser_agent = ${user_agent}\nlatency_human = ${latency_human} \nbytes_in = ${bytes_in}, bytes_out = ${bytes_out}\n\n",
	}))
	app.Use(echoMiddleware.Recover())
	app.Use(echoMiddleware.CORS())

	app.POST("/register", handlers.RegisterHandler(db))
	app.POST("/login", handlers.LoginHandler(db))
	app.POST("/results", handlers.NewResultHandler(db), middleware.JWTMiddleware)
	app.GET("/results", handlers.UserResultsHandler(db))
	app.GET("/results/:id", handlers.ResultHandler(db))
	app.GET("/records/:gamemode", handlers.GlobalRecordsHandler(db))
	app.GET("/records", handlers.UserRecordsHandler(db))

	app.Logger.Fatal(app.Start(":8080"))
}
