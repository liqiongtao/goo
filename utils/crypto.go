package utils

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

func MD5(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func SHA1(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(buf, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func HMacMd5(buf, key []byte) string {
	h := hmac.New(md5.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func HMacSha1(buf, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func HMacSha256(buf, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func Base64Encode(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

func Base64Decode(str string) []byte {
	var count = (4 - len(str)%4) % 4
	str += strings.Repeat("=", count)
	buf, _ := base64.StdEncoding.DecodeString(str)
	return buf
}

func SHAWithRSA(key, data []byte) (string, error) {
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
	return Base64Encode(buf), nil
}

func AESCBCEncrypt(rawData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// block 大小 16
	blockSize := block.BlockSize()
	// 填充原文
	rawData = pkcs7padding(rawData, blockSize)

	// 定义、初始向量IV
	cipherText := make([]byte, blockSize+len(rawData))
	if iv == nil {
		iv = cipherText[:blockSize]
	}

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

func AESCBCDecrypt(encryptData, key, iv []byte) ([]byte, error) {
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

	if iv == nil {
		iv = buf[:blockSize]
	}
	encryptData = buf[blockSize:]
	origData := make([]byte, l-blockSize)

	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(origData, encryptData)
	origData = pkcs7unpadding(origData)

	return origData, nil
}

func pkcs7padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs7unpadding(origData []byte) []byte {
	l := len(origData)
	unPadding := int(origData[l-1])
	if l < unPadding {
		return nil
	}
	return origData[:(l - unPadding)]
}

func SessionId() string {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}
