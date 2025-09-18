package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExchangeRate(t *testing.T) {
	req, err := http.NewRequest("GET", "/exchange?from=USD&to=EUR", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetExchangeRate)
	handler.ServeHTTP(rr, req)

	// Note: This test makes a real HTTP call to the external API
	// In a real scenario, you should mock the HTTP client for reliable tests
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]float64
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response, "rate")
	assert.Greater(t, response["rate"], 0.0)
}

func TestGetExchangeRateMissingParams(t *testing.T) {
	req, err := http.NewRequest("GET", "/exchange?from=USD", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetExchangeRate)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
