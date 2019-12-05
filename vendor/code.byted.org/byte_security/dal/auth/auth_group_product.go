// @Contact:    huaxinrui
// @Time:       2019/10/21 下午3:47

package auth

//Groups         string      `json:"groups"`          // 哪些组用户可以看此产品相关漏洞

type GroupProduct struct {
	Id int `gorm:"primary_key"`
	GroupId string `gorm:"index"`
	ProductId int `gorm:"index"`
}

func (GroupProduct) TableName() string {
	return "byte_security_auth_group_product"
}