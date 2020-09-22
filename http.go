package goo

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
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

func Upload(url, field, file string, data map[string]string) ([]byte, error) {
	fh, err := os.Open(file)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer fh.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(field, filepath.Base(file))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if _, err = io.Copy(part, fh); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	for k, v := range data {
		writer.WriteField(k, v)
	}

	if err = writer.Close(); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	request := NewRequest()
	request.SetHearder("Content-Type", writer.FormDataContentType())
	return request.Do("POST", url, body)
}
