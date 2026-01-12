package handlers

import (
	"book-api/config"
	"book-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, title, description, image_url,
		       release_year, price, total_page,
		       thickness, category_id
		FROM books
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.ImageURL,
			&book.ReleaseYear,
			&book.Price,
			&book.TotalPage,
			&book.Thickness,
			&book.CategoryID,
		)
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// VALIDASI release_year
	if book.ReleaseYear < 1980 || book.ReleaseYear > 2024 {
		c.JSON(400, gin.H{"error": "release_year must be between 1980-2024"})
		return
	}

	// KONVERSI thickness
	if book.TotalPage > 100 {
		book.Thickness = "tebal"
	} else {
		book.Thickness = "tipis"
	}

	_, err := config.DB.Exec(`
		INSERT INTO books
		(title, description, image_url, release_year, price, total_page, thickness, category_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		book.Title,
		book.Description,
		book.ImageURL,
		book.ReleaseYear,
		book.Price,
		book.TotalPage,
		book.Thickness,
		book.CategoryID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "book created"})
}

func GetBooksByCategory(c *gin.Context) {
	id := c.Param("id")

	rows, err := config.DB.Query(`
		SELECT id, title, thickness
		FROM books
		WHERE category_id = $1
	`, id)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		rows.Scan(&book.ID, &book.Title, &book.Thickness)
		books = append(books, book)
	}

	c.JSON(200, books)
}
