package controllers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "mini-project-buku-sb-go-73-Agil/models"
    "mini-project-buku-sb-go-73-Agil/database"
)

// GetCategories - Get all categories
func GetCategories(c *gin.Context) {
    rows, err := database.DB.Query("SELECT id, name, created_at, created_by FROM categories ORDER BY id")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var cat models.Category
        if err := rows.Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.CreatedBy); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        categories = append(categories, cat)
    }

    c.JSON(http.StatusOK, categories)
}

// CreateCategory - Create new category
func CreateCategory(c *gin.Context) {
    var category models.Category
    if err := c.ShouldBindJSON(&category); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    username, _ := c.Get("username")
    
    query := `INSERT INTO categories (name, created_by) VALUES ($1, $2) RETURNING id, created_at`
    err := database.DB.QueryRow(query, category.Name, username).Scan(&category.ID, &category.CreatedAt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, category)
}

// GetCategoryByID - Get category by ID
func GetCategoryByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var category models.Category
    query := "SELECT id, name, created_at, created_by FROM categories WHERE id = $1"
    err = database.DB.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy)
    
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        return
    }

    c.JSON(http.StatusOK, category)
}

// DeleteCategory - Delete category by ID
func DeleteCategory(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    // Check if category exists
    var exists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)"
    database.DB.QueryRow(checkQuery, id).Scan(&exists)
    
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        return
    }

    // Check if category has books
    var hasBooks bool
    booksQuery := "SELECT EXISTS(SELECT 1 FROM books WHERE category_id = $1)"
    database.DB.QueryRow(booksQuery, id).Scan(&hasBooks)
    
    if hasBooks {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with existing books"})
        return
    }

    // Delete category
    _, err = database.DB.Exec("DELETE FROM categories WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// GetBooksByCategory - Get books by category ID
func GetBooksByCategory(c *gin.Context) {
    categoryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
        return
    }

    rows, err := database.DB.Query(`
        SELECT b.id, b.title, b.description, b.image_url, b.release_year, 
               b.price, b.total_page, b.thickness, b.category_id,
               b.created_at, b.created_by
        FROM books b
        WHERE b.category_id = $1
        ORDER BY b.id
    `, categoryID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var books []models.Book
    for rows.Next() {
        var book models.Book
        err := rows.Scan(
            &book.ID, &book.Title, &book.Description, &book.ImageURL,
            &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness,
            &book.CategoryID, &book.CreatedAt, &book.CreatedBy,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        books = append(books, book)
    }

    c.JSON(http.StatusOK, books)
}