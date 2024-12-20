package service

import (
	"borrow_book/internal/domain/model"
	"borrow_book/internal/infra/database/query"
	"borrow_book/internal/repository"
	"context"
	"fmt"
)

type BorrowService interface {
	ListBorrowLists(ctx context.Context, filters, sorts []string, fields string) ([]model.Borrow, error)
	GetBorrow(ctx context.Context, id int) (*model.Borrow, error)
	CreateBorrow(ctx context.Context, bookID int, userName string, borrowedAt int64) (*model.Borrow, error)
	UpdateBorrow(ctx context.Context, id int, bookID int, userName string, borrowedAt int64) (*model.Borrow, error)
	DeleteBorrow(ctx context.Context, id int) error
}

type borrowService struct {
	repo repository.BorrowRepository
}

func NewBorrowService(repo repository.BorrowRepository) BorrowService {
	return &borrowService{repo: repo}
}

func (s *borrowService) ListBorrowLists(ctx context.Context, filters, sorts []string, fields string) ([]model.Borrow, error) {
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

	return s.repo.GetAllBorrowLists(ctx, opts)
}

func (s *borrowService) GetBorrow(ctx context.Context, id int) (*model.Borrow, error) {
	return s.repo.GetBorrowByID(ctx, id)
}

func (s *borrowService) CreateBorrow(ctx context.Context, bookID int, userName string, borrowedAt int64) (*model.Borrow, error) {
	newBorrow := model.Borrow{
		BookID:     bookID,
		UserName:   userName,
		BorrowedAt: borrowedAt,
	}
	id, err := s.repo.CreateBorrow(ctx, newBorrow)
	if err != nil {
		return nil, err
	}
	newBorrow.ID = id
	return &newBorrow, nil
}

func (s *borrowService) UpdateBorrow(ctx context.Context, id int, bookID int, userName string, borrowedAt int64) (*model.Borrow, error) {
	b, err := s.GetBorrow(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("borrow not found")
	}

	b.BookID = bookID
	b.UserName = userName
	b.BorrowedAt = borrowedAt

	err = s.repo.UpdateBorrow(ctx, *b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *borrowService) DeleteBorrow(ctx context.Context, id int) error {
	b, err := s.GetBorrow(ctx, id)
	if err != nil {
		return err
	}
	if b == nil {
		return fmt.Errorf("borrow not found with id %d", id)
	}
	return s.repo.DeleteBorrow(ctx, id)
}
