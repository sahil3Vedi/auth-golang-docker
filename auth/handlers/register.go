package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Types

type UserRegister struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register Handler Function

func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user UserRegister
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if username and email are unique
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2", user.Username, user.Email).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
			return
		}

		// Insert the user into the database
		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1,$2,$3)", user.Username, user.Email, (hashedPassword))
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register the user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}
