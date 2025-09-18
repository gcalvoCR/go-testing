package main

import (
	"os"

	"github.com/gcalvocr/go-testing/handlers"
	"github.com/gcalvocr/go-testing/logger"
	"github.com/gcalvocr/go-testing/server"
)

func main() {
	logger.Info("Starting Bank API application", nil)

	// Create and configure server
	srv := server.NewServer()

	// Initialize database
	err := srv.InitializeDatabase()
	if err != nil {
		logger.Error("Failed to initialize database", err)
		os.Exit(1)
	}
	defer func() {
		// Close database connection when server shuts down
		// This would be handled by a graceful shutdown in production
	}()

	// Get repository factory and set repositories in handlers
	repoFactory := srv.GetRepositoryFactory()
	if repoFactory != nil {
		handlers.SetAccountRepository(repoFactory.AccountRepo)
		handlers.SetTransactionRepository(repoFactory.TransactionRepo)
		logger.Info("Repositories initialized successfully", nil)
	} else {
		logger.Error("Repository factory is nil", nil)
		os.Exit(1)
	}

	// Setup routes
	srv.SetupRoutes()

	// Start server
	logger.Info("Starting server", map[string]interface{}{
		"port": srv.GetPort(),
	})

	err = srv.Start()
	if err != nil {
		logger.Error("Server failed to start", err)
		os.Exit(1)
	}
}
