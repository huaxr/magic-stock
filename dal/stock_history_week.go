// @Time:       2019/12/2 下午3:59

package dal

type TicketHistoryWeekly struct {
	ID      uint    `gorm:"primary_key" json:"id"`
	Code    string  `gorm:"index" json:"code"`
	Name    string  `json:"name"`
	Date    string  `gorm:"index" json:"date"`
	Shou    float64 `json:"shou"`
	Percent float64 `json:"percent"` // 周线涨跌幅
}

func (TicketHistoryWeekly) TableName() string {
	return "magic_stock_history_week"
}
