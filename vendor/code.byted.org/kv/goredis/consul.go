package goredis

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"code.byted.org/golf/consul"
	"code.byted.org/gopkg/logs"
)

const (
	HY_IDC = "hy"
	LF_IDC = "lf"
)

var (
	localIDC   string
	localOncer sync.Once
)

// LocalIDC return idc's name of current service
// first read env val RUNTIME_IDC_NAME
func LocalIDC() string {
	localOncer.Do(func() {
		if dc := os.Getenv("RUNTIME_IDC_NAME"); dc != "" {
			localIDC = strings.TrimSpace(dc)
		} else {
			sd, err := consul.NewServiceDiscovery()
			if err != nil {
				localIDC = ""
				return
			}
			localIDC = strings.TrimSpace(sd.Dc)
		}
	})
	return localIDC
}

type Instance struct {
	host string
	port string
	tags map[string]string
}

func NewInstance(host, port string, tags map[string]string) *Instance {
	for key, val := range tags {
		tags[key] = strings.TrimSpace(val)
	}
	return &Instance{
		host: strings.TrimSpace(host),
		port: strings.TrimSpace(port),
		tags: tags,
	}
}

func (it *Instance) Host() string {
	return it.host
}

func (it *Instance) Port() string {
	return it.port
}

func (it *Instance) Str() string {
	return it.host + ":" + it.port
}

type ConsulService struct {
	name string
}

func NewConsulService(name string) *ConsulService {
	return &ConsulService{name}
}

func (cs *ConsulService) Name() string {
	return cs.name
}

// Lookup return a list of instances
func (cs *ConsulService) Lookup(idc string) []*Instance {
	idc = strings.TrimSpace(idc)
	if len(idc) > 0 {
		idc = "." + idc
	}
	items, err := consul.TranslateOne(cs.name + ".service" + idc)
	if err != nil {
		logs.Errorf("Redisclient consul.TranslateOne error, cluster name:%s, error:%s", cs.name, err)
		return nil
	}

	var ret []*Instance
	for _, ins := range items {
		ret = append(ret, NewInstance(ins.Host, strconv.Itoa(ins.Port), ins.Tags))
	}
	return ret
}
