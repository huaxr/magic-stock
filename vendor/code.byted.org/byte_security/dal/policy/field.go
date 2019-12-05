package policy

import (
	"code.byted.org/gopkg/gorm"
)

//type HResults struct {
//	gorm.Model
//	Asset string // 资产 如 ip， domain
//	Uid string `gorm:"index"`
//	RuleStr string
//	Type string  // 事件类型
//	Level int    // 事件等级
//	Key int
//	//Detail common.JSON   `sql:"type:json" json:"object,omitempty"`
//	RawLog common.JSON   `sql:"type:json" json:"object,omitempty"`
//	RawMap common.JSON   `sql:"type:json" json:"object,omitempty"`
//	Users  common.JSON   `sql:"type:json" json:"object,omitempty"`
//	State string   // pending done
//	Handler string // 处理人
//}

type Field struct {
	gorm.Model
	Name string
	Type string

	MappingField string // asset
	MappingName  string // 资产 or 属性
	MappingTable string // 关联资产表 byte_security_asset_domain 如果是资产的话必填项

	Product   Product `gorm:"ForeignKey:ProductID"`
	ProductID int
}


func (Field) TableName() string {
	return "byte_security_policy_field"
}
