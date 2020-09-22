package goo

import (
	"fmt"
	"github.com/go-playground/validator"
	"strings"
)

func ValidationMessage(err error, msgs map[string]string) string {
	for _, i := range err.(validator.ValidationErrors) {
		key := fmt.Sprintf("%s_%s", strings.ToLower(i.Field()), strings.ToLower(i.Tag()))
		if msg, ok := msgs[key]; ok {
			return msg
		}
		msg := fmt.Sprintf("%s %s", i.Field(), i.Tag())
		return msg
	}
	return err.Error()
}
