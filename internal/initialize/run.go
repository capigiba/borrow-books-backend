package initialize

import (
	"borrow_book/internal/config"
	"borrow_book/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	appLogger := logger.NewLogger("Initialize")

	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Errorf("config loading error: %w", err)
		os.Exit(1)
	}

	// Initialize databases
	dbConn, err := InitDatabases(cfg, &appLogger)
	if err != nil {
		appLogger.Errorf("database initialization error: %w", err)
		os.Exit(1)
	}

	// // Initialize localization
	// localizer, err := localization.NewLocalizer(cfg, &appLogger)
	// if err != nil {
	// 	appLogger.Errorf("localization initialization error: %w", err)
	// 	os.Exit(1)
	// }

	// Initialize Router
	appRouter, err := InitAppRouter(dbConn)
	if err != nil {
		appLogger.Errorf("AppRouter initialization failed: %v", err)
		os.Exit(1)
	}

	// Initialize Server
	server := InitServer(cfg, appRouter)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-quit
		appLogger.Infof("Received signal '%v'. Shutting down server...", sig)
		if err := server.Shutdown(); err != nil {
			appLogger.Warnf("Server forced to shutdown: %v", err)
		}
		appLogger.Info("Server exiting")
	}()

	// Start the server
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		appLogger.Errorf("Failed to run the server: %v", err)
		os.Exit(1)
	}
}
