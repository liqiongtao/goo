package goo

import "testing"

func TestNewGRPCServer(t *testing.T) {
	s := NewGRPCServer(15001,
		GRPCServiceName("test"),
	)
	s.Serve()
}
