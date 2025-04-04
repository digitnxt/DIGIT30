package discovery

import (
	"github.com/hashicorp/consul/api"
)

// RegisterService registers a service with Consul.
// serviceID: Unique identifier for the service.
// serviceName: The name of the service.
// address: The service address.
// port: The service port.
// checkURL: The HTTP URL for health check.
func RegisterService(serviceID, serviceName, address string, port int, checkURL string) error {
	config := api.DefaultConfig()
	config.Address = "consul:8500"
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Check: &api.AgentServiceCheck{
			HTTP:     checkURL,
			Interval: "10s",
			Timeout:  "5s",
		},
	}
	return client.Agent().ServiceRegister(registration)
}

// DeregisterService deregisters a service from Consul using its service ID.
func DeregisterService(serviceID string) error {
	config := api.DefaultConfig()
	config.Address = "consul:8500"
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	return client.Agent().ServiceDeregister(serviceID)
}