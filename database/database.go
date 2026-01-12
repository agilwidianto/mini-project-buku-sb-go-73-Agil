package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "golang.org/x/crypto/bcrypt"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() error {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using environment variables")
    }

    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_SSLMODE"),
    )

    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }

    err = DB.Ping()
    if err != nil {
        return err
    }

    log.Println("Connected to database successfully")
    
    // Create tables if not exists
    if err := createTables(); err != nil {
        return fmt.Errorf("failed to create tables: %v", err)
    }
    
    // Insert default admin user if not exists
    if err := createDefaultUser(); err != nil {
        return fmt.Errorf("failed to create default user: %v", err)
    }

    return nil
}

func createTables() error {
    tables := []string{
        // Users table
        `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(100) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            created_by VARCHAR(100),
            modified_at TIMESTAMP,
            modified_by VARCHAR(100)
        )`,
        
        // Categories table
        `CREATE TABLE IF NOT EXISTS categories (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            created_by VARCHAR(100),
            modified_at TIMESTAMP,
            modified_by VARCHAR(100)
        )`,
        
        // Books table with check constraint
        `CREATE TABLE IF NOT EXISTS books (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            image_url VARCHAR(500),
            release_year INTEGER NOT NULL,
            price INTEGER NOT NULL,
            total_page INTEGER NOT NULL,
            thickness VARCHAR(50),
            category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            created_by VARCHAR(100),
            modified_at TIMESTAMP,
            modified_by VARCHAR(100),
            CONSTRAINT chk_release_year CHECK (release_year >= 1980 AND release_year <= 2024)
        )`,
    }

    for _, table := range tables {
        _, err := DB.Exec(table)
        if err != nil {
            return err
        }
    }
    
    log.Println("Tables created/verified successfully")
    return nil
}

func createDefaultUser() error {
    // Check if admin user already exists
    var count int
    err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
    if err != nil {
        return err
    }

    if count == 0 {
        // Hash password: "password123" menggunakan bcrypt
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        
        _, err = DB.Exec(`
            INSERT INTO users (username, password, created_by) 
            VALUES ($1, $2, $3)
        `, "admin", string(hashedPassword), "system")
        
        if err != nil {
            return err
        }
        
        log.Println("Default admin user created: admin/password123")
    }
    
    return nil
}