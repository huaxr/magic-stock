// @Time:       2019/12/2 下午4:15

package dal

// 子公司
type StockSubCompany struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Code     string `gorm:"index" json:"code"`
	Name     string `json:"name"`
	Relation string `json:"relation"` // 参股关系
	Percent  string `json:"percent"`  // 控股占比
	Type     string `json:"type"`     // 业务性质
}

func (StockSubCompany) TableName() string {
	return "magic_stock_sub_company"
}
