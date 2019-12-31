package policy

import (
	"code.byted.org/gopkg/gorm"
)

type Product struct {
	gorm.Model
	Name  string // 产品名称
	Desc  string `sql:"type:text"`
	Uuid  string
	Key   string // 分组关键字 psm|domain ， 默认是字符串， 会在Filed中保存
	State int    //启用0，灰度1，禁用2

	Code       string // 代号
	Doc        string // 文档
	Tag        string // （数据安全/代码安全）
	Type       string // 能力 0决策，1基线，2第三方扫描，3第三方监控，4第三方加固
	FieldsId   int    // 关联字段集（跟分析有关的能力）
	FieldsName string // 字段集名称
	Asset      string // 资产（跟分析有关的能力） psm, repo, domain, host
	AssetField string // 资产对应的key 如 host->ip
}

func (Product) TableName() string {
	return "byte_security_policy_product"
}
