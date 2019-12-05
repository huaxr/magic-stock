// 资产产品

package soc

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type Product struct {
	ID             uint        `gorm:"primary_key" json:"id"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DeletedAt      *time.Time  `sql:"index" json:"deleted_at"`
	Name           string      `json:"name"`
	Code           string      `json:"code"`
	Department     string      `json:"department"`
	Owner          string      `json:"owner"`
	InterfaceOwner string      `json:"interface_owner"` // 安全接口负责人
	Group          string      `json:"group"`
	BusinessId     int         `gorm:"index" json:"business_id"` // 所属业务线
	Psm            string      `json:"psm"`
	Extra          common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (Product) New() interface{} {
	return &Product{}
}

func (Product) TableName() string {
	return "byte_security_asset_product"
}

func (product *Product) GetOwnerNames() (name []string) {
	name = append(name, product.Owner)
	return
}

func (product Product) GetAssetKey() string {
	return "name"
}

func (product Product) GetAssetValue() string {
	return product.Name
}
