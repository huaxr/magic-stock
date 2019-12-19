// @Time:       2019/12/19 上午11:49

package dal

// 股票其它标签
type StockLabels struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

func (StockLabels) TableName() string {
	return "magic_stock_labels"
}
