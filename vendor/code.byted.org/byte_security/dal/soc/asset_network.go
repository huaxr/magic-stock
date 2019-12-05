// 外网资产
package soc

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type NetWork struct {
	Id      int         `gorm:"primary_key" json:"id"`
	Ip      string      `json:"ip"`
	Idc     IDC         `gorm:"ForeignKey:IdcID" json:"idc"`
	IdcID   int         `json:"idc_id"`
	Status  int         `json:"status"`
	Created time.Time   `json:"created"`
	Updated time.Time   `json:"updated"`
	Extra   common.JSON `sql:"type:json" json:"extra,omitempty"` // port
}

func (NetWork) New() interface{} {
	return &NetWork{}
}

type M2MOuter struct {
	Id         int `gorm:"primary_key"`
	NetworkId  int
	HostId     int
	UpdateTime time.Time `json:"updated"`
}

func (NetWork) TableName() string {
	return "byte_security_asset_network"
}

func (M2MOuter) TableName() string {
	return "byte_security_asset_host_network"
}

func (network NetWork) GetAssetKey() string {
	return "ip"
}

func (network NetWork) GetAssetValue() string {
	return network.Ip
}
