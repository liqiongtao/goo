package goo

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strings"
)

func (gooUtil) MD5(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func (gooUtil) SHA1(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func (gooUtil) SHA256(buf, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func (gooUtil) HMacMd5(buf, key []byte) string {
	h := hmac.New(md5.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func (gooUtil) HMacSha1(buf, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func (gooUtil) HMacSha256(buf, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func (gooUtil) Base64Encode(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

func (gooUtil) Base64Decode(str string) []byte {
	var count = (4 - len(str)%4) % 4
	str += strings.Repeat("=", count)
	buf, _ := base64.StdEncoding.DecodeString(str)
	return buf
}

func (ut gooUtil) SHAWithRSA(key, data []byte) (string, error) {
	pkey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}

	h := crypto.Hash.New(crypto.SHA1)
	h.Write(data)
	hashed := h.Sum(nil)

	buf, err := rsa.SignPKCS1v15(rand.Reader, pkey.(*rsa.PrivateKey), crypto.SHA1, hashed)
	if err != nil {
		return "", err
	}
	return ut.Base64Encode(buf), nil
}

func (ut gooUtil) AESCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// block 大小 16
	blockSize := block.BlockSize()
	// 填充原文
	rawData = ut.pkcs7padding(rawData, blockSize)

	// 定义、初始向量IV
	cipherText := make([]byte, blockSize+len(rawData))
	iv := cipherText[:blockSize]

	// 填充向量IV
	// ReadFull从rand.Reader精确地读取len(b)字节数据填充进iv
	// rand.Reader是一个全局、共享的密码用强随机数生成器
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)

	// 二进制 转 十六进制
	encryptData := make([]byte, len(cipherText)*2)
	hex.Encode(encryptData, cipherText)
	return encryptData, nil
}

func (ut gooUtil) AESCBCDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 十六进制 转 二进制
	buf := make([]byte, len(encryptData)/2)
	hex.Decode(buf, encryptData)

	blockSize := block.BlockSize()
	l := len(buf)

	// 校验长度
	if l < blockSize {
		return nil, errors.New("encrypt data too short")
	}
	// 校验数据块
	if l%blockSize != 0 {
		return nil, errors.New("encrypt data is not a multiple of the block size")
	}

	iv := buf[:blockSize]
	encryptData = buf[blockSize:]
	origData := make([]byte, l-blockSize)

	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(origData, encryptData)
	origData = ut.pkcs7unpadding(origData)

	return origData, nil
}

func (gooUtil) pkcs7padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func (gooUtil) pkcs7unpadding(origData []byte) []byte {
	l := len(origData)
	unPadding := int(origData[l-1])
	if l < unPadding {
		return nil
	}
	return origData[:(l - unPadding)]
}

func (gooUtil) SessionId() string {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}
