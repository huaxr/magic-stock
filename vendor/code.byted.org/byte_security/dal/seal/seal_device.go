package seal

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type Device struct {
	ID         int         `gorm:"primary_key"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  *time.Time  `sql:"index" json:"deleted_at"`
	Did        string      `gorm:"not null;unique"`
	Owner string `json:"owner"`
	DeviceName string      `json:"device_name"`
	MacAddr    string      `sql:"type:text"`
	Os         string      `json:"os"`
	OsVer      string      `json:"os_ver"`
	AppVer     string      `json:"app_ver"`
	Brand      string      `json:"brand"`
	Model      string      `gorm:"column:model"`
	Uid        string      `json:"uid"`
	Extra      common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (Device) New() interface{} {
	return &Device{}
}

func (Device) TableName() string {
	return "byte_security_seal_device"
}

func (device Device) GetAssetKey() string {
	return "did"
}

func (device Device) GetAssetValue() string {
	return device.Did
}
