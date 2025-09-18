package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gcalvocr/go-testing/handlers"
	"github.com/gcalvocr/go-testing/logger"
	"github.com/gcalvocr/go-testing/middleware"
	"github.com/gcalvocr/go-testing/repository"
	"github.com/gorilla/mux"
)

// Server holds the server configuration
type Server struct {
	router      *mux.Router
	port        string
	repoFactory *repository.RepositoryFactory
}

// NewServer creates a new server instance
func NewServer() *Server {
	return &Server{
		router: mux.NewRouter(),
		port:   getEnv("PORT", "8080"),
	}
}

// SetupRoutes configures all the API routes
func (s *Server) SetupRoutes() {
	// Add logging middleware
	s.router.Use(middleware.LoggingMiddleware)

	// Root route - API documentation
	s.router.HandleFunc("/", handlers.IndexHandler).Methods("GET")

	// Health check endpoint
	s.router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// Account routes
	s.router.HandleFunc("/accounts", handlers.GetAccounts).Methods("GET")
	s.router.HandleFunc("/accounts", handlers.CreateAccount).Methods("POST")
	s.router.HandleFunc("/accounts/{id}", handlers.GetAccountByID).Methods("GET")

	// Transaction routes
	s.router.HandleFunc("/accounts/{account_id}/transactions", handlers.GetTransactionsByAccountID).Methods("GET")
	s.router.HandleFunc("/transactions", handlers.CreateTransaction).Methods("POST")

	// Exchange rate route
	s.router.HandleFunc("/exchange", handlers.GetExchangeRate).Methods("GET")
}

// InitializeDatabase sets up the database connection and repositories
func (s *Server) InitializeDatabase() error {
	dbType := getEnv("DB_TYPE", "postgres")

	var connectionString string
	var err error

	switch repository.DatabaseType(dbType) {
	case repository.PostgreSQL:
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "password")
		dbname := getEnv("DB_NAME", "bankdb")

		connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		logger.Info("Connecting to PostgreSQL database", map[string]interface{}{
			"host":     host,
			"port":     port,
			"database": dbname,
		})

	case repository.MongoDB:
		connectionString = getEnv("MONGODB_URI", "mongodb://localhost:27017")

		logger.Info("Connecting to MongoDB database", map[string]interface{}{
			"uri": connectionString,
		})

	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	s.repoFactory, err = repository.NewRepositoryFactory(repository.DatabaseType(dbType), connectionString)
	if err != nil {
		logger.Error("Failed to initialize repository factory", err)
		return err
	}

	logger.Info("Database initialized successfully", map[string]interface{}{
		"db_type": dbType,
	})
	return nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	logger.Info("Server starting", map[string]interface{}{
		"port": s.port,
	})

	return http.ListenAndServe(":"+s.port, s.router)
}

// GetRouter returns the router (useful for testing)
func (s *Server) GetRouter() *mux.Router {
	return s.router
}

// GetPort returns the server port
func (s *Server) GetPort() string {
	return s.port
}

// GetRepositoryFactory returns the repository factory
func (s *Server) GetRepositoryFactory() *repository.RepositoryFactory {
	return s.repoFactory
}

// getEnv gets environment variable with fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
