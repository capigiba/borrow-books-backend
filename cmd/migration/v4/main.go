package main

import (
	"borrow_book/internal/config"
	"borrow_book/internal/initialize"
	"borrow_book/pkg/logger"
	"os"
)

var log logger.Logger

func main() {
	log = logger.NewLogger("migration v4")
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

	// Add UserName column to borrows table
	addColumnQuery := `
		ALTER TABLE borrows
		ADD COLUMN IF NOT EXISTS user_name VARCHAR(255);
	`
	_, err = tx.Exec(addColumnQuery)
	if err != nil {
		log.Fatalf("Error adding new column 'user_name': %v", err)
	}
	log.Infof("Added new column 'user_name' successfully.")
}
