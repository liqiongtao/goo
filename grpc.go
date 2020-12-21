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
)

type GRPCServer struct {
	Server  *grpc.Server
	options map[string]Option
}

func NewGRPCServer(port int64, opts ...Option) *GRPCServer {
	s := &GRPCServer{
		Server: grpc.NewServer(grpc_middleware.WithUnaryServerChain(GRPCInterceptor)),
		options: map[string]Option{
			grpcServerPort: NewOption(grpcServerPort, port),
		},
	}
	for _, opt := range opts {
		opt.Apply(s.options)
	}
	return s
}

func (s *GRPCServer) address() string {
	return fmt.Sprintf(":%d", s.options[grpcServerPort].Int64())
}

func (s *GRPCServer) serviceName() string {
	return s.options[grpcServiceName].String()
}

func (s *GRPCServer) consul() *Consul {
	if s.options[grpcConsul].Value == nil {
		return nil
	}
	return s.options[grpcConsul].Value.(*Consul)
}

func (s *GRPCServer) Serve() error {
	lis, err := net.Listen("tcp", s.address())
	if err != nil {
		return err
	}
	s.registerHealthServer()
	s.registerToConsul()
	AsyncFunc(func() {
		log.Println(fmt.Sprintf("server running %s, pid=%d", lis.Addr().String(), os.Getpid()))
	})
	AsyncFunc(func() {
		for {
			select {
			case <-Context.Done():
				os.Exit(0)
			}
		}
	})
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
		Log.Fatal(err.Error())
	}
}

func GRPCInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			Log.WithField("grpc-method", info.FullMethod).
				WithField("grep-request", req).
				Trace().
				Error(fmt.Sprintf("%v", e))
		}
	}()
	rsp, err = handler(ctx, req)
	if info.FullMethod != "/grpc.health.v1.Health/Check" {
		return
	}
	if err == nil {
		Log.WithField("grpc-method", info.FullMethod).
			WithField("grpc-request", req).
			WithField("grpc-response", rsp).
			Info()
	} else {
		Log.WithField("grpc-method", info.FullMethod).
			WithField("grpc-request", req).
			WithField("grpc-response", rsp).
			Trace().
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
