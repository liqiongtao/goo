package goo

import (
	"encoding/json"
)

type Response struct {
	Status  int           `json:"status,omitempty"`
	Code    int           `json:"code,omitempty"`
	Message string        `json:"message,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
	ErrMsg  []interface{} `json:"-"`
}

func (rsp *Response) String() string {
	buf, err := json.Marshal(rsp)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

func Success(data interface{}) *Response {
	if data == nil {
		data = map[string]interface{}{}
	}
	return &Response{
		Status:  1,
		Code:    200,
		Message: "ok",
		Data:    data,
	}
}

func Error(code int, message string, v ...interface{}) *Response {
	return &Response{
		Status:  0,
		Code:    code,
		Message: message,
		Data:    map[string]string{},
		ErrMsg:  v,
	}
}
