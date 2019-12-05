package policy

import (
	"code.byted.org/gopkg/gorm"
)

type StrategyVersionOffline struct {
	gorm.Model
	RuleExpress       []RuleExpress   //相当于规则的一部分， 但是缺少决策结果， 优先级字段（多对多)） m2m 中保存映射
	Strategy          StrategyOffline `gorm:"ForeignKey:StrategyID" json:"strategy"`
	StrategyOfflineID int            `json:"strategy_id"`
	VersionName       string          `json:"version_name"`
	Filter string  // 过滤条件: 先判断要不要进这个策略  e.g.  express - > bool
}

func (StrategyVersionOffline) TableName() string {
	return "byte_security_policy_strategy_version_offline"
}

//func (s *StrategyVersionOffline) GetVersionRulesCount() int {
//	common.Backend.DB.Model(&s).Related(&s.RuleExpress, "StrategyVersionId")
//	return len(s.RuleExpress)
//}

//func (s *StrategyVersion) CopyRules(r StrategyVersion) {
//	common.Backend.DB.Model(&r).Related(&r.RuleExpress, "StrategyVersionId")
//	for _, i := range r.RuleExpress {
//		r := Rule{StrategyVersionId:  i.StrategyVersionId, Name:i.Name, Express:i.Express, Status:i.Status}
//		common.Backend.DB.Save(&r)
//	}
//}
