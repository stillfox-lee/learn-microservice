package infra

import (
	"os"
	"strconv"

	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/log"
	"github.com/hashicorp/consul/api"
)

func newAgentServiceRegistration(name, id, address, port, checkUri string, tags []string, meta map[string]string) api.AgentServiceRegistration {
	portI, _ := strconv.Atoi(port)
	return api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Address: address,
		Port:    portI,
		Tags:    tags,
		Meta:    meta,
		Check: &api.AgentServiceCheck{
			HTTP:                           "http://" + address + ":" + port + checkUri,
			Interval:                       "10s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}
}

type DiscoveryClient struct {
	Client       consul.Client
	Registrar    sd.Registrar
	Registration api.AgentServiceRegistration
}

func NewDiscoveryClient(discoveryHost string, discoveryPort int) (*DiscoveryClient, error) {
	config := api.DefaultConfig()
	config.Address = discoveryHost + ":" + strconv.Itoa(discoveryPort)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	kitClient := consul.NewClient(client)

	return &DiscoveryClient{
		Client: kitClient,
	}, nil
}

func (c *DiscoveryClient) Register(name, id, address, port, checkUri string, tags []string, meta map[string]string) {
	registration := newAgentServiceRegistration(name, id, address, port, checkUri, tags, meta)
	registrar := consul.NewRegistrar(c.Client, &registration, log.NewLogfmtLogger(os.Stdout))
	c.Registrar = registrar
	c.Registration = registration

	c.Registrar.Register()
}

func (c *DiscoveryClient) Deregister() {
	c.Registrar.Deregister()
}
