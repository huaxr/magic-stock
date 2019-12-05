package policy

import "code.byted.org/gopkg/gorm"

type Location struct {
	gorm.Model
	Name string
	Value string
	Desc string `sql:"type:text"`
}


func (Location) TableName() string {
	return "byte_security_policy_location"
}