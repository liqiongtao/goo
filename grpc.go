package goo

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"runtime/debug"
)

type GRPCServer struct {
	ServiceName string
	Server      *grpc.Server
	lis         net.Listener
	consul      *Consul
}

func NewGRPCServer(port int64, serviceName string, consul *Consul) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(grpcInterceptor),
	}

	return &GRPCServer{
		ServiceName: serviceName,
		Server:      grpc.NewServer(opts...),
		lis:         lis,
		consul:      consul,
	}, nil
}

func (s *GRPCServer) Serve() error {
	go func() {
		log.Println(fmt.Sprintf("server running %s, pid=%d", s.lis.Addr().String(), os.Getpid()))
	}()
	s.registerHealthServer()
	s.registerToConsul()
	return s.Server.Serve(s.lis)
}

func (s *GRPCServer) registerHealthServer() {
	grpc_health_v1.RegisterHealthServer(s.Server, &GRPCHealth{})
}

func (s *GRPCServer) registerToConsul() {
	if s.ServiceName == "" || s.consul == nil {
		return
	}
	if err := s.consul.ServiceRegister(s.ServiceName); err != nil {
		log.Fatalln(err.Error())
	}
}

func grpcInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			Log.WithField("grpc-method", info.FullMethod).
				WithField("grep-request", req).
				WithField("error-stack", string(debug.Stack())).
				Error(fmt.Sprintf("%v", e))
		}
	}()
	rsp, err = handler(ctx, req)
	if err == nil {
		Log.WithField("grpc-method", info.FullMethod).
			WithField("grpc-request", req).
			WithField("grpc-response", rsp).
			Info()
	} else {
		Log.WithField("grpc-method", info.FullMethod).
			WithField("grpc-request", req).
			WithField("grpc-response", rsp).
			WithField("error-stack", string(debug.Stack())).
			Error(err.Error())
	}
	return
}

type GRPCHealth struct{}

func (GRPCHealth) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (GRPCHealth) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}
