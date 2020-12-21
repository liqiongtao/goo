package goo

import (
	"context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type gooConfig struct {
	ctx      context.Context
	yamlFile string
	conf     interface{}
}

func (cf *gooConfig) AutoReLoad(dur time.Duration) {
	AsyncFunc(func() {
		ti := time.NewTimer(dur)
		defer ti.Stop()
		for {
			select {
			case <-cf.ctx.Done():
				return
			case <-ti.C:
				if err := cf.load(); err != nil {
					Log.Trace().Error(err.Error())
				}
				ti.Reset(dur)
			}
		}
	})
}

func (cf *gooConfig) load() error {
	bts, err := ioutil.ReadFile(cf.yamlFile)
	if err != nil {
		Log.Trace().Error(err.Error())
		return err
	}
	if err := yaml.Unmarshal(bts, cf.conf); err != nil {
		Log.Trace().Error(err.Error())
		return err
	}
	return nil
}

func LoadConfig(yamlFile string, conf interface{}) error {
	cf := &gooConfig{
		ctx:      Context,
		yamlFile: yamlFile,
		conf:     conf,
	}
	return cf.load()
}
