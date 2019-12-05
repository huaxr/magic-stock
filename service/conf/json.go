package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"magic/stock/utils"
)

var Config config

type config struct {
	Store         string `json:"store"`
	SessionSecret string `json:"session_secret"`
	Host          string `json:"host"`
}

func init() {
	var ConfigFile string
	switch utils.TellEnv() {
	case "loc":
		fmt.Println("hi loc env")
		//ConfigFile = "conf/loc.json"
		ConfigFile = "/Users/huaxinrui/go/src/magic/stock/conf/loc.json"
	case "dev":
		ConfigFile = "/Users/huaxinrui/go/src/magic/stock/conf/dev.json"
	case "online":
		ConfigFile = "/home/tiger/go/src/magic/stock/output/conf/online.json"
	}

	d, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatalf("【*】read config error %v\n", err)
	}
	if err = json.Unmarshal(d, &Config); err != nil {
		log.Fatalf("【*】unmarshal conf json error %v\n", err)
	}
}
