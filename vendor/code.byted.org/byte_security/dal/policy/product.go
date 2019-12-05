package policy

import (
	"code.byted.org/gopkg/gorm"
)

type Product struct {
	gorm.Model
	Name  string // 产品名称
	Desc  string `sql:"type:text"`
	Uuid  string
	Key   string  // 分组关键字 psm|domain ， 默认是字符串， 会在Filed中保存
	State int //启用0，灰度1，禁用2
}

func (Product) TableName() string {
	return "byte_security_policy_product"
}

//// 创建 product
//func (p *Product) CreateProduct() {
//	p.Uuid = uuid.Must(uuid.NewV4()).String()
//	dal.Backend.DB.Save(&p)
//}
