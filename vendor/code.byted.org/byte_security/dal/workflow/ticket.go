package workflow

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

// 工单的输入是 事件， 有漏洞事件， waf事件等， 需要对应 不同的工单
// 不同的工单有不同的 工作流， 属于一对一的包含关系

type Ticket struct {
	gorm.Model
	Title           string
	Desc            string
	Type            string
	Creator         auth.User   `gorm:"ForeignKey:UserID"`
	UserId          int         `gorm:"index"` // 创建者
	State           string      `gorm:"index"` // ('pending', '处理中'),  ('finished', '完成'),  ('cancelled', '取消'), ('rejected', '拒绝')
	Detail          common.JSON `sql:"type:json" json:"detail"`
	Extra           common.JSON `sql:"type:json" json:"extra,omitempty"`
	TicketSource    string      `gorm:"index"` // event, normal, vuln
	Psm             string      `gorm:"index"` // 对应psm
	Asset           string      `gorm:"index"` // 对应资产
	AssetOwner      string      // 资产的拥有者
	AssetType       string      // domain, host
	Level           int         // 风险级别， 也就是事件的风险级别
	Group           string      // 资产所属组
	CurrentPriority int         // 当前节点
	BusinessID      int         `json:"business_id"`
	ProductID       int         `json:"product_id"`
}

func (Ticket) TableName() string {
	return "byte_security_workflow_ticket"
}
