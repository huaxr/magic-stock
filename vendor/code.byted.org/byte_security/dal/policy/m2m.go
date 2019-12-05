package policy

import "code.byted.org/gopkg/gorm"

type GroupAndStrategyOnline struct {
	gorm.Model
	GroupID          int
	StrategyOnlineId int
}

type ExpressAndStrategyVersionOffline struct {
	gorm.Model
	ExpressID                int
	StrategyVersionOfflineId int
	Status                   int // 决策结果
	Priority                 int // 优先级
}


type GroupAndLocation struct {
	gorm.Model
	GroupID    int
	LocationId int
}

func (GroupAndStrategyOnline) TableName() string {
	return "byte_security_policy_group_strategyonline"
}

func (GroupAndLocation) TableName() string {
	return "byte_security_policy_group_location"
}

func (ExpressAndStrategyVersionOffline) TableName() string {
	return "byte_security_policy_strategyversionoffline_express"
}

//// 多对多查询 GroupAndStrategy
//func GetGroupByStrategyId(strategy_id int) []Group {
//	var m2m []GroupAndStrategyOnline
//	var groups []Group
//	common.Backend.DB.Model(&GroupAndStrategyOnline{}).Where("strategy_id = ?", strategy_id).Find(&m2m)
//	for _, i := range m2m {
//		var g Group
//		common.Backend.DB.Model(&Group{}).Where("group_id = ?", i.GroupID).Find(&g)
//		if g.ID == 0 {
//			continue
//		}
//		groups = append(groups, g)
//	}
//	return groups
//}
//
//func CreateGroupAndLocation(group_id, location_id int) {
//	gl := GroupAndLocation{GroupID: group_id, LocationId: location_id}
//	common.Backend.DB.Save(&gl)
//}
//
//func CreateGroupAndStrategy(group_id, strategy_id int) {
//	gs := GroupAndStrategyOnline{GroupID: group_id, StrategyOnlineId: strategy_id}
//	common.Backend.DB.Save(&gs)
//}
