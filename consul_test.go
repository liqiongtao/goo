package goo

import (
	"fmt"
	"testing"
)

func TestConsulKV(t *testing.T) {
	c := NewConsul("", "", "")
	cc, _ := c.Client()
	fmt.Println(cc.KV().Get("", nil))
}

func TestConsulServiceRegister(t *testing.T) {
	c := NewConsul("", "", "")
	c.ServiceRegister("")
}

func TestConsulServiceDeregister(t *testing.T) {
	c := NewConsul("", "", "")
	c.ServiceDeregister("")
}
