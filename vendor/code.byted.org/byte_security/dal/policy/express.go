package policy

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

// 线下
type RuleExpress struct {
	gorm.Model
	Product   Product `gorm:"ForeignKey:ProductID"`
	ProductID int
	Name    string
	Desc    string `sql:"type:text"`
	Express string
	UseVars common.JSON `sql:"type:json" json:"use_vars,omitempty"`

	RiskLevel int // 风险级别
	RiskScore int // 风险分数
}

func (RuleExpress) TableName() string {
	return "byte_security_policy_express"
}
