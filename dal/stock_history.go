// @Time:       2019/12/2 下午2:30

package dal

type TicketHistory struct {
	ID           uint    `gorm:"primary_key" json:"id"`
	Code         string  `gorm:"index" json:"code"`
	Name         string  `json:"name"`
	Date         string  `gorm:"index" json:"date"`
	Kai          float64 `json:"kai"`
	High         float64 `json:"high"`
	Shou         float64 `json:"shou"`
	Low          float64 `json:"low"`
	TotalCount   float64 `json:"total_count"`   // 成交量 手
	TotalMoney   float64 `json:"total_money"`   // 成交额 万元
	Percent      float64 `json:"percent"`       // 涨跌幅 %
	Change       float64 `json:"change"`        // 涨跌额
	Amplitude    float64 `json:"amplitude"`     // 振幅
	TurnoverRate float64 `json:"turnover_rate"` // 换手率
	NumberRate   float64 `json:"number_rate"`   // 量比
}

func (TicketHistory) TableName() string {
	return "magic_stock_history"
}

type HistoryALL1 struct {
	ID           uint    `gorm:"primary_key" json:"id"`
	Code         string  `gorm:"index" json:"code"`
	Name         string  `json:"name"`
	Date         string  `gorm:"index" json:"date"`
	Kai          float64 `json:"kai"`
	High         float64 `json:"high"`
	Shou         float64 `json:"shou"`
	Low          float64 `json:"low"`
	TotalCount   float64 `json:"total_count"`   // 成交量 手
	TotalMoney   float64 `json:"total_money"`   // 成交额 万元
	Percent      float64 `json:"percent"`       // 涨跌幅 %
	Change       float64 `json:"change"`        // 涨跌额
	Amplitude    float64 `json:"amplitude"`     // 振幅
	TurnoverRate float64 `json:"turnover_rate"` // 换手率
	NumberRate   float64 `json:"number_rate"`   // 量比
}

func (HistoryALL1) TableName() string {
	return "magic_stock_history_lte_1500" // <= 1500 的id的股票
}

type HistoryALL2 struct {
	ID           uint    `gorm:"primary_key" json:"id"`
	Code         string  `gorm:"index" json:"code"`
	Name         string  `json:"name"`
	Date         string  `gorm:"index" json:"date"`
	Kai          float64 `json:"kai"`
	High         float64 `json:"high"`
	Shou         float64 `json:"shou"`
	Low          float64 `json:"low"`
	TotalCount   float64 `json:"total_count"`   // 成交量 手
	TotalMoney   float64 `json:"total_money"`   // 成交额 万元
	Percent      float64 `json:"percent"`       // 涨跌幅 %
	Change       float64 `json:"change"`        // 涨跌额
	Amplitude    float64 `json:"amplitude"`     // 振幅
	TurnoverRate float64 `json:"turnover_rate"` // 换手率
	NumberRate   float64 `json:"number_rate"`   // 量比
}

func (HistoryALL2) TableName() string {
	return "magic_stock_history_gte_1500"
}

type TicketHistoryWeekly struct {
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

func (TicketHistoryWeekly) TableName() string {
	return "magic_stock_history_week"
}

type TicketHistoryWeeklyALL struct {
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

func (TicketHistoryWeeklyALL) TableName() string {
	return "magic_stock_history_week_all" // 全部的周线数据
}

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
	return "magic_stock_history_month" //
}

type TicketHistoryMonthAll struct {
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

func (TicketHistoryMonthAll) TableName() string {
	return "magic_stock_history_month_all" // 全部的月线数据
}
