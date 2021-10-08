package goo

import (
	"bytes"
	"io"
	"mime/multipart"
)

const (
	CONTENT_TYPE_XML  = "application/xml; charset=utf-8"
	CONTENT_TYPE_JSON = "application/json; charset=utf-8"
	CONTENT_TYPE_FORM = "application/x-www-form-urlencoded; charset=utf-8"
)

func NewRequest() *Request {
	return &Request{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE_FORM,
		},
	}
}

func NewTlsRequest(caCrtFile, clientCrtFile, clientKeyFile string) *Request {
	return &Request{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE_FORM,
		},
		Tls: &Tls{
			CaCrtFile:     caCrtFile,
			ClientCrtFile: clientCrtFile,
			ClientKeyFile: clientKeyFile,
		},
	}
}

func Get(url string) ([]byte, error) {
	return NewRequest().Get(url)
}

func Post(url string, data []byte) ([]byte, error) {
	return NewRequest().Post(url, data)
}

func PostJson(url string, data []byte) ([]byte, error) {
	return NewRequest().JsonContentType().Post(url, data)
}

func Put(url string, data []byte) ([]byte, error) {
	return NewRequest().Put(url, data)
}

func Upload(url, field, fileName string, f io.Reader, data map[string]string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(field, fileName)
	if err != nil {
		Log.Error(err.Error())
		return nil, err
	}
	if _, err = io.Copy(part, f); err != nil {
		Log.Error(err.Error())
		return nil, err
	}

	for k, v := range data {
		writer.WriteField(k, v)
	}

	if err = writer.Close(); err != nil {
		Log.Error(err.Error())
		return nil, err
	}

	request := NewRequest()
	request.SetHearder("Content-Type", writer.FormDataContentType())
	return request.Do("POST", url, body)
}
