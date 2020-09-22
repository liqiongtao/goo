package goo

import (
	"encoding/json"
	"log"
)

type Response struct {
	Status  int         `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (rsp *Response) String() string {
	buf, err := json.Marshal(rsp)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

func Success(data interface{}) *Response {
	return &Response{
		Status:  1,
		Code:    200,
		Message: "ok",
		Data:    data,
	}
}

func Error(code int, message string, v ...interface{}) *Response {
	log.Println(v...)
	return &Response{
		Status:  0,
		Code:    code,
		Message: message,
		Data:    map[string]string{},
	}
}
