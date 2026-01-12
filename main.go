package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "mini-project-buku-sb-go-73-Agil/controllers"
    "mini-project-buku-sb-go-73-Agil/middleware"
    "mini-project-buku-sb-go-73-Agil/database"
)

func main() {
    // Initialize database (auto-creates tables)
    if err := database.ConnectDB(); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize router
    r := gin.Default()

    // Public routes
    r.POST("/api/users/login", controllers.Login)
    r.POST("/api/users/register", controllers.Register) // Optional

    // Protected routes with JWT middleware
    api := r.Group("/api")
    api.Use(middleware.JWTAuthMiddleware())
    {
        // Categories routes
        categories := api.Group("/categories")
        {
            categories.GET("", controllers.GetCategories)
            categories.POST("", controllers.CreateCategory)
            categories.GET("/:id", controllers.GetCategoryByID)
            categories.DELETE("/:id", controllers.DeleteCategory)
            categories.GET("/:id/books", controllers.GetBooksByCategory)
        }

        // Books routes
        books := api.Group("/books")
        {
            books.GET("", controllers.GetBooks)
            books.POST("", controllers.CreateBook)
            books.GET("/:id", controllers.GetBookByID)
            books.PUT("/:id", controllers.UpdateBook)
            books.DELETE("/:id", controllers.DeleteBook)
        }
    }

    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "OK",
            "database": "connected",
        })
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = os.Getenv("SERVER_PORT")
    }
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s\n", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}