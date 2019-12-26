package workflow

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type NodeTemplate struct {
	ID          uint        `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `sql:"index" json:"deleted_at"`
	Name        string      `json:"name"`                              // Node 的名称
	FlowID      int         `json:"flow_id" gorm:"index"`              // 关联的工作流ID
	Desc        string      `json:"desc"`                              // 节点的描述
	Priority    int         `json:"priority"`                          // 优先级链
	HandlerType string      `json:"handler_type"`                      // 处理方式：human指定人员，robot自动函数处理，auto自动获取
	Handler     string      `json:"handler"`                           //  处理人员：当type为human时，表示指定的处理人员；当type为robot时，表示机器人的ID;当type为auto时，duty/asset_owner/leader,对应值班表、资产负责人、操作用户leader
	Urgent      int         `json:"urgent"`                            // 是否可加急，0不可，1可以
	NodeType    string      `json:"node_type" gorm:"DEFAULT:approval"` // 节点类别：approval 需要审批，notify 仅通知
	Extra       common.JSON `sql:"type:json" json:"extra,omitempty"`
	Type        string      `json:"type"`       // TODO: DELETE
	State       string      `json:"state"`      // TODO: DELETE
	UserNames   common.JSON `json:"user_names"` // TODO: DELETE
}

func (NodeTemplate) TableName() string {
	return "byte_security_workflow_node_template"
}

type WorkFlow struct {
	ID            uint           `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     `sql:"index" json:"deleted_at"`
	Name          string         `json:"name"`  // 工作流名称
	Desc          string         `json:"desc"`  // 描述
	State         int            `json:"state"` // 是否可用
	Nodes         []NodeTemplate `gorm:"FOREIGNKEY:FlowID;ASSOCIATION_FOREIGNKEY:ID" json:"nodes"`
	Extra         common.JSON    `sql:"type:json" json:"extra,omitempty"`
	Code          string         `json:"code"`            // TODO: DELETE
	Type          string         `json:"type"`            // TODO: DELETE
	DefaultUserID string         `json:"default_user_id"` // TODO: DELETE
}

func (WorkFlow) TableName() string {
	return "byte_security_workflow_flow"
}
