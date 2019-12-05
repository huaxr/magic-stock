package auth

import (
	"code.byted.org/gopkg/gorm"
)

type User struct {
	gorm.Model
	UserName  string `gorm:"not null;unique"`
	RealName  string
	UserNum   string
	GroupId   string
	AvatarUrl string
	Leader    string
	Email     string
}

func (User) TableName() string {
	return "byte_security_auth_user"
}
