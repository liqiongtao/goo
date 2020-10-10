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

func AESECBEncrypt(data, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := cipher.BlockSize()
	paddingSize := blockSize - len(data)%blockSize
	if paddingSize != 0 {
		data = append(data, bytes.Repeat([]byte{byte(0)}, paddingSize)...)
	}
	encrypted := make([]byte, len(data))
	for bs, be := 0, blockSize; bs < len(data); bs, be = bs+blockSize, be+blockSize {
		cipher.Encrypt(encrypted[bs:be], data[bs:be])
	}
	return encrypted, nil
}

func AESECBDecrypt(buf, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := cipher.BlockSize()
	decrypted := make([]byte, len(buf))
	for bs, be := 0, blockSize; bs < len(buf); bs, be = bs+blockSize, be+blockSize {
		cipher.Decrypt(decrypted[bs:be], buf[bs:be])
	}
	paddingSize := int(decrypted[len(decrypted)-1])
	return decrypted[0 : len(decrypted)-paddingSize], nil
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

	// 定义密码数据
	var cipherData []byte

	// 如果iv为空，生成随机iv，并附加到加密数据前面，否则单独生成加密数据
	if iv == nil {
		// 初始化加密数据
		cipherData = make([]byte, blockSize+len(rawData))
		// 定义向量
		iv = cipherData[:blockSize]
		// 填充向量IV， ReadFull从rand.Reader精确地读取len(b)字节数据填充进iv，rand.Reader是一个全局、共享的密码用强随机数生成器
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
		// 加密
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(cipherData[blockSize:], rawData)
	} else {
		// 初始化加密数据
		cipherData = make([]byte, len(rawData))
		// 定义向量
		iv = iv[:blockSize]
		// 加密
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(cipherData, rawData)
	}

	return cipherData, nil
}

func AESCBCDecrypt(cipherData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// block 大小 16
	blockSize := block.BlockSize()

	// 加密串长度
	l := len(cipherData)

	// 校验长度
	if l < blockSize {
		return nil, errors.New("encrypt data too short")
	}

	// 定义原始数据
	var origData []byte

	// 如果iv为空，需要获取前16位作为随机iv
	if iv == nil {
		// 定义向量
		iv = cipherData[:blockSize]
		// 定义真实加密串
		cipherData = cipherData[blockSize:]
		// 初始化原始数据
		origData = make([]byte, l-blockSize)
	} else {
		// 定义向量
		iv = iv[:blockSize]
		// 初始化原始数据
		origData = make([]byte, l)
	}

	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(origData, cipherData)
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
