package policy

import "code.byted.org/gopkg/gorm"

type VariateOffline struct {
	gorm.Model
	Product   Product `gorm:"ForeignKey:ProductID"`
	ProductID int
	Name      string
	Desc      string
	Type      string
	SortType  int // 1 是字面值 3 表达式 2 列表
	//Deleted bool//用gorm.Model 中的DeletedAt
}

func (VariateOffline) TableName() string {
	return "byte_security_policy_variate_offline"
}

type VariateObjOffline struct {
	gorm.Model
	Variate   VariateOffline `gorm:"ForeignKey:VariateId"` // 一对多
	VariateId int
	Value     string `sql:"type:text"`
	//Deleted bool//用gorm.Model 中的DeletedAt
}

func (VariateObjOffline) TableName() string {
	return "byte_security_policy_variateobj_offline"
}
