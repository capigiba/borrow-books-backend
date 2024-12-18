package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"borrow_book/internal/config"
	"borrow_book/internal/domain/model"
	"borrow_book/internal/initialize"
	"borrow_book/internal/repository"
	"borrow_book/pkg/logger"

	_ "github.com/lib/pq"
)

// BookWithAuthorName is a helper struct for migration purposes
type BookWithAuthorName struct {
	Book       model.Book
	AuthorName string
}

var log logger.Logger

func main() {
	// Define command-line flags
	authorFilePath := flag.String("authors", "", "Path to the authors data file (e.g., authors.txt)")
	bookFilePath := flag.String("books", "", "Path to the books data file (e.g., books.txt)")
	flag.Parse()

	// Create a context
	ctx := context.Background()

	// Validate flags
	if *authorFilePath == "" || *bookFilePath == "" {
		fmt.Println("Usage: go run main.go --authors=path/to/authors.txt --books=path/to/books.txt")
		os.Exit(1)
	}

	log = logger.NewLogger("migration v3")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("Config loading error: %v", err)
		os.Exit(1)
	}

	// Initialize databases
	db, err := initialize.InitDatabases(cfg, &log)
	if err != nil {
		log.Errorf("Database initialization error: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	authorRepo := repository.NewAuthorRepository(db)
	bookRepo := repository.NewBookRepository(db)

	// Begin the migration within a transaction
	tx, err := db.Beginx()
	if err != nil {
		log.Fatalf("Error beginning transaction: %v", err)
	}

	// Defer a rollback in case anything fails. If the transaction is committed, this will do nothing.
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Fatalf("Panic occurred: %v. Transaction rolled back.", p)
		} else if err != nil {
			tx.Rollback()
			log.Fatalf("Error occurred: %v. Transaction rolled back.", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Fatalf("Error committing transaction: %v", err)
			}
			log.Infof("Migration v3 completed successfully.")
		}
	}()

	// Resolve absolute paths
	authorAbsPath, err := filepath.Abs(*authorFilePath)
	if err != nil {
		log.Fatalf("Error resolving author file path: %v", err)
	}

	bookAbsPath, err := filepath.Abs(*bookFilePath)
	if err != nil {
		log.Fatalf("Error resolving book file path: %v", err)
	}

	// Step 1: Insert authors
	authors, err := readAuthors(authorAbsPath)
	if err != nil {
		log.Fatalf("Error reading authors from file: %v", err)
	}

	if len(authors) > 0 {
		log.Infof("Inserting %d authors.", len(authors))
		err = insertAuthors(ctx, authors, authorRepo)
		if err != nil {
			log.Fatalf("Error inserting authors: %v", err)
		}
		log.Infof("Inserted %d authors successfully.", len(authors))
	} else {
		log.Infof("No authors to insert.")
	}

	// Step 2: Insert books
	books, err := readBooks(bookAbsPath)
	if err != nil {
		log.Fatalf("Error reading books from file: %v", err)
	}

	if len(books) > 0 {
		log.Infof("Inserting %d books.", len(books))
		err = insertBooks(ctx, books, authorRepo, bookRepo)
		if err != nil {
			log.Fatalf("Error inserting books: %v", err)
		}
		log.Infof("Inserted %d books successfully.", len(books))
	} else {
		log.Infof("No books to insert.")
	}
}

// readAuthors reads the authors from the given file path.
func readAuthors(filePath string) ([]model.Author, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening author file: %w", err)
	}
	defer file.Close()

	var authors []model.Author
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format in author file at line %d: %s", lineNumber, line)
		}

		name := strings.TrimSpace(parts[1])
		if name == "" {
			return nil, fmt.Errorf("empty name in author file at line %d", lineNumber)
		}

		authors = append(authors, model.Author{
			Name: name,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning author file: %w", err)
	}

	return authors, nil
}

// readBooks reads the books from the given file path and associates them with their author names.
func readBooks(filePath string) ([]BookWithAuthorName, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening book file: %w", err)
	}
	defer file.Close()

	var books []BookWithAuthorName
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid format in book file at line %d: %s", lineNumber, line)
		}

		//idStr := strings.TrimSpace(parts[0])
		title := strings.TrimSpace(parts[1])
		authorName := strings.TrimSpace(parts[2])
		publishedYearStr := strings.TrimSpace(parts[3])

		if title == "" {
			return nil, fmt.Errorf("empty title in book file at line %d", lineNumber)
		}

		if authorName == "" {
			return nil, fmt.Errorf("empty author name in book file at line %d", lineNumber)
		}

		publishedAt, err := strconv.ParseInt(publishedYearStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid PublishedAt in book file at line %d: %v", lineNumber, err)
		}

		books = append(books, BookWithAuthorName{
			Book: model.Book{
				Title:       title,
				PublishedAt: publishedAt,
			},
			AuthorName: authorName,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning book file: %w", err)
	}

	return books, nil
}

// insertAuthors inserts a slice of Author into the database ensuring no duplicates.
func insertAuthors(ctx context.Context, authors []model.Author, repo repository.AuthorRepository) error {
	for _, author := range authors {
		// Check if the author exists by Name
		existingAuthor, err := repo.GetAuthorByName(ctx, author.Name)
		if err != nil {
			return fmt.Errorf("checking existence for author Name '%s': %w", author.Name, err)
		}
		if existingAuthor != nil {
			log.Infof("Author already exists (ID: %d, Name: %s). Skipping insertion.", existingAuthor.ID, existingAuthor.Name)
			continue
		}

		// Insert the author using repository
		newID, err := repo.CreateAuthor(ctx, author)
		if err != nil {
			return fmt.Errorf("inserting author Name '%s': %w", author.Name, err)
		}
		log.Infof("Inserted author (ID: %d, Name: %s).", newID, author.Name)
	}

	return nil
}

// insertBooks inserts a slice of BookWithAuthorName into the database ensuring no duplicates.
func insertBooks(ctx context.Context, books []BookWithAuthorName, authorRepo repository.AuthorRepository, bookRepo repository.BookRepository) error {
	for _, bookWithAuthor := range books {
		book := bookWithAuthor.Book
		authorName := bookWithAuthor.AuthorName

		// Resolve AuthorName to AuthorID
		author, err := authorRepo.GetAuthorByName(ctx, authorName)
		if err != nil {
			return fmt.Errorf("resolving author name '%s': %w", authorName, err)
		}
		if author == nil {
			return fmt.Errorf("author '%s' not found for book '%s'", authorName, book.Title)
		}
		book.AuthorID = author.ID

		// Check if the book exists by Title and AuthorID
		existingBook, err := bookRepo.GetBookByTitleAndAuthorID(ctx, book.Title, book.AuthorID)
		if err != nil {
			return fmt.Errorf("checking existence for book Title '%s' and AuthorID %d: %w", book.Title, book.AuthorID, err)
		}
		if existingBook != nil {
			log.Infof("Book already exists (ID: %d, Title: %s). Skipping insertion.", existingBook.ID, existingBook.Title)
			continue
		}

		// Insert the book using repository
		newID, err := bookRepo.CreateBook(ctx, book)
		if err != nil {
			return fmt.Errorf("inserting book Title '%s': %w", book.Title, err)
		}
		log.Infof("Inserted book (ID: %d, Title: %s).", newID, book.Title)
	}

	return nil
}
