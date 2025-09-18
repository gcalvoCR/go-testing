# Bank API - Go Testing Project

A simple REST API for bank operations built with Go, designed for teaching testing concepts.

## Features

- **Account Management**: Create and retrieve bank accounts
- **Transaction Processing**: Handle deposits and withdrawals with balance validation
- **Exchange Rates**: Fetch real-time currency exchange rates from external API
- **Multi-Database Support**: PostgreSQL and MongoDB with repository pattern
- **DTOs**: Clean data transfer objects for API communication
- **Repository Pattern**: Clean abstraction layer for database operations
- **Structured Logging**: Comprehensive logging with Logrus for observability
- **Docker Support**: Complete containerization with Docker Compose
- **Testing Suite**: Unit tests, integration tests, and Postman collections

## Quick Start

### With PostgreSQL (Default)
1. **Start the application:**
   ```bash
   make build
   make up
   ```

2. **Access the API:**
   - **📖 API Documentation:** `http://localhost:8080/` (HTML page with all endpoints)
   - **Health check:** `GET http://localhost:8080/health`
   - View logs: `make logs`
   - Run tests: `make test`

### With MongoDB
1. **Start with MongoDB:**
   ```bash
   make build
   make up-mongo
   ```

2. **Or run locally with MongoDB:**
   ```bash
   make run-local-mongo
   ```

### Switching Between Databases
To switch between PostgreSQL and MongoDB, simply change the `DB_TYPE` environment variable:

```bash
# For PostgreSQL (default)
DB_TYPE=postgres make up

# For MongoDB
DB_TYPE=mongodb make up-mongo
```

The application uses the Repository pattern to abstract database operations, making it easy to switch between different database implementations without changing the business logic.

### Testing with Postman
- Import `postman/Bank_API_Collection.postman_collection.json`
- Import `postman/Bank_API_Environment.postman_environment.json`
- Start testing all endpoints

## API Endpoints

### Documentation
- `GET /` - **📖 Interactive API Documentation** (HTML page with examples)

### Health Check
- `GET /health` - API health status

### Accounts
- `GET /accounts` - List all accounts
- `POST /accounts` - Create new account
- `GET /accounts/{id}` - Get account by ID

### Transactions
- `GET /accounts/{account_id}/transactions` - Get account transactions
- `POST /transactions` - Create transaction (deposit/withdrawal)

### Exchange Rates
- `GET /exchange?from=USD&to=EUR` - Get exchange rate

## Project Structure

```
├── main.go                 # Application orchestration (entry point)
├── server/                 # Server setup and configuration
│   └── server.go
├── handlers/               # HTTP request handlers
│   ├── account.go
│   ├── transaction.go
│   └── exchange.go
├── models/                 # Legacy data models
│   ├── account.go
│   └── transaction.go
├── dto/                    # Data Transfer Objects
│   ├── account.go
│   └── transaction.go
├── repository/             # Repository pattern implementation
│   ├── interface.go        # Repository interfaces
│   ├── postgres.go         # PostgreSQL implementation
│   └── mongodb.go          # MongoDB implementation
├── db/                     # Legacy database connection
│   └── db.go
├── logger/                 # Structured logging
│   └── logger.go
├── middleware/             # HTTP middleware
│   └── logging.go
├── tests/                  # Integration tests
│   └── integration_test.go
├── postman/                # Postman collections for testing
│   ├── Bank_API_Collection.postman_collection.json
│   ├── Bank_API_Environment.postman_environment.json
│   └── README.md
├── Dockerfile              # Docker image definition
├── docker-compose.yml      # Multi-container setup
├── Makefile               # Build and development commands
└── .env                   # Environment configuration
```

## Development Commands

```bash
# Build and run
make build          # Build Docker image
make up            # Start with PostgreSQL (default)
make up-mongo      # Start with MongoDB
make down          # Stop services
make logs          # View application logs

# Local development
make run-local     # Run locally with PostgreSQL
make run-local-mongo  # Run locally with MongoDB

# Testing
make test          # Run all tests
make test-unit     # Run unit tests
make test-integration  # Run integration tests

# Database
make db-logs       # View PostgreSQL logs
make db-shell      # Access PostgreSQL shell
make mongo-logs    # View MongoDB logs
make mongo-shell   # Access MongoDB shell

# Cleanup
make clean         # Remove containers and volumes
```

## Configuration

Environment variables in `.env`:

### Database Configuration
- `DB_TYPE` - Database type: `postgres` or `mongodb` (default: postgres)
- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - PostgreSQL user (default: postgres)
- `DB_PASSWORD` - PostgreSQL password (default: password)
- `DB_NAME` - PostgreSQL database name (default: bankdb)
- `MONGODB_URI` - MongoDB connection URI (default: mongodb://localhost:27017)

### Application Configuration
- `LOG_LEVEL` - Logging level: debug, info, warn, error (default: info)
- `PORT` - Server port (default: 8080)

## Testing with Postman

1. Import the collection: `postman/Bank_API_Collection.postman_collection.json`
2. Import the environment: `postman/Bank_API_Environment.postman_environment.json`
3. Set the environment to "Bank API Environment"
4. Start testing!

The collection includes:
- All API endpoints with proper headers and bodies
- Sample requests with realistic data
- Expected response examples
- Error handling examples

## Sample API Usage

### Create Account
```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "balance": 1000.00, "currency": "USD"}'
```

### Create Transaction
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "amount": 500.00, "type": "deposit"}'
```

### Get Exchange Rate
```bash
curl "http://localhost:8080/exchange?from=USD&to=EUR"
```

## Architecture

This project follows clean architecture principles with modern Go patterns:

### Core Components
- **`main.go`**: Pure orchestration - only starts the application
- **`server/`**: Server setup and configuration (routes, middleware)
- **`handlers/`**: HTTP request handling logic
- **`dto/`**: Data Transfer Objects for clean API communication
- **`repository/`**: Repository pattern for database abstraction
- **`models/`**: Legacy data structures (being phased out)
- **`logger/`**: Centralized logging configuration
- **`middleware/`**: HTTP middleware (logging, etc.)

### Design Patterns
- **Repository Pattern**: Clean abstraction over database operations
- **DTO Pattern**: Separation of internal models from API contracts
- **Dependency Injection**: Handlers receive repository instances
- **Clean Architecture**: Clear separation of concerns

### Database Abstraction
The system supports multiple databases through the repository pattern:

```go
// Repository interfaces define contracts
type AccountRepository interface {
    Create(ctx context.Context, account *dto.AccountDTO) error
    GetByID(ctx context.Context, id string) (*dto.AccountDTO, error)
    // ... other methods
}

// Multiple implementations
type PostgreSQLAccountRepository struct { /* PostgreSQL implementation */ }
type MongoDBAccountRepository struct { /* MongoDB implementation */ }
```

### DTOs (Data Transfer Objects)
Clean separation between internal data models and API contracts:

```go
// Request DTOs
type CreateAccountRequest struct {
    Name     string  `json:"name" validate:"required"`
    Balance  float64 `json:"balance" validate:"min=0"`
    Currency string  `json:"currency" validate:"required,len=3"`
}

// Response DTOs
type AccountResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Balance   float64   `json:"balance"`
    Currency  string    `json:"currency"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Learning Objectives

This project demonstrates:
- Clean architecture and separation of concerns
- REST API development with Go
- Database integration with PostgreSQL
- Structured logging and observability
- Unit and integration testing (tests don't depend on main.go)
- Docker containerization
- API testing with Postman
- Error handling and validation
- Transaction management
- Middleware implementation

## Technologies Used

- **Go** - Programming language
- **Gorilla Mux** - HTTP router
- **PostgreSQL** - Primary database
- **MongoDB** - Alternative database
- **Logrus** - Structured logging
- **Docker & Docker Compose** - Containerization
- **Postman** - API testing
- **Testify** - Testing framework

## License

This project is for educational purposes.