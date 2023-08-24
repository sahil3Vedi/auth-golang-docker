package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"auth/handlers"
)

// Main Func

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	// Connect to PostgreSQL database
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbInfo := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create a new Gin router
	router := gin.Default()

	// Handle the register user API endpoint
	router.POST("/register", handlers.RegisterHandler(db))
	router.POST("/login", handlers.LoginHandler(db, jwtSecret))

	// Run the HTTP server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start the server", err)
	}
}
