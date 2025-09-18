package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// APISuite is a test suite for API endpoints
type APISuite struct {
	suite.Suite
	server *httptest.Server
	client *http.Client
}

// SetupSuite runs once before all tests in the suite
func (suite *APISuite) SetupSuite() {
	// Create a test server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", suite.handleHealth)
	mux.HandleFunc("/users", suite.handleUsers)
	mux.HandleFunc("/users/", suite.handleUserByID)

	suite.server = httptest.NewServer(mux)
	suite.client = suite.server.Client()
}

// TearDownSuite runs once after all tests in the suite
func (suite *APISuite) TearDownSuite() {
	suite.server.Close()
}

// SetupTest runs before each test
func (suite *APISuite) SetupTest() {
	// Reset any test state if needed
}

// TearDownTest runs after each test
func (suite *APISuite) TearDownTest() {
	// Clean up after each test
}

// TestHealthEndpoint tests the health check endpoint
func (suite *APISuite) TestHealthEndpoint() {
	resp, err := suite.client.Get(suite.server.URL + "/health")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
}

// TestGetUsers tests the users endpoint
func (suite *APISuite) TestGetUsers() {
	resp, err := suite.client.Get(suite.server.URL + "/users")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
}

// TestCreateUser tests user creation
func (suite *APISuite) TestCreateUser() {
	// This would test POST /users
	suite.T().Skip("POST endpoint not implemented in this example")
}

// TestGetUserByID tests getting a user by ID
func (suite *APISuite) TestGetUserByID() {
	resp, err := suite.client.Get(suite.server.URL + "/users/123")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
}

// TestInvalidUserID tests invalid user ID
func (suite *APISuite) TestInvalidUserID() {
	resp, err := suite.client.Get(suite.server.URL + "/users/invalid")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	defer resp.Body.Close()
}

// TestNonExistentEndpoint tests 404 handling
func (suite *APISuite) TestNonExistentEndpoint() {
	resp, err := suite.client.Get(suite.server.URL + "/nonexistent")

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

	defer resp.Body.Close()
}

// HTTP handler methods for the test server
func (suite *APISuite) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}

func (suite *APISuite) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[{"id": 1, "name": "John"}, {"id": 2, "name": "Jane"}]`))
}

func (suite *APISuite) handleUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple path parsing (in real code, use a proper router)
	path := r.URL.Path
	if len(path) < 8 { // "/users/" is 7 chars
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userID := path[7:] // Extract ID after "/users/"
	if userID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// In a real implementation, you'd look up the user
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id": ` + userID + `, "name": "User ` + userID + `"}`))
}

// TestAPISuite runs the API test suite
func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}
