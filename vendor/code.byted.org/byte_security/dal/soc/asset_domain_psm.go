// @Contact:    huaxinrui
// @Time:       2019/8/15 下午4:10

package soc

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type DomainPSM struct {
	gorm.Model
	Domain   Domain `gorm:"ForeignKey:DomainId"`
	DomainId int `gorm:"index"`
	Psm      PSM `gorm:"ForeignKey:PsmId"`
	PsmId    int `gorm:"index"`
	Paths    common.JSON `sql:"type:json" json:"paths"` // 域名下面不同的路径 []
}

func (DomainPSM) TableName() string {
	return "byte_security_asset_domain_psm"
}
