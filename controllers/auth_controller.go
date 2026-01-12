package controllers

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    "mini-project-buku-sb-go-73-Agil/models"
    "mini-project-buku-sb-go-73-Agil/database"
    "os"
)

func Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Query user from database
    var user models.User
    query := "SELECT id, username, password FROM users WHERE username = $1"
    err := database.DB.QueryRow(query, req.Username).Scan(&user.ID, &user.Username, &user.Password)
    
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Check password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Create JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "user_id":  user.ID,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "your-secret-key-change-in-production"
    }

    tokenString, err := token.SignedString([]byte(jwtSecret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, models.LoginResponse{Token: tokenString})
}

// Register - Optional endpoint for user registration
func Register(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    // Insert user into database
    query := "INSERT INTO users (username, password, created_by) VALUES ($1, $2, $3) RETURNING id"
    var userID int
    err = database.DB.QueryRow(query, req.Username, string(hashedPassword), "system").Scan(&userID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "user_id": userID,
    })
}