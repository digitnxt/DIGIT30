package main

import (
	"log"
	"os"

	"account/config"
	"account/controller"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load environment variables from .env file if available.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize the database.
	config.DatabaseInit()

	// Create an Echo instance.
	e := echo.New()

	// Define account routes.
	accountRoute := e.Group("/account")
	accountRoute.POST("/", controller.CreateAccount)
	accountRoute.GET("/:id", controller.GetAccount)
	accountRoute.PUT("/:id", controller.UpdateAccount)
	accountRoute.DELETE("/:id", controller.DeleteAccount)

	// Use PORT environment variable or default to 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server.
	log.Fatal(e.Start(":" + port))
}
