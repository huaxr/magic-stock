// @Time:       2019/11/21 下午1:56

package soc

import "time"

type Business struct {
	ID             uint       `gorm:"primary_key" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `sql:"index" json:"deleted_at"`
	Name           string     `json:"name"`            //  业务线名称
	Code           string     `json:"code"`            //
	Owner          string     `json:"owner"`           // 业务线负责人
	InterfaceOwner string     `json:"interface_owner"` // 安全接口负责人
	Products       []Product  `gorm:"FOREIGNKEY:BusinessId;ASSOCIATION_FOREIGNKEY:ID" json:"products"`
}

func (Business) TableName() string {
	return "byte_security_asset_business"
}
