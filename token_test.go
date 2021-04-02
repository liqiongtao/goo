package goo

import (
	"log"
	"testing"
)

func TestCreateToken(t *testing.T) {
	tokenStr, err := CreateToken("1234567890", 100)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(tokenStr)
}

func TestParseToken(t *testing.T) {
	tokenStr := "R7aM0LcLlxXIAXTDaLx47Ccp3i6dULm+yVLxN4xemkTLCYfkPFAa41wsUKnTRqQN4eufXumA9htEXKBFMwigKQf+xmkYD/POywnGKXiUt6E="
	token, err := ParseToken(tokenStr, "1234567890")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(token)
}
