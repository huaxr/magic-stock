package policy

import "code.byted.org/gopkg/gorm"

type StrategyOnline struct {
	gorm.Model
	Product               Product                 `gorm:"ForeignKey:ProductID" json:"product"`
	ProductID             int                    `json:"product_id"`
	StrategyVersionOnline []StrategyVersionOnline `json:"strategy_version_list"` // 一对多,json后面的那个名称需要改吗？
	Name                  string                  `json:"name"`
	Desc                  string                  `sql:"type:text" json:"desc"`
	VersionActiveId       int                    `json:"version_active_id"`
	Priority              int                     `json:"priority"`
	OfflineId int
}

func (StrategyOnline) TableName() string {
	return "byte_security_policy_strategy_online"
}
