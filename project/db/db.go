package db

// import (
// 	"database/sql"
// 	"fmt"

// 	"github.com/gcalvocr/go-testing/logger"
// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// func InitDB(host, port, user, password, dbname string) error {
// 	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)

// 	logger.Info("Initializing database connection", map[string]interface{}{
// 		"host":     host,
// 		"port":     port,
// 		"database": dbname,
// 	})

// 	var err error
// 	DB, err = sql.Open("postgres", dataSourceName)
// 	if err != nil {
// 		logger.Error("Failed to open database connection", err)
// 		return err
// 	}

// 	err = DB.Ping()
// 	if err != nil {
// 		logger.Error("Failed to ping database", err)
// 		return err
// 	}

// 	logger.Info("Database connection established successfully", nil)
// 	return nil
// }

// func CloseDB() {
// 	if DB != nil {
// 		logger.Info("Closing database connection", nil)
// 		DB.Close()
// 	}
// }

// func CreateTables() error {
// 	logger.Info("Creating database tables", nil)

// 	accountTable := `
// 	CREATE TABLE IF NOT EXISTS accounts (
// 		id SERIAL PRIMARY KEY,
// 		name VARCHAR(255) NOT NULL,
// 		balance DECIMAL(10,2) NOT NULL DEFAULT 0,
// 		currency VARCHAR(3) NOT NULL DEFAULT 'USD',
// 		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 	);`

// 	transactionTable := `
// 	CREATE TABLE IF NOT EXISTS transactions (
// 		id SERIAL PRIMARY KEY,
// 		account_id INT REFERENCES accounts(id),
// 		amount DECIMAL(10,2) NOT NULL,
// 		type VARCHAR(50) NOT NULL,
// 		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 	);`

// 	_, err := DB.Exec(accountTable)
// 	if err != nil {
// 		logger.Error("Failed to create accounts table", err)
// 		return err
// 	}

// 	_, err = DB.Exec(transactionTable)
// 	if err != nil {
// 		logger.Error("Failed to create transactions table", err)
// 		return err
// 	}

// 	logger.Info("Database tables created successfully", nil)
// 	return nil
// }
