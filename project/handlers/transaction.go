package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gcalvocr/go-testing/dto"
	"github.com/gcalvocr/go-testing/logger"
	"github.com/gcalvocr/go-testing/repository"
	"github.com/gorilla/mux"
)

var transactionRepo repository.TransactionRepository

// SetTransactionRepository sets the transaction repository (called from main)
func SetTransactionRepository(repo repository.TransactionRepository) {
	transactionRepo = repo
}

func GetTransactionsByAccountID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["account_id"]

	logger.Info("Getting transactions for account", map[string]interface{}{
		"account_id": accountID,
	})

	if transactionRepo == nil {
		logger.Error("Transaction repository not initialized", nil)
		http.Error(w, "Database not available", http.StatusInternalServerError)
		return
	}

	transactions, err := transactionRepo.GetByAccountID(r.Context(), accountID)
	if err != nil {
		logger.Error("Failed to get transactions", err)
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		response[i] = dto.TransactionResponse{
			ID:        tx.ID,
			AccountID: tx.AccountID,
			Amount:    tx.Amount,
			Type:      tx.Type,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
		}
	}

	logger.Info("Retrieved transactions successfully", map[string]interface{}{
		"account_id": accountID,
		"count":      len(transactions),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("Failed to decode transaction JSON", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	logger.Info("Creating transaction", map[string]interface{}{
		"account_id": req.AccountID,
		"amount":     req.Amount,
		"type":       req.Type,
	})

	if accountRepo == nil || transactionRepo == nil {
		logger.Error("Repositories not initialized", nil)
		http.Error(w, "Database not available", http.StatusInternalServerError)
		return
	}

	// Get current account to check balance
	account, err := accountRepo.GetByID(r.Context(), req.AccountID)
	if err != nil {
		logger.Error("Failed to get account for transaction", err)
		http.Error(w, "Failed to retrieve account", http.StatusInternalServerError)
		return
	}

	if account == nil {
		logger.Warn("Account not found for transaction", map[string]interface{}{
			"account_id": req.AccountID,
		})
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Calculate new balance
	var newBalance float64
	if req.Type == "deposit" {
		newBalance = account.Balance + req.Amount
		logger.Info("Processing deposit", map[string]interface{}{
			"account_id":      req.AccountID,
			"current_balance": account.Balance,
			"deposit_amount":  req.Amount,
			"new_balance":     newBalance,
		})
	} else if req.Type == "withdrawal" {
		newBalance = account.Balance - req.Amount
		if newBalance < 0 {
			logger.Warn("Insufficient funds for withdrawal", map[string]interface{}{
				"account_id":        req.AccountID,
				"current_balance":   account.Balance,
				"withdrawal_amount": req.Amount,
			})
			http.Error(w, "Insufficient funds", http.StatusBadRequest)
			return
		}
		logger.Info("Processing withdrawal", map[string]interface{}{
			"account_id":        req.AccountID,
			"current_balance":   account.Balance,
			"withdrawal_amount": req.Amount,
			"new_balance":       newBalance,
		})
	} else {
		logger.Warn("Invalid transaction type", map[string]interface{}{
			"account_id": req.AccountID,
			"type":       req.Type,
		})
		http.Error(w, "Invalid transaction type", http.StatusBadRequest)
		return
	}

	// Update account balance
	err = accountRepo.UpdateBalance(r.Context(), req.AccountID, newBalance)
	if err != nil {
		logger.Error("Failed to update account balance", err)
		http.Error(w, "Failed to update account balance", http.StatusInternalServerError)
		return
	}

	// Create transaction
	transaction := &dto.TransactionDTO{
		AccountID: req.AccountID,
		Amount:    req.Amount,
		Type:      req.Type,
	}

	err = transactionRepo.Create(r.Context(), transaction)
	if err != nil {
		logger.Error("Failed to create transaction", err)
		// Rollback balance update by reverting it
		_ = accountRepo.UpdateBalance(r.Context(), req.AccountID, account.Balance)
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	response := dto.TransactionResponse{
		ID:        transaction.ID,
		AccountID: transaction.AccountID,
		Amount:    transaction.Amount,
		Type:      transaction.Type,
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}

	logger.Info("Transaction completed successfully", map[string]interface{}{
		"transaction_id": transaction.ID,
		"account_id":     transaction.AccountID,
		"amount":         transaction.Amount,
		"type":           transaction.Type,
		"new_balance":    newBalance,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
