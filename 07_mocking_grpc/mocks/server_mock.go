package mocks

import (
	"context"
	pb "grpc-example/hellopb"
)

type ServerMock struct {
	pb.UnimplementedGreeterServer
}

func (s *ServerMock) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Salut " + in.GetName()}, nil
}
