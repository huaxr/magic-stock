// @Time:       2019/11/21 下午2:00

package auth

// 产品订阅者可查看
type UserBusiness struct {
	Id         int `gorm:"primary_key" json:"id"`
	BusinessId int `gorm:"index"`
	UserId     int `gorm:"index"`
}

func (UserBusiness) TableName() string {
	return "byte_security_auth_user_business"
}
