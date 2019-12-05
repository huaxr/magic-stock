// @Contact:    huaxinrui
// @Time:       2019/7/18 下午1:56

package workflow

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

// all tickets type
type TypeList struct {
	ID        uint        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *time.Time  `sql:"index" json:"deleted_at"`
	Name      string      `json:"name"`
	FormName  string      `json:"form_name"` //  对应前端组件名称
	FlowId    int         `json:"flow_id"`   // 对应的workflow_flow id
	Type      string      `json:"type"`      // 类别：normal会展示在前端，backend是后台自动生成工单
	Desc      string      `json:"desc" sql:"type:text"`
	Doc       string      `json:"doc"`   // 相关文档，url或文档内容
	State     int         `json:"state"` // 是否可用
	Extra     common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (TypeList) TableName() string {
	return "byte_security_workflow_typelist"
}
