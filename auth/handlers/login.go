package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Helper Funcs

func authenticate(username, password string, db *sql.DB) (*UserLogin, error) {
	var user UserLogin
	query := "SELECT username,password FROM users WHERE username = $1"
	err := db.QueryRow(query, username).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Authentication Failed - Username Doesnt Exist")
		}
		return nil, err
	}

	fmt.Println("DB User:", user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("Authentication Failed - Incorrect Password")
	}

	return &user, nil
}

func generateJWTToken(user *UserLogin, secretKey string) (string, error) {
	fmt.Println("JWT User:", user)
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Login Handler Function

func LoginHandler(db *sql.DB, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user UserLogin
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		// authenticates password
		authenticatedUser, err := authenticate(user.Username, user.Password, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		}
		fmt.Println("Authenticated user", authenticatedUser)
		// if authenticated, then a JWT token with an expiry time is created
		token, err := generateJWTToken(authenticatedUser, jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"token": token})
	}
}
