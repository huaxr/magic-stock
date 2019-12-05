package policy

import (
	"code.byted.org/gopkg/gorm"
)


// 线下策略没有生效版本
type StrategyOffline struct {
	gorm.Model
	Product                Product                  `gorm:"ForeignKey:ProductID" json:"product"`
	ProductID              int                     `json:"product_id"`
	StrategyVersionOffline []StrategyVersionOffline `json:"strategy_version_list"` // 一对多,json后面的那个名称需要改吗？
	Name                   string                   `gorm:"type:varchar(100)" json:"name"`
	Desc                   string                   `sql:"type:text" json:"desc"`
	Priority int `json:"priority"`
}

func (StrategyOffline) TableName() string {
	return "byte_security_policy_strategy_offline"
}
