package discovery

import (
	"strconv"
	"strings"

	"code.byted.org/gopkg/consul"
)

// ConsulDiscover discover this service with specifical idc by consul
func ConsulDiscover(serviceName, idc string) ([]*Instance, error) {
	idc = strings.TrimSpace(idc)
	items, err := consul.Lookup(serviceName, consul.WithIDC(consul.IDC(idc)))
	if err != nil {
		return nil, err
	}

	var ret []*Instance
	for _, ins := range items {
		addr := strings.Split(ins.Addr, ":")
		ret = append(ret, NewInstance(addr[0], addr[1], map[string]string{
			"cluster": ins.Cluster,
			"env":     ins.Env,
			"weight":  strconv.Itoa(ins.Weight),
		}))
	}
	return ret, nil
}

// ConsulDiscoverer .
type ConsulDiscoverer struct{}

// NewConsulDiscoverer .
func NewConsulDiscoverer() *ConsulDiscoverer {
	return &ConsulDiscoverer{}
}

// Discover .
func (c *ConsulDiscoverer) Discover(serviceName, idc string) ([]*Instance, error) {
	return ConsulDiscover(serviceName, idc)
}
