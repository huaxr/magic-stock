// @Contact:    huaxinrui
// @Time:       2019/10/28 上午11:18

package auth

// 产品订阅者可查看
type UserProduct struct {
	Id        int `gorm:"primary_key" json:"id"`
	ProductId int `gorm:"index"`
	UserId    int `gorm:"index"`
}

func (UserProduct) TableName() string {
	return "byte_security_auth_user_product"
}
