package main

import (
	"log"
	"os"

	"account/internal/database"
	"account/internal/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	// Replace with your actual PostgreSQL connection string.
	connStr := "postgres://admin:password@localhost:5432/registry?sslmode=disable"
	database.InitDB(connStr)
	defer database.DB.Close()

	// Create tables if they don't exist
	database.CreateTables(database.DB)

	e := echo.New()

	// Register CRUD routes for accounts.
	e.POST("/accounts", handlers.CreateAccount)
	e.GET("/accounts", handlers.ListAccounts)
	e.GET("/accounts/:id", handlers.GetAccount)
	e.PUT("/accounts/:id", handlers.UpdateAccount)
	e.DELETE("/accounts/:id", handlers.DeleteAccount)

	// Get port from environment or default to 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(e.Start(":" + port))
}
