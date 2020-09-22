package upload

import (
	"github.com/gin-gonic/gin"
	"github.com/liqiongtao/goo/utils"
	"io/ioutil"
	"os"
	"path"
)

var (
	Local      = gooLocal{}
	upload_dir = "static/"
)

type gooLocal struct {
}

func (l gooLocal) Upload(c *gin.Context) (string, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	md5str := utils.MD5(data)
	fpath := md5str[0:2] + "/" + md5str[2:4] + "/"

	if err := os.MkdirAll(upload_dir+fpath, 0755); err != nil {
		return "", err
	}

	fname := fpath + md5str[8:24] + path.Ext(header.Filename)

	fw, err := os.Create(upload_dir + fname)
	if err != nil {
		return "", err
	}
	defer fw.Close()

	if _, err := fw.Write(data); err != nil {
		return "", err
	}

	return fname, nil
}
