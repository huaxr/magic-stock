package auth

import (
	"code.byted.org/gopkg/gorm"
)

type User struct {
	gorm.Model
	UserName  string `gorm:"not null;unique"`
	RealName  string
	Email     string
	UserNum   string
	GroupId   string
	AvatarUrl string
	Leader    string
}

func (User) TableName() string {
	return "byte_security_auth_user"
}
