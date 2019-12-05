package workflow

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type Node struct {
	gorm.Model
	Type     string //('review', '审批'), ('run', '执行')
	Name     string // Node 的名称
	Desc     string
	FlowId   int `gorm:"index"`
	TicketId int `gorm:"index"`
	Priority int       // 优先级链
	State    string  `gorm:"index"`  // 状态  pending, rejected, finished, abort
	Operator auth.User `gorm:"ForeignKey:UserId"`
	Extra    common.JSON `sql:"type:json" json:"extra,omitempty"`
}

type NodeTemplate struct {
	gorm.Model
	Type     string //('review', '审批'), ('run', '执行')
	Name     string // Node 的名称
	Desc     string
	FlowId   int
	Priority int    // 优先级链
	State    string // 状态  pending, rejected, finished
	//Operators []auth.User `gorm:"many2many:byte_security_workflow_nodetemp_user;association_jointable_foreignkey:user_id;jointable_foreignkey:node_id"`
	Operator auth.User `gorm:"ForeignKey:UserId"`
	UserNames common.JSON `sql:"type:json" json:"extra,omitempty"` // 可以有多个处理人
	Extra    common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (Node) TableName() string {
	return "byte_security_workflow_node"
}

func (NodeTemplate) TableName() string {
	return "byte_security_workflow_node_template"
}
