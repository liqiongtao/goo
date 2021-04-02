package goo

import (
	"encoding/json"
	"errors"
	"github.com/liqiongtao/goo/utils"
)

type Token struct {
	AppId    string
	OpenId   int64
	nonceStr string
}

func (t *Token) Bytes() []byte {
	buf, _ := json.Marshal(t)
	return buf
}

func (t *Token) String() string {
	return string(t.Bytes())
}

func CreateToken(appId string, openId int64) (tokenStr string, err error) {
	token := &Token{
		AppId:    appId,
		OpenId:   openId,
		nonceStr: utils.NonceStr(),
	}

	var (
		key    = utils.MD5([]byte(appId))
		iv     = key[8:24]
		encBuf []byte
	)

	encBuf, err = utils.AESCBCEncrypt(token.Bytes(), []byte(key), []byte(iv))
	if err != nil {
		Log.Error(err.Error())
		return
	}

	tokenStr = utils.Base64Encode(encBuf)
	return
}

func ParseToken(tokenStr, appId string) (token *Token, err error) {
	var (
		tokenBuf = utils.Base64Decode(tokenStr)
		key      = utils.MD5([]byte(appId))
		iv       = key[8:24]
		decBuf   []byte
	)

	decBuf, err = utils.AESCBCDecrypt(tokenBuf, []byte(key), []byte(iv))
	if err != nil {
		Log.Error(err.Error())
		return
	}
	token = new(Token)
	if err = json.Unmarshal(decBuf, token); err != nil {
		Log.Error(err.Error())
		return
	}
	if token.AppId != appId {
		err = errors.New("appid invalid")
		return
	}
	return
}
