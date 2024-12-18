package repository

import (
	"borrow_book/internal/domain/model"
	"borrow_book/internal/infra/database/query"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// AuthorRepository defines the interface for author-related data operations
type AuthorRepository interface {
	GetAllAuthors(ctx context.Context, opts query.QueryOptions) ([]model.Author, error)
	GetAuthorByID(ctx context.Context, id int) (*model.Author, error)
	GetAuthorByName(ctx context.Context, name string) (*model.Author, error)
	CreateAuthor(ctx context.Context, a model.Author) (int, error)
	UpdateAuthor(ctx context.Context, a model.Author) error
	DeleteAuthor(ctx context.Context, id int) error
}

// authorRepository is the concrete implementation of AuthorRepository
type authorRepository struct {
	db *sqlx.DB
}

// NewAuthorRepository creates a new instance of AuthorRepository
func NewAuthorRepository(db *sqlx.DB) AuthorRepository {
	return &authorRepository{db: db}
}

func (r *authorRepository) GetAllAuthors(ctx context.Context, opts query.QueryOptions) ([]model.Author, error) {
	q, args := query.BuildSelectQuery("authors", opts)
	var authors []model.Author
	err := r.db.SelectContext(ctx, &authors, q, args...)
	return authors, err
}

func (r *authorRepository) GetAuthorByID(ctx context.Context, id int) (*model.Author, error) {
	var author model.Author
	err := r.db.GetContext(ctx, &author, "SELECT id, name FROM authors WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &author, err
}

func (r *authorRepository) GetAuthorByName(ctx context.Context, name string) (*model.Author, error) {
	var author model.Author
	err := r.db.GetContext(ctx, &author, "SELECT id, name FROM authors WHERE name=$1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &author, err
}

func (r *authorRepository) CreateAuthor(ctx context.Context, a model.Author) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO authors (name) VALUES ($1) RETURNING id",
		a.Name,
	).Scan(&id)
	return id, err
}

func (r *authorRepository) UpdateAuthor(ctx context.Context, a model.Author) error {
	res, err := r.db.ExecContext(ctx,
		"UPDATE authors SET name=$1 WHERE id=$2",
		a.Name, a.ID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}

func (r *authorRepository) DeleteAuthor(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM authors WHERE id=$1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no rows deleted")
	}
	return nil
}
