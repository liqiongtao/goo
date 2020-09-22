package captcha

import (
	"errors"

	"github.com/mojocn/base64Captcha"
)

func Captcha(width, height int) map[string]interface{} {
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}

	var configDigit = base64Captcha.ConfigDigit{
		Height:     height,
		Width:      width,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 6,
	}

	id, cap := base64Captcha.GenerateCaptcha("", configDigit)
	base64image := base64Captcha.CaptchaWriteToBase64Encoding(cap)

	return map[string]interface{}{
		"id":          id,
		"base64image": base64image,
	}
}

func CaptchaVerify(captchaId, captchaCode string) error {
	if captchaId == "" {
		return errors.New("图片验证码ID为空")
	}

	if captchaCode == "" {
		return errors.New("图片验证码为空")
	}

	verifyResult := base64Captcha.VerifyCaptcha(captchaId, captchaCode)
	if !verifyResult {
		return errors.New("图片验证码错误")
	}

	return nil
}
