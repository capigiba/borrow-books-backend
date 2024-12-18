package main

import (
	"borrow_book/internal/config"
	"borrow_book/internal/initialize"
	"borrow_book/pkg/logger"
	"os"
)

var log logger.Logger

func main() {
	log = logger.NewLogger("migration v2")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("config loading error: %v", err)
		os.Exit(1)
	}

	// Initialize databases
	db, err := initialize.InitDatabases(cfg, &log)
	if err != nil {
		log.Errorf("database initialization error: %v", err)
		os.Exit(1)
	}
	defer db.Close()

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

	// Step 1: Add a new column 'borrowed_at_new' of type BIGINT
	addColumnQuery := `
		ALTER TABLE borrows
		ADD COLUMN IF NOT EXISTS borrowed_at_new BIGINT;
	`
	_, err = tx.Exec(addColumnQuery)
	if err != nil {
		log.Fatalf("Error adding new column 'borrowed_at_new': %v", err)
	}
	log.Infof("Added new column 'borrowed_at_new'.")

	// Step 2: Populate 'borrowed_at_new' with UNIX timestamps from 'borrowed_at'
	updateColumnQuery := `
		UPDATE borrows
		SET borrowed_at_new = EXTRACT(EPOCH FROM borrowed_at)::BIGINT;
	`
	_, err = tx.Exec(updateColumnQuery)
	if err != nil {
		log.Fatalf("Error updating 'borrowed_at_new' with UNIX timestamps: %v", err)
	}
	log.Infof("Populated 'borrowed_at_new' with UNIX timestamps.")

	// Step 3: Drop the old 'borrowed_at' column
	dropOldColumnQuery := `
		ALTER TABLE borrows
		DROP COLUMN IF EXISTS borrowed_at;
	`
	_, err = tx.Exec(dropOldColumnQuery)
	if err != nil {
		log.Fatalf("Error dropping old column 'borrowed_at': %v", err)
	}
	log.Infof("Dropped old column 'borrowed_at'.")

	// Step 4: Rename 'borrowed_at_new' to 'borrowed_at'
	renameColumnQuery := `
		ALTER TABLE borrows
		RENAME COLUMN borrowed_at_new TO borrowed_at;
	`
	_, err = tx.Exec(renameColumnQuery)
	if err != nil {
		log.Fatalf("Error renaming 'borrowed_at_new' to 'borrowed_at': %v", err)
	}
	log.Infof("Renamed 'borrowed_at_new' to 'borrowed_at'.")
}
