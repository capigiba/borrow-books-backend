package service

import (
	"borrow_book/internal/domain/model"
	"borrow_book/internal/infra/database/query"
	"borrow_book/internal/repository"
	"context"
	"fmt"
)

// AuthorService defines the interface for author-related operations
type AuthorService interface {
	ListAuthors(ctx context.Context, filters, sorts []string, fields string) ([]model.Author, error)
	GetAuthor(ctx context.Context, id int) (*model.Author, error)
	CreateAuthor(ctx context.Context, name string) (*model.Author, error)
	UpdateAuthor(ctx context.Context, id int, name string) (*model.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
}

// authorService is the concrete implementation of AuthorService
type authorService struct {
	repo repository.AuthorRepository
}

// NewAuthorService creates a new instance of AuthorService
func NewAuthorService(repo repository.AuthorRepository) AuthorService {
	return &authorService{repo: repo}
}

func (s *authorService) ListAuthors(ctx context.Context, filters, sorts []string, fields string) ([]model.Author, error) {
	f, err := query.ParseFilters(filters)
	if err != nil {
		return nil, err
	}
	srts, err := query.ParseSorts(sorts)
	if err != nil {
		return nil, err
	}
	fs := query.ParseFields(fields)

	opts := query.QueryOptions{
		Filters: f,
		Sorts:   srts,
		Fields:  fs,
	}
	return s.repo.GetAllAuthors(ctx, opts)
}

func (s *authorService) GetAuthor(ctx context.Context, id int) (*model.Author, error) {
	return s.repo.GetAuthorByID(ctx, id)
}

func (s *authorService) CreateAuthor(ctx context.Context, name string) (*model.Author, error) {
	newAuthor := model.Author{
		Name: name,
	}
	id, err := s.repo.CreateAuthor(ctx, newAuthor)
	if err != nil {
		return nil, err
	}
	newAuthor.ID = id
	return &newAuthor, nil
}

func (s *authorService) UpdateAuthor(ctx context.Context, id int, name string) (*model.Author, error) {
	// Optionally get existing author first for validation
	a, err := s.GetAuthor(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("author not found")
	}

	a.Name = name

	err = s.repo.UpdateAuthor(ctx, *a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *authorService) DeleteAuthor(ctx context.Context, id int) error {
	// Optionally check if it exists
	return s.repo.DeleteAuthor(ctx, id)
}
