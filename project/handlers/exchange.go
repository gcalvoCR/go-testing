package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gcalvocr/go-testing/logger"
)

type ExchangeResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		logger.Warn("Missing exchange rate parameters", map[string]interface{}{
			"from": from,
			"to":   to,
		})
		http.Error(w, "Missing from or to parameter", http.StatusBadRequest)
		return
	}

	logger.Info("Fetching exchange rate", map[string]interface{}{
		"from": from,
		"to":   to,
	})

	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", from)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Failed to fetch exchange rate from external API", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var exchangeResp ExchangeResponse
	err = json.NewDecoder(resp.Body).Decode(&exchangeResp)
	if err != nil {
		logger.Error("Failed to decode exchange rate response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rate, exists := exchangeResp.Rates[to]
	if !exists {
		logger.Warn("Currency not found in exchange rate response", map[string]interface{}{
			"from": from,
			"to":   to,
		})
		http.Error(w, "Currency not found", http.StatusNotFound)
		return
	}

	logger.Info("Exchange rate retrieved successfully", map[string]interface{}{
		"from": from,
		"to":   to,
		"rate": rate,
	})

	result := map[string]float64{"rate": rate}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
