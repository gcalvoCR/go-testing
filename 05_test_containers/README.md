# Testcontainers in Go: Complete Guide & Implementation

This project demonstrates comprehensive testcontainers usage in Go applications, focusing on MongoDB integration testing. Learn how to implement, leverage, and optimize test containers for robust testing strategies.

## 🎯 What are Test Containers?

Testcontainers is a library that provides lightweight, throwaway instances of databases, message brokers, web browsers, or anything that can run in a Docker container. These containers are perfect for integration testing where you need real dependencies without the complexity of managing external services.

## 📋 Table of Contents

- [Why Test Containers?](#-why-test-containers)
- [Advantages & Disadvantages](#-advantages--disadvantages)
- [When to Use Test Containers](#-when-to-use-test-containers)
- [Progressive Implementation Guide](#-progressive-implementation-guide)
- [Project Structure](#-project-structure)
- [Quick Start](#-quick-start)
- [Advanced Patterns](#-advanced-patterns)
- [Best Practices](#-best-practices)
- [Troubleshooting](#-troubleshooting)

## 🚀 Why Test Containers?

Test containers solve the fundamental problem of testing code that depends on external services. Instead of:

- ❌ Mocking complex database behaviors
- ❌ Maintaining shared test databases
- ❌ Dealing with flaky external dependencies
- ❌ Testing against stale or inconsistent data

You get:

- ✅ Real database behavior and constraints
- ✅ Isolated, clean test environments
- ✅ Consistent test execution across environments
- ✅ Ability to test complex interactions

## ⚖️ Advantages & Disadvantages

### Advantages

**🔧 Realism**
- Test against real database implementations
- Catch integration issues early
- Validate actual query performance
- Test database-specific features (indexes, transactions, etc.)

**🔄 Isolation**
- Each test gets a fresh database instance
- No test interference or data pollution
- Parallel test execution without conflicts
- Clean state for every test run

**🚀 Developer Experience**
- No external dependencies to manage
- Consistent environments across team
- Easy CI/CD integration
- Fast feedback loops

**📈 Scalability**
- Run tests in parallel without resource conflicts
- Scale test suites without infrastructure complexity
- Cost-effective for large test suites

### Disadvantages

**⏱️ Performance**
- Container startup time (typically 1-3 seconds)
- Slower than in-memory databases or mocks
- Resource overhead for large test suites

**🔧 Complexity**
- Docker dependency
- Additional setup and configuration
- Learning curve for advanced patterns

**💰 Resource Usage**
- Memory and CPU overhead
- Disk space for container images
- Network resources for container communication

## 🎯 When to Use Test Containers

### ✅ Perfect For

- **Database Integration Tests**: Testing actual database operations, queries, and constraints
- **Microservices**: Testing service-to-service communication with real dependencies
- **Message Queues**: Testing producers/consumers with real brokers (RabbitMQ, Kafka)
- **Web Applications**: Testing against real browsers with Selenium
- **API Testing**: Testing against real API servers
- **Complex Business Logic**: Where mocking would be too complex or error-prone

### ❌ Not Ideal For

- **Unit Tests**: Where fast, isolated testing is needed
- **Simple CRUD**: Where in-memory databases suffice
- **Performance Tests**: Where you need to test against production-scale data
- **CI/CD Pipelines**: Where startup time is critical (though still viable)

## 📚 Progressive Implementation Guide

### Level 1: Basic Setup

**1. Add Dependencies**
```bash
go get github.com/testcontainers/testcontainers-go/modules/mongodb
```

**2. Basic Test Structure**
```go
func TestBasicMongoDB(t *testing.T) {
    ctx := context.Background()

    // Start container
    mongoContainer, err := mongodb.Run(ctx, "mongo:7")
    if err != nil {
        t.Fatal(err)
    }
    defer mongoContainer.Terminate(ctx)

    // Get connection
    uri, err := mongoContainer.ConnectionString(ctx)
    if err != nil {
        t.Fatal(err)
    }

    // Your test logic here
    // ...
}
```

### Level 2: Test Helper Pattern

**Create Reusable Helpers**
```go
type MongoDBTestHelper struct {
    Container *mongodb.MongoDBContainer
    Config    *config.DatabaseConfig
    Repo      repository.UserRepository
}

func SetupMongoDBContainer(ctx context.Context, t *testing.T) *MongoDBTestHelper {
    // Implementation as shown in testhelpers/mongodb.go
}
```

**Usage in Tests**
```go
func TestWithHelper(t *testing.T) {
    ctx := context.Background()
    helper := testhelpers.SetupMongoDBContainer(ctx, t)
    defer helper.Cleanup(ctx, t)

    // Use helper.Repo for testing
}
```

### Level 3: Configuration Management

**Environment-Based Config**
```go
type DatabaseConfig struct {
    URI        string
    Database   string
    Collection string
    Username   string
    Password   string
}

func LoadDatabaseConfig() *DatabaseConfig {
    return &DatabaseConfig{
        URI:        getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
        Database:   getEnvOrDefault("MONGO_DATABASE", "testdb"),
        Collection: getEnvOrDefault("MONGO_COLLECTION", "users"),
        // ...
    }
}
```

### Level 4: Advanced Testing Patterns

**Parallel Testing**
```go
func TestParallelOperations(t *testing.T) {
    t.Parallel() // Run tests in parallel

    ctx := context.Background()
    helper := testhelpers.SetupMongoDBContainer(ctx, t)
    defer helper.Cleanup(ctx, t)

    // Each test gets its own container
}
```

**Integration Testing**
```go
func TestFullApplicationFlow(t *testing.T) {
    ctx := context.Background()
    helper := testhelpers.SetupMongoDBContainer(ctx, t)
    defer helper.Cleanup(ctx, t)

    // Test the full application stack
    server, err := NewServer(helper.Repo)
    require.NoError(t, err)

    err = server.RunDemo()
    require.NoError(t, err)
}
```

### Level 5: CI/CD Integration

**GitHub Actions Example**
```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: go test -v ./...
```

## 🏗️ Project Structure

```
.
├── config/
│   └── config.go              # Configuration management
├── demo/
│   └── demo.go                # Demo operations
├── models/
│   └── user.go                # Data models
├── repository/
│   ├── interface.go           # Repository contracts
│   ├── mongodb.go             # MongoDB implementation
│   └── mongodb_test.go        # Repository tests
├── testhelpers/
│   └── mongodb.go             # Test utilities
├── server.go                  # Application server
├── server_test.go             # Integration tests
├── main.go                    # Application entry point
└── docker-compose.yml         # Local development setup
```

## 🚀 Quick Start

**1. Prerequisites**
- Go 1.21+
- Docker
- Make sure Docker daemon is running

**2. Clone and Setup**
```bash
git clone <repository-url>
cd testcontainers-mongo-demo
go mod download
```

**3. Run Tests**
```bash
# Run all tests
go test -v ./...

# Run specific package tests
go test -v ./repository

# Run with race detection
go test -race -v ./...
```

**4. Run Application**
```bash
# Start local MongoDB (optional)
docker-compose up -d

# Run the demo
go run main.go
```

## 🔧 Advanced Patterns

### Custom Container Configuration

```go
func TestWithCustomConfig(t *testing.T) {
    ctx := context.Background()

    container, err := mongodb.Run(ctx,
        "mongo:7",
        mongodb.WithUsername("testuser"),
        mongodb.WithPassword("testpass"),
        mongodb.WithDatabase("testdb"),
    )
    // ...
}
```

### Multiple Containers

```go
func TestMultiService(t *testing.T) {
    ctx := context.Background()

    // MongoDB
    mongoContainer, err := mongodb.Run(ctx, "mongo:7")
    // ...

    // Redis
    redisContainer, err := redis.Run(ctx, "redis:7")
    // ...

    // Test service interactions
}
```

### Network Management

```go
func TestNetworkedServices(t *testing.T) {
    ctx := context.Background()

    network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
        NetworkRequest: testcontainers.NetworkRequest{
            Name: "test-network",
        },
    })
    // ...

    mongoContainer, err := mongodb.Run(ctx, "mongo:7",
        testcontainers.WithNetwork([]string{"test-network"}, network),
    )
    // ...
}
```

## 📋 Best Practices

### 🧪 Test Organization

**1. Test Naming**
```go
func TestUserRepository_Create(t *testing.T)    // Unit test
func TestServer_RunDemo(t *testing.T)           // Integration test
func TestInitializeDatabase(t *testing.T)       // End-to-end test
```

**2. Test Categories**
- **Unit Tests**: Test individual functions with mocks
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete workflows

### 🛠️ Resource Management

**1. Proper Cleanup**
```go
func TestWithCleanup(t *testing.T) {
    helper := testhelpers.SetupMongoDBContainer(ctx, t)
    defer helper.Cleanup(ctx, t) // Always cleanup

    // Test logic
}
```

**2. Timeout Management**
```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    helper := testhelpers.SetupMongoDBContainer(ctx, t)
    defer helper.Cleanup(ctx, t)
}
```

### 📊 Assertions & Validation

**1. Use Testify for Better Assertions**
```go
require.NoError(t, err)
require.Equal(t, expected, actual)
require.Len(t, users, 1)
```

**2. Validate Important Behaviors**
```go
// Test that ID is generated
require.False(t, user.ID.IsZero())

// Test data integrity
retrieved, err := repo.GetByID(ctx, user.ID)
require.NoError(t, err)
require.Equal(t, "John", retrieved.Name)
```

## 🔍 Troubleshooting

### Common Issues

**1. Container Startup Failures**
```
Error: container startup timeout
```
**Solutions:**
- Ensure Docker daemon is running
- Check available resources (memory, disk space)
- Try different container versions
- Add startup timeouts

**2. Connection Issues**
```
Error: connection refused
```
**Solutions:**
- Verify container is ready before connecting
- Use `mongoContainer.ConnectionString(ctx)`
- Add retry logic for connection attempts

**3. Test Flakiness**
```
Error: tests pass locally but fail in CI
```
**Solutions:**
- Use consistent container versions
- Ensure proper cleanup between tests
- Add retry logic for transient failures
- Check CI environment Docker configuration

### Debug Tips

**1. Enable Debug Logging**
```go
import "github.com/testcontainers/testcontainers-go"

testcontainers.Logger = log.New(os.Stdout, "", log.LstdFlags)
```

**2. Inspect Container State**
```go
// Check if container is running
state, err := mongoContainer.State(ctx)
if err != nil {
    t.Logf("Container state error: %v", err)
}
```

**3. Manual Container Inspection**
```bash
# List running containers
docker ps

# Check container logs
docker logs <container-id>

# Connect to container
docker exec -it <container-id> mongo
```

## 📚 Additional Resources

- [Testcontainers Go Documentation](https://golang.testcontainers.org/)
- [Testcontainers Official Website](https://testcontainers.org/)
- [MongoDB Go Driver](https://docs.mongodb.com/drivers/go/)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/TableDrivenTests)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.