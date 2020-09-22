package upload

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/liqiongtao/goo/utils"
	"io/ioutil"
	"os"
	"path"
)

// -----------------------------------------------
// 定义main文件，然后调用ossScript()
// go build -ldflags "-s -w" -o oss
// -----------------------------------------------

var (
	endpoint        = "oss-cn-beijing.aliyuncs.com"
	accessKeyId     = ""
	accessKeySecret = ""
	bucketName      = ""
)

func ossScript(endpoint, accessKeyId, accessKeySecret, bucketName string) {
	if total := len(os.Args); total < 2 {
		fmt.Println("请选择上传文件", total, os.Args)
		return
	}

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("[oss-client]", err.Error())
		return
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("[oss-bucket]", err.Error())
		return
	}

	filename := os.Args[1]

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("[read-file]", err.Error())
		return
	}

	md5str := utils.MD5(body)
	filename = fmt.Sprintf("%s/%s/%s%s", md5str[0:2], md5str[2:4], md5str[8:24], path.Ext(filename))

	if err := bucket.PutObject(filename, bytes.NewReader(body)); err != nil {
		fmt.Println("[oss-upload]", err.Error())
		return
	}

	fmt.Println("https://" + bucketName + "." + endpoint + "/" + filename)
}
