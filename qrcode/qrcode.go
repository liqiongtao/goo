package qrcode

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

const (
	default_size = 256
)

func New(url string) *gooQRCode {
	return &gooQRCode{
		Url: url,
	}
}

type gooQRCode struct {
	Size int
	Url  string
}

func (qr *gooQRCode) Get() ([]byte, error) {
	if qr.Url == "" {
		return nil, errors.New("url 为空")
	}

	if qr.Size == 0 {
		qr.Size = default_size
	}

	qr.Url, _ = url.QueryUnescape(qr.Url)

	png, err := qrcode.Encode(qr.Url, qrcode.Medium, qr.Size)
	if err != nil {
		return nil, err
	}

	return png, nil
}

func (qr *gooQRCode) Base64Image() (string, error) {
	png, err := qr.Get()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(png)), nil
}

func (qr *gooQRCode) Output(c gin.Context) error {
	png, err := qr.Get()
	if err != nil {
		return err
	}

	c.Header("Content-Type", "image/png")
	c.Writer.Write(png)

	return nil
}
