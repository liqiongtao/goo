package sms

import (
	"fmt"
	"sync"
	"time"
)

var __cache = cache{}

type cache struct {
	Data sync.Map
}

func (ca cache) set(appid, mobile, action, code string, expireIn int64) {
	key := fmt.Sprintf("%s_%s_%s", appid, mobile, action)
	ca.Data.Store(key, &codeInfo{Code: code, ExpireOut: time.Now().Unix() + expireIn})
	go ca.recovery()
}

func (ca cache) get(appid, mobile, action string) *codeInfo {
	defer func() {
		go ca.recovery()
	}()

	key := fmt.Sprintf("%s_%s_%s", appid, mobile, action)
	rst, ok := ca.Data.Load(key)
	if !ok {
		return nil
	}
	return rst.(*codeInfo)
}

func (ca cache) recovery() {
	ca.Data.Range(func(k, v interface{}) bool {
		if time.Now().Unix() > v.(*codeInfo).ExpireOut {
			ca.Data.Delete(k)
		}
		return true
	})
}
