// 资产主机
// refer api: http://console.byted.org/tag/static/doc/index.html#api-_

package soc

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type Host struct {
	Id      int       `gorm:"primary_key" json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	// Ip 可以根据 lh 获取tag标签
	Ip           string    `gorm:"type:varchar(15);unique_index" json:"ip"` // >>> lh ip
	HostName     string    `gorm:"column:hostname" json:"host_name"`
	Os           string    `json:"os"`
	Kernel       string    `json:"kernel"`
	HoneyPot     int       `json:"honeypot"`
	HoneyPotPort int       `json:"honeypot_port"`
	State        int       `json:"state"`
	User         string    `sql:"type:text" json:"user"`
	Dpkg         string    `sql:"type:text" json:"dpkg"`
	Crontab      string    `sql:"type:text" json:"crontab"`
	AgentVersion string    `json:"agent_version"`
	Idc          string    `json:"idc"`
	Eth          string    `sql:"type:text" json:"eth"`
	Ports        string    `sql:"type:text" json:"ports"`
	Outer        []NetWork //`gorm:"many2many:byte_security_asset_host_network;association_jointable_foreignkey:network_id;jointable_foreignkey:host_id"`
	Owner        string `gorm:"index"`
	GroupId      string `gorm:"index"`
	Tags  string      `json:"tags"`
	Extra common.JSON `sql:"type:json" json:"extra,omitempty"` // tags
}

func (Host) New() interface{} {
	return &Host{}
}

func (Host) TableName() string {
	return "byte_security_asset_host"
}

func (host Host) GetAssetKey() string {
	return "ip"
}

func (host Host) GetAssetValue() string {
	return host.Ip
}
