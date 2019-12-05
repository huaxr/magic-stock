// @Time:       2019/12/2 下午2:30

package dal

type TicketHistory struct {
	ID           uint   `gorm:"primary_key"`
	Code         string `gorm:"index"`
	Name         string
	Date         string `gorm:"index"`
	Kai          float64
	High         float64
	Shou         float64
	Low          float64
	TotalCount   float64 // 成交量 手
	TotalMoney   float64 // 成交额 万元
	Percent      float64 // 涨跌幅 %
	Change       float64 // 涨跌额
	Amplitude    float64 // 振幅
	TurnoverRate float64 // 换手率
}

func (TicketHistory) TableName() string {
	return "magic_stock_history"
}
