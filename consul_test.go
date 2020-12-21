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
	c := NewConsul("http://dc.weflys.com", "xz", "xz527")
	c.ServiceRegister("xz/services/base-auth-test")
}

func TestConsulServiceDeregister(t *testing.T) {
	c := NewConsul("http://dc.weflys.com", "xz", "xz527")
	c.ServiceDeregister("base_auth_test")
}
