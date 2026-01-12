package controllers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "mini-project-buku-sb-go-73-Agil/models"
    "mini-project-buku-sb-go-73-Agil/database"
)

// GetBooks - Get all books
func GetBooks(c *gin.Context) {
    rows, err := database.DB.Query(`
        SELECT b.id, b.title, b.description, b.image_url, b.release_year, 
               b.price, b.total_page, b.thickness, b.category_id,
               b.created_at, b.created_by,
               c.id as category_id, c.name as category_name
        FROM books b
        LEFT JOIN categories c ON b.category_id = c.id
        ORDER BY b.id
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var books []models.Book
    for rows.Next() {
        var book models.Book
        var categoryID *int
        var categoryName *string
        
        err := rows.Scan(
            &book.ID, &book.Title, &book.Description, &book.ImageURL,
            &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness,
            &book.CategoryID, &book.CreatedAt, &book.CreatedBy,
            &categoryID, &categoryName,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if categoryID != nil && categoryName != nil {
            book.Category = &models.Category{
                ID:   *categoryID,
                Name: *categoryName,
            }
        }
        books = append(books, book)
    }

    c.JSON(http.StatusOK, books)
}

// CreateBook - Create new book
func CreateBook(c *gin.Context) {
    var req models.BookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Calculate thickness based on total page
    thickness := "tipis"
    if req.TotalPage > 100 {
        thickness = "tebal"
    }

    username, _ := c.Get("username")
    
    query := `
        INSERT INTO books (title, description, image_url, release_year, 
                          price, total_page, thickness, category_id, created_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at
    `
    
    var book models.Book
    book.Title = req.Title
    book.Description = req.Description
    book.ImageURL = req.ImageURL
    book.ReleaseYear = req.ReleaseYear
    book.Price = req.Price
    book.TotalPage = req.TotalPage
    book.Thickness = thickness
    book.CategoryID = req.CategoryID
    book.CreatedBy = username.(string)

    err := database.DB.QueryRow(query,
        book.Title, book.Description, book.ImageURL, book.ReleaseYear,
        book.Price, book.TotalPage, book.Thickness, book.CategoryID, book.CreatedBy,
    ).Scan(&book.ID, &book.CreatedAt)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, book)
}

// GetBookByID - Get book by ID
func GetBookByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var book models.Book
    var categoryID *int
    var categoryName *string
    
    query := `
        SELECT b.id, b.title, b.description, b.image_url, b.release_year, 
               b.price, b.total_page, b.thickness, b.category_id,
               b.created_at, b.created_by,
               c.id as category_id, c.name as category_name
        FROM books b
        LEFT JOIN categories c ON b.category_id = c.id
        WHERE b.id = $1
    `
    
    err = database.DB.QueryRow(query, id).Scan(
        &book.ID, &book.Title, &book.Description, &book.ImageURL,
        &book.ReleaseYear, &book.Price, &book.TotalPage, &book.Thickness,
        &book.CategoryID, &book.CreatedAt, &book.CreatedBy,
        &categoryID, &categoryName,
    )
    
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    if categoryID != nil && categoryName != nil {
        book.Category = &models.Category{
            ID:   *categoryID,
            Name: *categoryName,
        }
    }

    c.JSON(http.StatusOK, book)
}

// UpdateBook - Update book by ID
func UpdateBook(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    // Check if book exists
    var exists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)"
    database.DB.QueryRow(checkQuery, id).Scan(&exists)
    
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    var req models.BookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Calculate thickness based on total page
    thickness := "tipis"
    if req.TotalPage > 100 {
        thickness = "tebal"
    }

    username, _ := c.Get("username")
    
    query := `
        UPDATE books 
        SET title = $1, description = $2, image_url = $3, release_year = $4,
            price = $5, total_page = $6, thickness = $7, category_id = $8,
            modified_at = CURRENT_TIMESTAMP, modified_by = $9
        WHERE id = $10
        RETURNING id, created_at, modified_at
    `
    
    var book models.Book
    err = database.DB.QueryRow(query,
        req.Title, req.Description, req.ImageURL, req.ReleaseYear,
        req.Price, req.TotalPage, thickness, req.CategoryID,
        username, id,
    ).Scan(&book.ID, &book.CreatedAt, &book.ModifiedAt)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    book.Title = req.Title
    book.Description = req.Description
    book.ImageURL = req.ImageURL
    book.ReleaseYear = req.ReleaseYear
    book.Price = req.Price
    book.TotalPage = req.TotalPage
    book.Thickness = thickness
    book.CategoryID = req.CategoryID
    book.ModifiedBy = username.(string)

    c.JSON(http.StatusOK, book)
}

// DeleteBook - Delete book by ID
func DeleteBook(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    // Check if book exists
    var exists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)"
    database.DB.QueryRow(checkQuery, id).Scan(&exists)
    
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    _, err = database.DB.Exec("DELETE FROM books WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}