package workflow

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type WorkFlow struct {
	gorm.Model
	Code          string
	Type          string         // 对应的all中的name
	State         int            // 是否可用
	Nodes         []NodeTemplate `gorm:"ForeignKey:FlowId"` // 一个工作流包含多个节点
	DefaultUserId int            // 默认处理人 the admin to handle the work flow, which means every event will trigger this guy ..
	Extra         common.JSON    `sql:"type:json" json:"extra,omitempty"`
}

func (WorkFlow) TableName() string {
	return "byte_security_workflow_flow"
}
