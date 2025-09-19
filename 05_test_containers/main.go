package main

import (
	"context"
	"errors"
	"log"
	"os"
	"testcontainers-mongo-demo/config"
	"testcontainers-mongo-demo/repository"

	"github.com/joho/godotenv"
)

func InitializeDataBase() (repository.UserRepository, error) {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		// defaults to mongo
		dbType = string(repository.MongoDB)
	}

	cfg := config.LoadDatabaseConfig()

	switch repository.DatabaseType(dbType) {
	case repository.MongoDB:
		repo, err := repository.NewMongoUserRepository(cfg)
		if err != nil {
			return nil, err
		}
		return repo, nil
	default:
		return nil, errors.New("DB type not implemented")
	}
}

func RunApp() error {
	repo, err := InitializeDataBase()
	if err != nil {
		return err
	}

	// Create server with connection initialization
	server, err := NewServer(repo)
	if err != nil {
		return err
	}

	defer server.repo.Close(context.Background())

	return server.RunDemo()
}

func main() {
	godotenv.Load()
	if err := RunApp(); err != nil {
		log.Fatal(err)
	}
}
