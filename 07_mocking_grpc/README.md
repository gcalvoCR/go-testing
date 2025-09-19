# gRPC Mocking Example in Go

This project demonstrates how to create a basic gRPC service in Go and, more importantly, how to mock gRPC servers for testing external dependencies. This is particularly useful when you want to test your client code without relying on actual network calls or external services.

## Project Structure

```
.
├── hello.proto              # Protocol buffer definition
├── hellopb/                 # Generated Go code from proto
│   ├── hello.pb.go
│   └── hello_grpc.pb.go
├── server/                  # Real gRPC server implementation
│   └── server.go
├── client/                  # gRPC client implementation
│   ├── client.go
│   └── client_test.go       # Tests using mocked server
├── mocks/                   # Mock implementations
│   └── server_mock.go
├── testhelpers/             # Testing utilities
│   └── helpers.go
├── makefile                 # Build automation
├── go.mod                   # Go module file
└── go.sum                   # Go dependencies
```

## Basic gRPC Example

### Running the Server

```bash
go run server/server.go
```

The server will start on `localhost:50051`.

### Running the Client

In another terminal:

```bash
go run client/client.go
```

This will connect to the server and print a greeting message.

## Testing with Mocked gRPC Server

The main focus of this project is demonstrating how to mock gRPC servers for testing.

### Key Components for Testing

1. **Mock Server** (`mocks/server_mock.go`):
   - Implements the same interface as the real server
   - Returns predictable responses for testing
   - No network dependencies

2. **Test Helpers** (`testhelpers/helpers.go`):
   - `StartTestGRPCServer`: Creates an in-memory gRPC server for testing
   - Handles server lifecycle (start/stop)
   - Provides client connection for tests

3. **Test Example** (`client/client_test.go`):
   - Uses the mock server instead of a real network server
   - Tests client logic in isolation
   - Fast and reliable

### Running Tests

```bash
make test
# or
go test ./client/ -v
```

### How Mocking Works

1. **Test Server Setup**:
   ```go
   address, conn, cleanup := testhelpers.StartTestGRPCServer(t, func(s *grpc.Server) {
       hellopb.RegisterGreeterServer(s, &mocks.ServerMock{})
   })
   defer cleanup()
   ```

2. **Mock Implementation**:
   ```go
   type ServerMock struct {
       pb.UnimplementedGreeterServer
   }

   func (s *ServerMock) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
       return &pb.HelloReply{Message: "Salut " + in.GetName()}, nil
   }
   ```

3. **Client Testing**:
   ```go
   client := pb.NewGreeterClient(conn)
   response, err := client.SayHello(context.Background(), in)
   // Assert expected behavior
   ```

## Benefits of Mocking gRPC Services

- **Isolation**: Test client logic without external dependencies
- **Speed**: No network calls means faster test execution
- **Reliability**: Tests don't fail due to network issues or server downtime
- **Control**: Predictable responses for edge case testing
- **CI/CD Friendly**: No need to spin up external services in CI pipelines

## Generating Protocol Buffer Code

If you modify `hello.proto`, regenerate the Go code:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       hello.proto
```

## Dependencies

- Go 1.19+
- Protocol Buffers compiler (protoc)
- Go plugins for protoc

Install protoc plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Extending the Example

You can extend this pattern to:
- Add more RPC methods to the service
- Create multiple mock implementations for different test scenarios
- Add middleware for authentication/authorization testing
- Implement streaming RPCs with mocked streams

This approach ensures your tests remain fast, reliable, and focused on the code under test rather than external dependencies.