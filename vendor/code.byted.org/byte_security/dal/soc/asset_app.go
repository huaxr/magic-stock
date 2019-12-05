// 资产APP

package soc

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type App struct {
	ID             uint        `gorm:"primary_key" json:"id"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DeletedAt      *time.Time  `sql:"index" json:"deleted_at"`
	Name           string      `json:"name"`
	Platform       string      `json:"platform"`
	Language       string      `json:"language"`
	Desc           string      `gorm:"size:255"  json:"desc"`
	DownloadUrl    string      `json:"download_url"`
	ProductID      int         `json:"product_id"`
	Product        Product     `gorm:"ForeignKey:ProductID"  json:"product"`
	BusinessId     int         `gorm:"index" json:"business_id"` // 所属业务线      `gorm:"index" json:"business_code"` // 所属业务线
	Owner          string      `json:"owner"`
	InterfaceOwner string      `json:"interface_owner"`
	Extra          common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (App) New() interface{} {
	return &App{}
}

func (App) TableName() string {
	return "byte_security_asset_app"
}

func (app App) GetAssetKey() string {
	return "name"
}

func (app App) GetAssetValue() string {
	return app.Name
}
