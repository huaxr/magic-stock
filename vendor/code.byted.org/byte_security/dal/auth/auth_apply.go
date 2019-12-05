package auth

import "code.byted.org/gopkg/gorm"


type AuthApply struct {
	gorm.Model
	Who string
	Url string `gorm:"size:255"`
	Body string `gorm:"size:255"`
	Message string `gorm:"size:255"`
}


func (AuthApply) TableName() string {
	return "byte_security_auth_apply"
}