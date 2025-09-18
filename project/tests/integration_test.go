package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gcalvocr/go-testing/server"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	// Create server instance
	srv := server.NewServer()
	srv.SetupRoutes()

	// Create test request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Get router and serve request
	router := srv.GetRouter()
	router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestGetExchangeRateIntegration(t *testing.T) {
	// Skip if running in CI or if external API might not be available
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping integration test in CI environment")
	}

	// Create server instance
	srv := server.NewServer()
	srv.SetupRoutes()

	// Create test request
	req, err := http.NewRequest("GET", "/exchange?from=USD&to=EUR", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Get router and serve request
	router := srv.GetRouter()
	router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)
	// The response should contain a rate field
	assert.Contains(t, rr.Body.String(), "rate")
}

func TestGetAccountsEndpoint(t *testing.T) {
	// Create server instance
	srv := server.NewServer()
	srv.SetupRoutes()

	// Create test request
	req, err := http.NewRequest("GET", "/accounts", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Get router and serve request
	router := srv.GetRouter()
	router.ServeHTTP(rr, req)

	// Since DB is not initialized in tests, we expect an internal server error
	// But the endpoint should be reachable and the middleware should work
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestServerCreation(t *testing.T) {
	// Test that we can create a server without main.go
	srv := server.NewServer()
	assert.NotNil(t, srv)
	assert.Equal(t, "8080", srv.GetPort())

	// Test that we can setup routes
	srv.SetupRoutes()
	router := srv.GetRouter()
	assert.NotNil(t, router)
}
