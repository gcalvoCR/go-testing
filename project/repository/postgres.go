package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gcalvocr/go-testing/dto"
	"github.com/gcalvocr/go-testing/logger"
	_ "github.com/lib/pq"
)

// PostgreSQLAccountRepository implements AccountRepository for PostgreSQL
type PostgreSQLAccountRepository struct {
	db *sql.DB
}

// PostgreSQLTransactionRepository implements TransactionRepository for PostgreSQL
type PostgreSQLTransactionRepository struct {
	db *sql.DB
}

// newPostgreSQLFactory creates PostgreSQL repository instances
func newPostgreSQLFactory(connectionString string) (*RepositoryFactory, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	logger.Info("Connected to PostgreSQL database", nil)

	// Create tables if they don't exist
	if err := createPostgreSQLTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &RepositoryFactory{
		AccountRepo:     &PostgreSQLAccountRepository{db: db},
		TransactionRepo: &PostgreSQLTransactionRepository{db: db},
	}, nil
}

// createPostgreSQLTables creates the necessary tables
func createPostgreSQLTables(db *sql.DB) error {
	accountTable := `
	CREATE TABLE IF NOT EXISTS accounts (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		balance DECIMAL(15,2) NOT NULL DEFAULT 0,
		currency VARCHAR(3) NOT NULL DEFAULT 'USD',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	transactionTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id VARCHAR(36) PRIMARY KEY,
		account_id VARCHAR(36) REFERENCES accounts(id),
		amount DECIMAL(15,2) NOT NULL,
		type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(accountTable); err != nil {
		return err
	}

	if _, err := db.Exec(transactionTable); err != nil {
		return err
	}

	logger.Info("PostgreSQL tables created successfully", nil)
	return nil
}

// Account repository methods
func (r *PostgreSQLAccountRepository) Create(ctx context.Context, account *dto.AccountDTO) error {
	query := `
		INSERT INTO accounts (id, name, balance, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		account.ID, account.Name, account.Balance, account.Currency,
		account.CreatedAt, account.UpdatedAt)

	if err != nil {
		logger.Error("Failed to create account in PostgreSQL", err)
		return err
	}

	logger.Info("Account created in PostgreSQL", map[string]interface{}{
		"account_id": account.ID,
		"name":       account.Name,
	})
	return nil
}

func (r *PostgreSQLAccountRepository) GetByID(ctx context.Context, id string) (*dto.AccountDTO, error) {
	query := `
		SELECT id, name, balance, currency, created_at, updated_at
		FROM accounts WHERE id = $1`

	var account dto.AccountDTO
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID, &account.Name, &account.Balance,
		&account.Currency, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Account not found
		}
		logger.Error("Failed to get account from PostgreSQL", err)
		return nil, err
	}

	return &account, nil
}

func (r *PostgreSQLAccountRepository) GetAll(ctx context.Context) ([]*dto.AccountDTO, error) {
	query := `SELECT id, name, balance, currency, created_at, updated_at FROM accounts`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error("Failed to query accounts from PostgreSQL", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []*dto.AccountDTO
	for rows.Next() {
		var account dto.AccountDTO
		err := rows.Scan(
			&account.ID, &account.Name, &account.Balance,
			&account.Currency, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			logger.Error("Failed to scan account from PostgreSQL", err)
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (r *PostgreSQLAccountRepository) Update(ctx context.Context, id string, update *dto.UpdateAccountRequest) error {
	query := `UPDATE accounts SET `
	args := []interface{}{}
	argCount := 0

	if update.Name != nil {
		argCount++
		query += fmt.Sprintf("name = $%d, ", argCount)
		args = append(args, *update.Name)
	}

	if update.Balance != nil {
		argCount++
		query += fmt.Sprintf("balance = $%d, ", argCount)
		args = append(args, *update.Balance)
	}

	if update.Currency != nil {
		argCount++
		query += fmt.Sprintf("currency = $%d, ", argCount)
		args = append(args, *update.Currency)
	}

	if argCount == 0 {
		return nil // Nothing to update
	}

	argCount++
	query += fmt.Sprintf("updated_at = $%d WHERE id = $%d", argCount, argCount+1)
	args = append(args, time.Now(), id)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error("Failed to update account in PostgreSQL", err)
		return err
	}

	logger.Info("Account updated in PostgreSQL", map[string]interface{}{
		"account_id": id,
	})
	return nil
}

func (r *PostgreSQLAccountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("Failed to delete account from PostgreSQL", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	logger.Info("Account deleted from PostgreSQL", map[string]interface{}{
		"account_id":    id,
		"rows_affected": rowsAffected,
	})
	return nil
}

func (r *PostgreSQLAccountRepository) GetByName(ctx context.Context, name string) (*dto.AccountDTO, error) {
	query := `
		SELECT id, name, balance, currency, created_at, updated_at
		FROM accounts WHERE name = $1`

	var account dto.AccountDTO
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&account.ID, &account.Name, &account.Balance,
		&account.Currency, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Account not found
		}
		logger.Error("Failed to get account by name from PostgreSQL", err)
		return nil, err
	}

	return &account, nil
}

func (r *PostgreSQLAccountRepository) UpdateBalance(ctx context.Context, id string, newBalance float64) error {
	query := `UPDATE accounts SET balance = $1, updated_at = $2 WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, newBalance, time.Now(), id)
	if err != nil {
		logger.Error("Failed to update account balance in PostgreSQL", err)
		return err
	}

	logger.Info("Account balance updated in PostgreSQL", map[string]interface{}{
		"account_id":  id,
		"new_balance": newBalance,
	})
	return nil
}

// Transaction repository methods
func (r *PostgreSQLTransactionRepository) Create(ctx context.Context, transaction *dto.TransactionDTO) error {
	query := `
		INSERT INTO transactions (id, account_id, amount, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		transaction.ID, transaction.AccountID, transaction.Amount,
		transaction.Type, transaction.CreatedAt, transaction.UpdatedAt)

	if err != nil {
		logger.Error("Failed to create transaction in PostgreSQL", err)
		return err
	}

	logger.Info("Transaction created in PostgreSQL", map[string]interface{}{
		"transaction_id": transaction.ID,
		"account_id":     transaction.AccountID,
		"amount":         transaction.Amount,
		"type":           transaction.Type,
	})
	return nil
}

func (r *PostgreSQLTransactionRepository) GetByID(ctx context.Context, id string) (*dto.TransactionDTO, error) {
	query := `
		SELECT id, account_id, amount, type, created_at, updated_at
		FROM transactions WHERE id = $1`

	var transaction dto.TransactionDTO
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID, &transaction.AccountID, &transaction.Amount,
		&transaction.Type, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Transaction not found
		}
		logger.Error("Failed to get transaction from PostgreSQL", err)
		return nil, err
	}

	return &transaction, nil
}

func (r *PostgreSQLTransactionRepository) GetByAccountID(ctx context.Context, accountID string) ([]*dto.TransactionDTO, error) {
	query := `
		SELECT id, account_id, amount, type, created_at, updated_at
		FROM transactions WHERE account_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		logger.Error("Failed to query transactions by account ID from PostgreSQL", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []*dto.TransactionDTO
	for rows.Next() {
		var transaction dto.TransactionDTO
		err := rows.Scan(
			&transaction.ID, &transaction.AccountID, &transaction.Amount,
			&transaction.Type, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			logger.Error("Failed to scan transaction from PostgreSQL", err)
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (r *PostgreSQLTransactionRepository) GetAll(ctx context.Context) ([]*dto.TransactionDTO, error) {
	query := `SELECT id, account_id, amount, type, created_at, updated_at FROM transactions`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error("Failed to query all transactions from PostgreSQL", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []*dto.TransactionDTO
	for rows.Next() {
		var transaction dto.TransactionDTO
		err := rows.Scan(
			&transaction.ID, &transaction.AccountID, &transaction.Amount,
			&transaction.Type, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			logger.Error("Failed to scan transaction from PostgreSQL", err)
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (r *PostgreSQLTransactionRepository) Update(ctx context.Context, id string, transaction *dto.TransactionDTO) error {
	query := `
		UPDATE transactions SET account_id = $1, amount = $2, type = $3, updated_at = $4
		WHERE id = $5`

	transaction.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		transaction.AccountID, transaction.Amount, transaction.Type,
		transaction.UpdatedAt, id)

	if err != nil {
		logger.Error("Failed to update transaction in PostgreSQL", err)
		return err
	}

	logger.Info("Transaction updated in PostgreSQL", map[string]interface{}{
		"transaction_id": id,
	})
	return nil
}

func (r *PostgreSQLTransactionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM transactions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("Failed to delete transaction from PostgreSQL", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	logger.Info("Transaction deleted from PostgreSQL", map[string]interface{}{
		"transaction_id": id,
		"rows_affected":  rowsAffected,
	})
	return nil
}

func (r *PostgreSQLTransactionRepository) GetTransactionSummary(ctx context.Context, accountID string) (*dto.TransactionSummary, error) {
	query := `
		SELECT
			COUNT(*) as total_transactions,
			COALESCE(SUM(CASE WHEN type = 'deposit' THEN amount ELSE 0 END), 0) as total_deposits,
			COALESCE(SUM(CASE WHEN type = 'withdrawal' THEN amount ELSE 0 END), 0) as total_withdrawals,
			MAX(created_at) as last_transaction_at
		FROM transactions
		WHERE account_id = $1`

	var summary dto.TransactionSummary
	summary.AccountID = accountID

	err := r.db.QueryRowContext(ctx, query, accountID).Scan(
		&summary.TotalTransactions,
		&summary.TotalDeposits,
		&summary.TotalWithdrawals,
		&summary.LastTransactionAt)

	if err != nil {
		logger.Error("Failed to get transaction summary from PostgreSQL", err)
		return nil, err
	}

	// Get current balance from accounts table
	balanceQuery := `SELECT balance FROM accounts WHERE id = $1`
	err = r.db.QueryRowContext(ctx, balanceQuery, accountID).Scan(&summary.CurrentBalance)
	if err != nil {
		logger.Error("Failed to get current balance from PostgreSQL", err)
		return nil, err
	}

	return &summary, nil
}
