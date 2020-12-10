package goo

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
)

type Consul struct {
	config *api.Config
	client *api.Client
}

func NewConsul(address, username, password string) *Consul {
	return &Consul{
		config: &api.Config{
			Address:  address,
			HttpAuth: &api.HttpBasicAuth{Username: username, Password: password},
		},
	}
}

func (c *Consul) Client() (*api.Client, error) {
	if c.client == nil {
		client, err := api.NewClient(c.config)
		if err != nil {
			return nil, err
		}
		c.client = client
	}
	return c.client, nil
}

func (c *Consul) Get(key string) ([]byte, error) {
	client, err := c.Client()
	if err != nil {
		return nil, err
	}
	kvp, _, err := client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	return kvp.Value, nil
}

func (c *Consul) ServiceRegister(key string) error {
	client, err := c.Client()
	if err != nil {
		return err
	}

	kvp, _, err := client.KV().Get(key, nil)
	if err != nil {
		return err
	}

	service := new(api.AgentServiceRegistration)
	if err := json.Unmarshal(kvp.Value, service); err != nil {
		return err
	}

	return client.Agent().ServiceRegister(service)
}

func (c *Consul) ServiceDeregister(serviceID string) error {
	client, err := c.Client()
	if err != nil {
		return err
	}

	return client.Agent().ServiceDeregister(serviceID)
}
