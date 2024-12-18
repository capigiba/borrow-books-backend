package initialize

import (
	"borrow_book/internal/config"
	"borrow_book/internal/infra/database/postgres"
	"borrow_book/pkg/logger"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// InitDatabase establishes a database connection and runs migrations.
func InitDatabases(cfg *config.Config, log *logger.Logger) (*sqlx.DB, error) {
	pgDB, err := postgres.NewPostgresDB(cfg.Database.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	log.Info("Database connection established successfully!")
	return pgDB, nil
}
