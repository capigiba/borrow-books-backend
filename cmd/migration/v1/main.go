package main

import (
	"borrow_book/internal/config"
	"borrow_book/internal/initialize"
	"borrow_book/pkg/logger"
	"os"
)

var log logger.Logger

func main() {
	log = logger.NewLogger("migration v1")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("config loading error: %w", err)
		os.Exit(1)
	}

	// Initialize databases
	db, err := initialize.InitDatabases(cfg, &log)
	if err != nil {
		log.Errorf("database initialization error: %w", err)
		os.Exit(1)
	}

	// Begin the migration within a transaction
	tx, err := db.Begin()
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
			log.Infof("Migration completed successfully.")
		}
	}()

	// Step 1: Add a new column 'published_at_new' of type BIGINT
	addColumnQuery := `
		ALTER TABLE books
		ADD COLUMN IF NOT EXISTS published_at_new BIGINT;
	`
	_, err = tx.Exec(addColumnQuery)
	if err != nil {
		log.Fatalf("Error adding new column 'published_at_new': %v", err)
	}
	log.Infof("Added new column 'published_at_new'.")

	// Step 2: Populate 'published_at_new' with UNIX timestamps from 'published_at'
	updateColumnQuery := `
		UPDATE books
		SET published_at_new = EXTRACT(EPOCH FROM published_at)::BIGINT;
	`
	_, err = tx.Exec(updateColumnQuery)
	if err != nil {
		log.Fatalf("Error updating 'published_at_new' with UNIX timestamps: %v", err)
	}
	log.Infof("Populated 'published_at_new' with UNIX timestamps.")

	// Step 3: Drop the old 'published_at' column
	dropOldColumnQuery := `
		ALTER TABLE books
		DROP COLUMN IF EXISTS published_at;
	`
	_, err = tx.Exec(dropOldColumnQuery)
	if err != nil {
		log.Fatalf("Error dropping old column 'published_at': %v", err)
	}
	log.Infof("Dropped old column 'published_at'.")

	// Step 4: Rename 'published_at_new' to 'published_at'
	renameColumnQuery := `
		ALTER TABLE books
		RENAME COLUMN published_at_new TO published_at;
	`
	_, err = tx.Exec(renameColumnQuery)
	if err != nil {
		log.Fatalf("Error renaming 'published_at_new' to 'published_at': %v", err)
	}
	log.Infof("Renamed 'published_at_new' to 'published_at'.")
}
