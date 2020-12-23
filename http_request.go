package goo

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Headers map[string]string
	Tls     *Tls
}

func (r *Request) SetHearder(name, value string) *Request {
	r.Headers[name] = value
	return r
}

func (r *Request) SetContentType(contentType string) *Request {
	r.SetHearder("Content-Type", contentType)
	return r
}

func (r *Request) JsonContentType() *Request {
	r.SetHearder("Content-Type", CONTENT_TYPE_JSON)
	return r
}

func (r *Request) getClient() *http.Client {
	client := &http.Client{}
	if r.Tls != nil {
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(r.Tls.CaCrt())
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{r.Tls.ClientCrt()},
			},
		}
	}
	return client
}

func (r *Request) Do(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		Log.Error(err.Error())
		return nil, err
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	rsp, err := r.getClient().Do(req)
	if err != nil {
		Log.Error( err.Error())
		return nil, err
	}

	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		Log.Error(err.Error())
		return nil, err
	}

	return buf, nil
}

func (r *Request) Get(url string) ([]byte, error) {
	return r.Do("GET", url, nil)
}

func (r *Request) Post(url string, data []byte) ([]byte, error) {
	return r.Do("POST", url, bytes.NewReader(data))
}
