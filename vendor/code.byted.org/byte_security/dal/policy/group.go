package policy

import "code.byted.org/gopkg/gorm"

// 分组和策略是多对多的关系
// 其中产品相当于一个命名空间
type Group struct {
	gorm.Model
	Product   Product `gorm:"ForeignKey:ProductID"`
	ProductID int

	Desc  string `sql:"type:text"`
	State int    //启用0，灰度1，禁用2
	Unique string `gorm:"unique"`

	GroupKeys []GroupKey       // 这个就是分组关键字， 在新建分组的时候填入psm/domain 保存在这个表中
	// 线下策略和分组没有关系
	Strategys []StrategyOnline //`gorm:"many2many:byte_security_waf_group_strategy;association_jointable_foreignkey:strategy_id;jointable_foreignkey:group_id"`
	Locations []Location       //`gorm:"many2many:byte_security_waf_group_strategy;association_jointable_foreignkey:location_id;jointable_foreignkey:group_id"`
}

func (Group) TableName() string {
	return "byte_security_policy_group"
}
