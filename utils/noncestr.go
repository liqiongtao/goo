package utils

import (
	"strings"
	"time"
)

func NonceStr() string {
	return strings.ToLower(Id2Code(time.Now().UnixNano()))
}
