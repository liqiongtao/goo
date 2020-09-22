package goo

type icache interface {
	init()
	ping()
}

func NewCache(cache icache) icache {
	cache.init()
	go cache.ping()
	return cache
}

var __cache icache

func InitCache(cache icache) {
	__cache = NewCache(cache)
}

func Cache() icache {
	return __cache
}
