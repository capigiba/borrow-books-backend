package main

import (
	"borrow_book/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Attempt to connect with retries
	var pool *pgxpool.Pool
	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.New(context.Background(), cfg.Database.PostgresURL)
		if err != nil {
			log.Printf("Attempt %d: Unable to connect to PostgreSQL: %v", i+1, err)
		} else {
			err = pool.Ping(context.Background())
			if err == nil {
				break
			}
			log.Printf("Attempt %d: PostgreSQL ping failed: %v", i+1, err)
			pool.Close()
		}
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL after %d attempts: %v", maxRetries, err)
	}
	defer pool.Close()

	fmt.Println("Successfully connected to PostgreSQL!")

	// Create tables
	err = createTables(context.Background(), pool)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	fmt.Println("Database tables created successfully.")
}

// createTables executes SQL statements to create the required tables.
func createTables(ctx context.Context, pool *pgxpool.Pool) error {
	// SQL statements to create tables
	createAuthorsTable := `
    CREATE TABLE IF NOT EXISTS authors (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL
    );
    `

	createBooksTable := `
    CREATE TABLE IF NOT EXISTS books (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        author_id INTEGER NOT NULL,
        published_at BIGINT NOT NULL,
        FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE
    );
    `

	createBorrowsTable := `
    CREATE TABLE IF NOT EXISTS borrows (
        id SERIAL PRIMARY KEY,
        book_id INTEGER NOT NULL,
        borrowed_at BIGINT NOT NULL,
        FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
    );
    `

	// Execute SQL statements sequentially
	statements := []string{
		createAuthorsTable,
		createBooksTable,
		createBorrowsTable,
	}

	for _, stmt := range statements {
		_, err := pool.Exec(ctx, stmt)
		if err != nil {
			return fmt.Errorf("error executing statement: %v\nStatement: %s", err, stmt)
		}
	}

	return nil
}
