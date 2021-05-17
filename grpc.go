package goo

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"path"
	"runtime"
	"strings"
	"time"
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
	s.registerHealthServer()
	s.registerToConsul()

	g := NewGRPCGraceful("tcp", s.address(), s.Server)
	return g.Serve()
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
	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		return
	}
	lg := Log.WithField("grpc-method", info.FullMethod).WithField("grpc-request", req)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key, val := range md {
			lg.WithField(key, val)
		}
	}
	defer func() {
		if e := recover(); e != nil {
			lg.Error(fmt.Sprintf("%v", e))
		}
	}()
	rsp, err = handler(ctx, req)
	lg.WithField("grpc-response", rsp)
	if err == nil {
		lg.Info()
		return
	}
	lg.Error(err.Error())
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

func GRPCContext(c *gin.Context) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 8*time.Second)
	md := metadata.New(map[string]string{})
	if c != nil {
		md.Set("request-id", fmt.Sprintf("%d", c.GetInt("__trace_id")))
		md.Set("request-server", c.GetString("__server_name"))
		_, file, line, _ := runtime.Caller(1)
		file = strings.Replace(file, path.Dir(c.GetString("__base_dir"))+"/", "", -1)
		md.Set("request-file", fmt.Sprintf("%s %dL", file, line))
	}
	return metadata.NewOutgoingContext(ctx, md)
}
