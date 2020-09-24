package goo

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Prefix   string `yaml:"prefix"`
	AutoPing bool   `yaml:"auto_ping"`
}

type gooRedis struct {
	ctx    context.Context
	conf   RedisConfig
	client *redis.Client
}

func NewRedis(ctx context.Context, conf RedisConfig) *gooRedis {
	r := &gooRedis{
		ctx:  ctx,
		conf: conf,
	}
	r.new()
	if conf.AutoPing {
		go r.ping()
	}
	return r
}

func (r *gooRedis) new() {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.conf.Addr,
		Password: r.conf.Password,
		DB:       r.conf.DB,
	})
}

func (r *gooRedis) ping() {
	dur := 60 * time.Second
	ti := time.NewTimer(dur)
	defer ti.Stop()
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ti.C:
			if err := r.client.Ping().Err(); err != nil {
				Log.Error("[redis-ping]", err.Error())
			}
			ti.Reset(dur)
		}
	}
}

var __redis *gooRedis

func Redis() *redis.Client {
	return __redis.client
}

func RedisInit(conf RedisConfig) {
	__redis = NewRedis(ctx, conf)
}
