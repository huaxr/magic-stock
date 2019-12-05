// @Contact:    huaxinrui
// @Time:       2019/11/7 下午5:49

package auth

type UserRole struct {
	ID       uint  `gorm:"primary_key"`
	RoleId int `gorm:"index"`
	UserId int `gorm:"index"`
	Limit string // 限制应用资源的模块权限 如只能是 事件管理role下的 hids 类型事件 范围在 policy_product 下
}

func (UserRole) TableName() string {
	return "byte_security_auth_user_role"
}