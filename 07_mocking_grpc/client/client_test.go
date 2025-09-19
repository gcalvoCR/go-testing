package main

import (
	"context"
	"fmt"
	"grpc-example/hellopb"
	pb "grpc-example/hellopb"
	"grpc-example/mocks"
	"grpc-example/testhelpers"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func Test_SayHello(t *testing.T) {
	address, conn, cleanup := testhelpers.StartTestGRPCServer(t, func(s *grpc.Server) {
		hellopb.RegisterGreeterServer(s, &mocks.ServerMock{})
	})
	defer cleanup()

	client := pb.NewGreeterClient(conn)

	in := &pb.HelloRequest{
		Name: "Gabriel",
	}

	// Real test without making call
	response, err := client.SayHello(context.Background(), in)
	require.NoError(t, err)
	require.Equal(t, response.Message, "Salut Gabriel")

	fmt.Println("The server run in the following address:", address)

}
