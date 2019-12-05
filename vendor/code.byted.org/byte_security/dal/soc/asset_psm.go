// @Contact:    huaxinrui
// @Time:       2019/8/15 下午4:08

package soc

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type PSM struct {
	gorm.Model
	Psm         string      `json:"psm" gorm:"index"`
	Owner       string      `json:"owner"`
	GroupId     string      `json:"group_id"`
	ScmRepos    common.JSON `sql:"type:json" json:"scm_repos"`
	Subscribers common.JSON `sql:"type:json" json:"subscribers"`
	NodeID      int         `json:"node_id"`
	Path        string      `json:"path"`
	Extra       common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (PSM) New() interface{} {
	return &PSM{}
}

func (PSM) TableName() string {
	return "byte_security_asset_psm"
}

func (psm PSM) GetAssetKey() string {
	return "psm"
}

func (psm PSM) GetAssetValue() string {
	return psm.Psm
}
