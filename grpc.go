package goo

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

type GrpcServer struct {
	serviceName string
	Server      *grpc.Server
	lis         net.Listener
	consul      *Consul
}

func NewGrpcServer(port int64, serviceName string, consul *Consul) (*GrpcServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	return &GrpcServer{
		serviceName: serviceName,
		Server:      grpc.NewServer(),
		lis:         lis,
		consul:      consul,
	}, nil
}

func (s *GrpcServer) Serve() error {
	s.registerHealthServer()
	s.registerToConsul()
	return s.Server.Serve(s.lis)
}

func (s *GrpcServer) registerHealthServer() {
	grpc_health_v1.RegisterHealthServer(s.Server, &Health{})
}

func (s *GrpcServer) registerToConsul() {
	if err := s.consul.ServiceRegister(s.serviceName); err != nil {
		log.Fatalln(err.Error())
	}
}

type Health struct{}

func (Health) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (Health) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}
