package goo

import (
	"log"
	"testing"
)

func TestNewGRPCServer(t *testing.T) {
	s := NewGRPCServer(15001,
		GRPCServiceName("test"),
	)
	if err := s.Serve(); err != nil {
		log.Println(err.Error())
	}
}
