package main

import (
	"context"
	"testcontainers-mongo-demo/demo"
	"testcontainers-mongo-demo/repository"
)

type Server struct {
	repo repository.UserRepository
}

func NewServer(repo repository.UserRepository) (*Server, error) {
	return &Server{
		repo: repo,
	}, nil
}

func (s *Server) RunDemo() error {
	ctx := context.Background()
	return demo.RunUserDemo(ctx, s.repo)
}
