package goo

import (
	"fmt"
	"testing"
)

func TestNewConsul(t *testing.T) {
	c := NewConsul("", "", "")
	cc, _ := c.Client()
	fmt.Println(cc.KV().Get("", nil))
}
