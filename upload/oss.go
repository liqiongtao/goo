package upload

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/liqiongtao/goo/utils"
	"log"
	"path"
)

var OSS *gooOSS

type OSSConfig struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket"`
	Domain          string `yaml:"domain"`
}

func InitOSS(config OSSConfig) {
	var err error
	OSS, err = NewOSS(config)
	if err != nil {
		log.Panic(err.Error())
	}
}

func NewOSS(config OSSConfig) (*gooOSS, error) {
	o := &gooOSS{
		Config: config,
	}

	client, err := o.getClient()
	if err != nil {
		return nil, err
	}

	o.Client = client

	bucket, err := o.getBucket()
	if err != nil {
		return nil, err
	}

	o.Bucket = bucket

	return o, nil
}

type gooOSS struct {
	Config OSSConfig
	Client *oss.Client
	Bucket *oss.Bucket
}

func (o *gooOSS) Upload(filename string, body []byte) (string, error) {
	md5str := utils.MD5(body)
	filename = fmt.Sprintf("%s/%s/%s%s", md5str[0:2], md5str[2:4], md5str[8:24], path.Ext(filename))

	if err := o.Bucket.PutObject(filename, bytes.NewReader(body)); err != nil {
		return "", err
	}

	if o.Config.Domain != "" {
		return o.Config.Domain + filename, nil
	}

	url := "https://" + o.Config.Bucket + "." + o.Config.Endpoint + "/" + filename
	return url, nil
}

func (o *gooOSS) getClient() (*oss.Client, error) {
	return oss.New(o.Config.Endpoint, o.Config.AccessKeyId, o.Config.AccessKeySecret)
}

func (o *gooOSS) getBucket() (*oss.Bucket, error) {
	return o.Client.Bucket(o.Config.Bucket)
}
