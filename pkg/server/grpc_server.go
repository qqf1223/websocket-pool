package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	srv         *grpc.Server
	Name        string
	Addr        string
	Port        int
	ServiceOpts []grpc.ServerOption
	RegisterFn  func(srv *grpc.Server)
}

func (s *GRPCServer) Run() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.srv = grpc.NewServer(s.ServiceOpts...)
	if s.RegisterFn != nil {
		s.RegisterFn(s.srv)
	}

	log.Printf("[INFO] GRPC server start, listen: %s\n", s.Addr)
	if err := s.srv.Serve(l); err != nil {
		return err
	}
	return nil
}
