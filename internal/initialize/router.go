package initialize

import (
	"borrow_book/internal/router"

	"github.com/jmoiron/sqlx"
)

// InitRouter sets up the application router using dependency injection.
func InitAppRouter(pgDB *sqlx.DB) (*router.AppRouter, error) {
	appRouter, err := InitializeApp(pgDB)
	if err != nil {
		return nil, err
	}
	return appRouter, nil
}
