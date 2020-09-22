package goo

import (
	"crypto/tls"
	"io/ioutil"
	"log"
)

type Tls struct {
	CaCrtFile     string
	ClientCrtFile string
	ClientKeyFile string
}

func (this *Tls) CaCrt() []byte {
	if this.CaCrtFile == "" {
		return caCert
	}
	bts, err := ioutil.ReadFile(this.CaCrtFile)
	if err != nil {
		log.Println(err.Error())
	}
	return bts
}

func (this *Tls) ClientCrt() tls.Certificate {
	crt, err := tls.LoadX509KeyPair(this.ClientCrtFile, this.ClientKeyFile)
	if err != nil {
		log.Println(err.Error())
	}
	return crt
}
