package service

import (
	"borrow_book/internal/repository"
	"context"
)

type ExtraService interface {
	RunRawQuery(ctx context.Context, q string) ([]map[string]interface{}, error)
}

type extraService struct {
	repo repository.ExtraRepository
}

func NewExtraService(repo repository.ExtraRepository) ExtraService {
	return &extraService{repo: repo}
}

func (s *extraService) RunRawQuery(ctx context.Context, q string) ([]map[string]interface{}, error) {
	// Potentially add security checks or whitelists here.
	return s.repo.ExecuteRawQuery(ctx, q)
}
