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
	Server  *grpc.Server
	options map[string]Option
}

func NewGRPCServer(port int, opts ...Option) *GRPCServer {
	s := &GRPCServer{options: map[string]Option{
		"port": NewOption(grpcServerPort, port),
	}}
	for _, opt := range opts {
		opt.Apply(s.options)
	}
	return s
}

func (s *GRPCServer) address() string {
	return fmt.Sprintf(":%d", s.options[grpcServerPort].Int())
}

func (s *GRPCServer) serviceName() string {
	return s.options[grpcServiceName].String()
}

func (s *GRPCServer) consul() *Consul {
	return s.options[grpcConsul].Value.(*Consul)
}

func (s *GRPCServer) Serve() error {
	lis, err := net.Listen("tcp", s.address())
	if err != nil {
		return err
	}
	s.Server = grpc.NewServer(grpc_middleware.WithUnaryServerChain(GRPCInterceptor))
	s.registerHealthServer()
	s.registerToConsul()
	go func() {
		log.Println(fmt.Sprintf("server running %s, pid=%d", lis.Addr().String(), os.Getpid()))
	}()
	return s.Server.Serve(lis)
}

func (s *GRPCServer) registerHealthServer() {
	grpc_health_v1.RegisterHealthServer(s.Server, &GRPCHealth{})
}

func (s *GRPCServer) registerToConsul() {
	if s.serviceName() == "" || s.consul() == nil {
		return
	}
	if err := s.consul().ServiceRegister(s.serviceName()); err != nil {
		log.Fatalln(err.Error())
	}
}

func GRPCInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, err error) {
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
