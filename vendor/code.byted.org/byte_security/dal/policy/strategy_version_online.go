package policy

import (
	"code.byted.org/gopkg/gorm"
)

type StrategyVersionOnline struct {
	gorm.Model
	Rules []Rule  `gorm:"ForeignKey:StrategyVersionOnlineId"` //一对多的关系
	//RuleExpress []RuleExpress
	Strategy         StrategyOnline `gorm:"ForeignKey:StrategyID" json:"strategy"`
	StrategyOnlineID int           `json:"strategy_id"`
	VersionName      string         `json:"version_name"`
	Filter string `json:"filter"`
}

func (StrategyVersionOnline) TableName() string {
	return "byte_security_policy_strategy_version_online"
}

//func (s *StrategyVersion) CopyRules(r StrategyVersion) {
//	common.Backend.DB.Model(&r).Related(&r.RuleExpress, "StrategyVersionId")
//	for _, i := range r.RuleExpress {
//		r := Rule{StrategyVersionId:  i.StrategyVersionId, Name:i.Name, Express:i.Express, Status:i.Status}
//		common.Backend.DB.Save(&r)
//	}
//}
