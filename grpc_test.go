package goo

import (
	"log"
	"testing"
)

func TestNewGRPCServer(t *testing.T) {
	s := NewGRPCServer(15001)
	if err := s.Serve(); err != nil {
		log.Println(err.Error())
	}
}

func TestNewGRPCServer2(t *testing.T) {
	s := NewGRPCServer(15001,
		GRPCServiceName("test"),
		GRPCConsul("127.0.0.1:18001", "test", "123456"),
	)
	if err := s.Serve(); err != nil {
		log.Println(err.Error())
	}
}
