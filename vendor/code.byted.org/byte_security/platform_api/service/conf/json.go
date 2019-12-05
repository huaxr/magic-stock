package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"code.byted.org/byte_security/platform_api/utils"
)

func init() {
	Init()
}

var Config config

type config struct {
	Store           string `json:"store"`
	Mode            string `json:"mode"`
	Addr            string `json:"addr"`
	Psm             string `json:"psm"`
	CMK             string `json:"cmk"`
	SessionSecret   string `json:"session_secret"`
	RedisCluster    string `json:"redis_cluster"`
	Host            string `json:"host"`
	HostMethod      string `json:"host_method"`
	SSO             string `json:"sso"`
	ES              string `json:"es"`
	Sentry          string `json:"sentry"`
	DomainURLPrefix string `json:"domain_url_prefix"`
}

func Init() {
	var ConfigFile string
	switch utils.TellEnv() {
	case "loc":
		fmt.Println("hi loc env")
		//ConfigFile = "conf/loc.json"
		ConfigFile = "/Users/huaxinrui/go/src/code.byted.org/byte_security/platform_api/conf/loc.json"
	case "boe":
		ConfigFile = "conf/boe.json"
	case "tce":
		ConfigFile = "conf/tce.json"
	}

	d, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatalf("【*】read config error %v\n", err)
	}
	if err = json.Unmarshal(d, &Config); err != nil {
		log.Fatalf("【*】unmarshal conf json error %v\n", err)
	}
}
