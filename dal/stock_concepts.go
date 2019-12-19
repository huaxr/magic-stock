// @Time:       2019/12/9 下午2:17

package dal

// 概念股
type StockConcept struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

func (StockConcept) TableName() string {
	return "magic_stock_concept"
}
