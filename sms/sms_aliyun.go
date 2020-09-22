package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type AliyunConfig struct {
	Region       string `yaml:"region"`
	Appid        string `yaml:"appid"`
	Secret       string `yaml:"secret"`
	SignName     string `yaml:"sign_name"`
	TemplateCode string `yaml:"template_code"`
}

type gooAliyun struct {
	conf AliyunConfig
	*sdk.Client
}

func Aliyun(conf AliyunConfig) *gooAliyun {
	ali := &gooAliyun{
		conf: conf,
	}
	ali.Client, _ = sdk.NewClientWithAccessKey(ali.conf.Region, ali.conf.Appid, ali.conf.Secret)
	return ali
}

func (ali *gooAliyun) Send(mobile, action string) (string, error) {
	code := ali.getSmsCode(mobile, action)

	request := requests.NewCommonRequest()

	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"

	request.QueryParams["RegionId"] = ali.conf.Region
	request.QueryParams["PhoneNumbers"] = mobile
	request.QueryParams["SignName"] = ali.conf.SignName
	request.QueryParams["TemplateCode"] = ali.conf.TemplateCode
	request.QueryParams["TemplateParam"] = fmt.Sprintf("{\"code\": %s}", code)

	response, err := ali.Client.ProcessCommonRequest(request)
	if err != nil {
		log.Println("[aliyun-sms]", "[send-error]", request.String(), err.Error())
		return "", err
	}

	rsp := map[string]string{}
	if err := json.Unmarshal(response.GetHttpContentBytes(), &rsp); err != nil {
		log.Println("[aliyun-sms]", "[send-error]", request.String(), err.Error())
		return "", err
	}

	if rsp["Code"] != "OK" {
		log.Println("[aliyun-sms]", "[send-error]", request.String(), rsp)
		return "", errors.New(rsp["Message"])
	}

	__cache.set(ali.conf.Appid, mobile, action, code, expireIn)

	return code, nil
}

func (ali *gooAliyun) Verify(mobile, action, code string) error {
	rst := __cache.get(ali.conf.Appid, mobile, action)
	if rst == nil {
		return errors.New("验证码无效")
	}
	if rst.ExpireOut < time.Now().Unix() {
		return errors.New("验证码过期")
	}
	if rst.Code != code {
		return errors.New("验证码错误")
	}
	return nil
}

func (ali *gooAliyun) getSmsCode(mobile, action string) string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(8)+1) +
		strconv.Itoa(rand.Intn(8)+1) +
		strconv.Itoa(rand.Intn(8)+1) +
		strconv.Itoa(rand.Intn(8)+1) +
		strconv.Itoa(rand.Intn(8)+1) +
		strconv.Itoa(rand.Intn(8)+1)
}
