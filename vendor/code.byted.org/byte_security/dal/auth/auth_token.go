// @Contact:    huaxinrui
// @Time:       2019/9/4 下午5:46

package auth


import (
	"code.byted.org/gopkg/gorm"
)


type Token struct {
	gorm.Model
	Token string `gorm:"not null;unique"`
	Name string
	Owner string
	Path string  // 指定路径
}


func (Token) TableName() string {
	return "byte_security_auth_token"
}