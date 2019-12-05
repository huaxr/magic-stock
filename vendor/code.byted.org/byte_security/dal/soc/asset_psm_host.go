package soc

import (
	"code.byted.org/gopkg/gorm"
)

type PSMHost struct {
	gorm.Model
	PsmId     int    `json:"psm_id"`
	Psm       string `json:"psm"`
	HostId    int    `json:"host_id"`
	IP        string `json:"ip" gorm:"index"`
	ServiceId int    `json:"service_id"`
}

func (PSMHost) TableName() string {
	return "byte_security_asset_psm_host"
}
