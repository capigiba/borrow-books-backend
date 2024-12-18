package service

import (
	"borrow_book/internal/domain/model"
	"borrow_book/internal/infra/database/query"
	"borrow_book/internal/repository"
	"context"
	"fmt"
)

type BookService interface {
	ListBooks(ctx context.Context, filters, sorts []string, fields string) ([]model.Book, error)
	GetBook(ctx context.Context, id int) (*model.Book, error)
	CreateBook(ctx context.Context, title string, authorID int, publishedAt int64) (*model.Book, error)
	UpdateBook(ctx context.Context, id int, title string, authorID int, publishedAt int64) (*model.Book, error)
	DeleteBook(ctx context.Context, id int) error
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) ListBooks(ctx context.Context, filters, sorts []string, fields string) ([]model.Book, error) {
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
	return s.repo.GetAllBooks(ctx, opts)
}

func (s *bookService) GetBook(ctx context.Context, id int) (*model.Book, error) {
	return s.repo.GetBookByID(ctx, id)
}

func (s *bookService) CreateBook(ctx context.Context, title string, authorID int, publishedAt int64) (*model.Book, error) {
	newBook := model.Book{
		Title:       title,
		AuthorID:    authorID,
		PublishedAt: publishedAt,
	}
	id, err := s.repo.CreateBook(ctx, newBook)
	if err != nil {
		return nil, err
	}
	newBook.ID = id
	return &newBook, nil
}

func (s *bookService) UpdateBook(ctx context.Context, id int, title string, authorID int, publishedAt int64) (*model.Book, error) {
	// Optionally get existing book first for validation
	b, err := s.GetBook(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("book not found")
	}

	b.Title = title
	b.AuthorID = authorID
	b.PublishedAt = publishedAt

	err = s.repo.UpdateBook(ctx, *b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *bookService) DeleteBook(ctx context.Context, id int) error {
	// Optionally check if it exists
	return s.repo.DeleteBook(ctx, id)
}
