package testhelpers

import (
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StartTestGRPCServer is a generic Test Server
// similar to the httptest.NewServer implementation
func StartTestGRPCServer(t *testing.T, register func(*grpc.Server)) (string, *grpc.ClientConn, func()) {
	lis, err := net.Listen("tcp", "localhost:0") // Listen on random port
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	register(server) // Let caller register their services

	go func() {
		if err := server.Serve(lis); err != nil {
			t.Logf("server closed: %v", err)
		}
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	cleanup := func() {
		server.Stop()
		if err := lis.Close(); err != nil {
			t.Logf("server closed: %v", err)
		}
		if err := conn.Close(); err != nil {
			t.Logf("server closed: %v", err)
		}
	}

	return lis.Addr().String(), conn, cleanup
}
