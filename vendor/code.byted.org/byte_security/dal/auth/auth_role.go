// @Contact:    huaxinrui
// @Time:       2019/8/21 上午10:33

package auth

type Role struct {
	ID       uint  `gorm:"primary_key"`
	RoleName string `gorm:"not null;unique"`
}

func (Role) TableName() string {
	return "byte_security_auth_role"
}