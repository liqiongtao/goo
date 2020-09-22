package goo

import (
	"context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type gooConfig struct {
	ctx      context.Context
	yamlFile string
	conf     interface{}
}

func (cf *gooConfig) AutoReLoad(dur time.Duration) {
	go func() {
		ti := time.NewTimer(dur)
		defer ti.Stop()
		for {
			select {
			case <-cf.ctx.Done():
				return
			case <-ti.C:
				if err := cf.load(); err != nil {
					log.Println("[conf-load-err]", err.Error())
				}
				ti.Reset(dur)
			}
		}
	}()
}

func (cf *gooConfig) load() error {
	bts, err := ioutil.ReadFile(cf.yamlFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(bts, cf.conf); err != nil {
		return err
	}
	return nil
}

func LoadConfig(yamlFile string, conf interface{}) *gooConfig {
	cf := &gooConfig{
		ctx:      ctx,
		yamlFile: yamlFile,
		conf:     conf,
	}
	if err := cf.load(); err != nil {
		panic(err.Error())
	}
	return cf
}
