// 机房信息
package soc

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type IDC struct {
	Id       int         `gorm:"primary_key" json:"id"`
	Name     string      `json:"name"`
	Location string      `json:"location"`
	SlaverId int         `json:"slaver_id"`
	Code     string      `json:"code"`
	Created  time.Time   `json:"created"`
	Updated  time.Time   `json:"updated"`
	Extra    common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (IDC) New() interface{} {
	return &IDC{}
}

func (IDC) TableName() string {
	return "byte_security_asset_idc"
}

func (idc IDC) GetAssetKey() string {
	return "name"
}

func (idc IDC) GetAssetValue() string {
	return idc.Name
}
