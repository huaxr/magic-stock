// @Contact:    huaxinrui
// @Time:       2019/7/30 下午5:04

package auth


import "code.byted.org/gopkg/gorm"


type Policy struct {
	gorm.Model
	Group, Path, Method string
}


func (Policy) TableName() string {
	return "byte_security_auth_policy"
}