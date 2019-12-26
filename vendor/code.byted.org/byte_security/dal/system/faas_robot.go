package system

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type FAASRobot struct {
	ID          uint        `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `sql:"index" json:"deleted_at"`
	Name        string      `json:"name"`         // faas机器人名称
	Desc        string      `json:"desc"`         // faas机器人描述
	Owner       string      `json:"owner"`        // 机器人Owner
	URL         string      `json:"url"`          // faas调用URL
	Token       string      `json:"token"`        //
	CallBackURL string      `json:"callback_url"` //  处理人员：当type为human时，表示指定的处理人员；当type为robot时，表示机器人的ID;当type为auto时，duty/asset_owner/leader,对应值班表、资产负责人、操作用户leader
	Extra       common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (FAASRobot) TableName() string {
	return "faas_robot"
}
