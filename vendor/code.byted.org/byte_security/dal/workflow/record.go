package workflow

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type Record struct {
	gorm.Model
	Ticket   Ticket `gorm:"ForeignKey:TicketId"`
	TicketId int `gorm:"index"`

	User   auth.User `gorm:"ForeignKey:UserId"`
	UserId int `gorm:"index"`

	Node   Node `gorm:"ForeignKey:NodeId"`
	NodeId int `gorm:"index"`

	Content string `sql:"type:text"`
	Type    string

	Extra common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (Record) TableName() string {
	return "byte_security_workflow_record"
}
