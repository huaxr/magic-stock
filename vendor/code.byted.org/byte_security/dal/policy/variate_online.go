package policy

import "code.byted.org/gopkg/gorm"

type VariateOnline struct {
	gorm.Model
	Product Product `gorm:"ForeignKey:ProductID"`
	ProductID int
	StrategyVersionOnline   StrategyVersionOnline `gorm:"ForeignKey:StrategyVersionOnlineId"`
	StrategyVersionOnlineId uint  // 属于策略版本
	Name     string `gorm:"type:varchar(255)"`
	Desc string
	Type     string
	SortType int // 1 2列表 3
	IsStore int // 变量结果是否放在结论中
}


func (VariateOnline) TableName() string {
	return "byte_security_policy_variate_online"
}


type VariateObjOnline struct {
	gorm.Model
	Variate   VariateOnline `gorm:"ForeignKey:VariateId"`
	VariateId int
	Value     string `sql:"type:text"`
	//Deleted bool//用gorm.Model 中的DeletedAt
}

func (VariateObjOnline) TableName() string {
	return "byte_security_policy_variateobj_online"
}
