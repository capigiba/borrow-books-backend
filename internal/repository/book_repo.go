package repository

import (
	"borrow_book/internal/domain/model"
	"borrow_book/internal/infra/database/query"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type BookRepository interface {
	GetAllBooks(ctx context.Context, opts query.QueryOptions) ([]model.Book, error)
	GetBookByID(ctx context.Context, id int) (*model.Book, error)
	CreateBook(ctx context.Context, b model.Book) (int, error)
	UpdateBook(ctx context.Context, b model.Book) error
	DeleteBook(ctx context.Context, id int) error
}

type bookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) GetAllBooks(ctx context.Context, opts query.QueryOptions) ([]model.Book, error) {
	q, args := query.BuildSelectQuery("books", opts)
	var books []model.Book
	err := r.db.SelectContext(ctx, &books, q, args...)
	return books, err
}

func (r *bookRepository) GetBookByID(ctx context.Context, id int) (*model.Book, error) {
	var book model.Book
	err := r.db.GetContext(ctx, &book, "SELECT id, title, author_id, published_at FROM books WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &book, err
}

func (r *bookRepository) CreateBook(ctx context.Context, b model.Book) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO books (title, author_id, published_at) VALUES ($1, $2, $3) RETURNING id",
		b.Title, b.AuthorID, b.PublishedAt,
	).Scan(&id)
	return id, err
}

func (r *bookRepository) UpdateBook(ctx context.Context, b model.Book) error {
	res, err := r.db.ExecContext(ctx,
		"UPDATE books SET title=$1, author_id=$2, published_at=$3 WHERE id=$4",
		b.Title, b.AuthorID, b.PublishedAt, b.ID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}

func (r *bookRepository) DeleteBook(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM books WHERE id=$1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no rows deleted")
	}
	return nil
}
