package policy

import "code.byted.org/gopkg/gorm"

type GroupKey struct {
	gorm.Model
	Group   Group `gorm:"ForeignKey:GroupID"`
	GroupID uint

	Name  string //分组名称，比如psm，domain
	Value string //值，比如domain对应127.0.0.1
}

func (GroupKey) TableName() string {
	return "byte_security_policy_group_key"
}
