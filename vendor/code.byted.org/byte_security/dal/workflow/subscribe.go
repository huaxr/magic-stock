// @Contact:    huaxinrui
// @Time:       2019/7/9 下午3:57

package workflow

import (
	"code.byted.org/gopkg/gorm"
)

// 订阅者
type Subscribe struct {
	gorm.Model
	TicketID int `gorm:"index"`
	NodeID int `gorm:"index"`
	UserId int `gorm:"index"`
}


func (Subscribe) TableName() string {
	return "byte_security_workflow_subscribe"
}
