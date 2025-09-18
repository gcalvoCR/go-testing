package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gcalvocr/go-testing/dto"
	"github.com/gcalvocr/go-testing/logger"
	"github.com/gcalvocr/go-testing/repository"
	"github.com/gorilla/mux"
)

var accountRepo repository.AccountRepository

// SetAccountRepository sets the account repository (called from main)
func SetAccountRepository(repo repository.AccountRepository) {
	accountRepo = repo
}

func GetAccounts(w http.ResponseWriter, r *http.Request) {
	logger.Info("Getting all accounts", nil)

	if accountRepo == nil {
		logger.Error("Account repository not initialized", nil)
		http.Error(w, "Database not available", http.StatusInternalServerError)
		return
	}

	accounts, err := accountRepo.GetAll(r.Context())
	if err != nil {
		logger.Error("Failed to get accounts", err)
		http.Error(w, "Failed to retrieve accounts", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	response := make([]dto.AccountResponse, len(accounts))
	for i, acc := range accounts {
		response[i] = dto.AccountResponse{
			ID:        acc.ID,
			Name:      acc.Name,
			Balance:   acc.Balance,
			Currency:  acc.Currency,
			CreatedAt: acc.CreatedAt,
			UpdatedAt: acc.UpdatedAt,
		}
	}

	logger.Info("Retrieved accounts successfully", map[string]interface{}{
		"count": len(accounts),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAccountByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	logger.Info("Getting account by ID", map[string]interface{}{
		"account_id": id,
	})

	if accountRepo == nil {
		logger.Error("Account repository not initialized", nil)
		http.Error(w, "Database not available", http.StatusInternalServerError)
		return
	}

	account, err := accountRepo.GetByID(r.Context(), id)
	if err != nil {
		logger.Error("Failed to get account", err)
		http.Error(w, "Failed to retrieve account", http.StatusInternalServerError)
		return
	}

	if account == nil {
		logger.Warn("Account not found", map[string]interface{}{
			"account_id": id,
		})
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	response := dto.AccountResponse{
		ID:        account.ID,
		Name:      account.Name,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	logger.Info("Retrieved account successfully", map[string]interface{}{
		"account_id": id,
		"name":       account.Name,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Health check requested", nil)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("Failed to decode account JSON", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	logger.Info("Creating new account", map[string]interface{}{
		"name":     req.Name,
		"balance":  req.Balance,
		"currency": req.Currency,
	})

	if accountRepo == nil {
		logger.Error("Account repository not initialized", nil)
		http.Error(w, "Database not available", http.StatusInternalServerError)
		return
	}

	// Create account DTO
	account := &dto.AccountDTO{
		Name:     req.Name,
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	err = accountRepo.Create(r.Context(), account)
	if err != nil {
		logger.Error("Failed to create account", err)
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	response := dto.AccountResponse{
		ID:        account.ID,
		Name:      account.Name,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	logger.Info("Account created successfully", map[string]interface{}{
		"account_id": account.ID,
		"name":       account.Name,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
