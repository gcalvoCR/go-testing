# Testcontainers MongoDB Demo

This project demonstrates how to use testcontainers in Go to test MongoDB operations.

## Structure

- `models/user.go`: User model
- `repository/interface.go`: Repository interface
- `repository/mongodb.go`: MongoDB implementation
- `main.go`: Main application
- `main_test.go`: Tests using testcontainers

## Setup

1. Start local MongoDB:
   ```bash
   docker-compose up -d
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

3. Run tests:
   ```bash
   go test
   ```

## Environment Variables

- `MONGODB_URI`: MongoDB connection string (default: `mongodb://admin:password@localhost:27017`)

The docker-compose.yml sets up MongoDB with authentication:
- Username: `admin`
- Password: `password`
- Database: `admin`

For local with auth: `mongodb://admin:password@localhost:27017`