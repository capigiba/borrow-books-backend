// wire.go
//go:build wireinject
// +build wireinject

package initialize

import (
	"borrow_book/internal/handler"
	"borrow_book/internal/repository"
	"borrow_book/internal/router"
	"borrow_book/internal/service"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

func InitializeApp(db *sqlx.DB) (*router.AppRouter, error) {
	wire.Build(
		handler.ProviderSetHandler,
		service.ProviderSetService,
		repository.ProviderSetRepository,
		router.ProviderSetRouter,
	)
	return &router.AppRouter{}, nil
}
