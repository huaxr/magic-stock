package policy

import "code.byted.org/gopkg/gorm"

type Rule struct {
	gorm.Model
	StrategyVersionOnline   StrategyVersionOnline `gorm:"ForeignKey:StrategyVersionOnlineId"`
	StrategyVersionOnlineId int
	Name                    string
	Express                 string
	Status                  int
	Priority int
}


func (Rule) TableName() string {
	return "byte_security_policy_rule"
}
