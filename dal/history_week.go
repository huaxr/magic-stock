// @Time:       2019/12/2 下午3:59

package dal

type TicketHistoryWeekly struct {
	ID      uint   `gorm:"primary_key"`
	Code    string `gorm:"index"`
	Name    string
	Date    string `gorm:"index"`
	Shou    float64
	Percent float64 // 周线涨跌幅
}

func (TicketHistoryWeekly) TableName() string {
	return "magic_stock_history_week"
}
