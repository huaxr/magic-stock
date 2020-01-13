// @Time:       2019/12/2 下午3:59

package dal

type TicketHistoryMonth struct {
	ID           uint    `gorm:"primary_key" json:"id"`
	Code         string  `gorm:"index" json:"code"`
	Name         string  `json:"name"`
	Date         string  `gorm:"index" json:"date"`
	Kai          float64 `json:"kai"`
	High         float64 `json:"high"`
	Shou         float64 `json:"shou"`
	Low          float64 `json:"low"`
	TotalCount   float64 `json:"total_count"`   // 成交量 手
	Percent      float64 `json:"percent"`       // 涨跌幅 %
	Change       float64 `json:"change"`        // 涨跌额
	Amplitude    float64 `json:"amplitude"`     // 振幅
	TurnoverRate float64 `json:"turnover_rate"` // 换手率
}

func (TicketHistoryMonth) TableName() string {
	return "magic_stock_history_month" // 全部的月线数据
}
