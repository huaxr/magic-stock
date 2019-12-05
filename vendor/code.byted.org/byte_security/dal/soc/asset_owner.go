// @Contact:    huaxinrui
// @Time:       2019/8/2 下午3:44

package soc

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/gopkg/gorm"
)


// 资产的组， 用户信息
type AssetOwner struct {
	gorm.Model
	Type    string                    // 资产类型
	AssetId int        `gorm:"index"` // 资产ID
	User    auth.User  `gorm:"ForeignKey:UserId"`
	UserId  int        `gorm:"index"` // 所属人
	Group   auth.Group `gorm:"ForeignKey:GroupId"`
	GroupId int        `gorm:"index"` // 所属组
}

func (AssetOwner) TableName() string {
	return "byte_security_asset_type_owner"
}
