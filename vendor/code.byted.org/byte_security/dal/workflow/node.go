package workflow

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type Node struct {
	ID          uint        `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `sql:"index" json:"deleted_at"`
	Name        string      `json:"name"`                                // Node 的名称
	Desc        string      `json:"desc"`                                // 后期去掉冗余字段
	TicketID    int         `json:"ticket_id" gorm:"index"`              // 关联的ticket
	FlowID      int         `json:"flow_id" gorm:"index"`                // 关联的工作流ID
	Priority    int         `json:"priority"`                            // 优先级链
	State       string      `json:"state" gorm:"index"`                  // 状态 pending:处理中, rejected:驳回, finished:完成, suspended:挂起, urgent:已加急
	HandlerType string      `json:"handler_type"`                        // 处理方式：human指定人员，robot自动函数处理，auto自动获取
	Handler     string      `json:"handler"`                             // 当handler_type为机器人时，填写机器人ID
	NodeType    string      `json:"node_type" gorm:"DEFAULT:'approval'"` // 节点类别：approval 需要审批，notify 仅通知
	Extra       common.JSON `sql:"type:json" json:"extra,omitempty"`
	Type        string      `json:"type"` // TODO: DELETE
}

func (Node) TableName() string {
	return "byte_security_workflow_node"
}
