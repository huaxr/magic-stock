// @Time:       2019/12/2 下午4:15

package dal

// http://vip.stock.finance.sina.com.cn/corp/go.php/vFD_ProfitStatement/stockid/600536/ctrl/part/displaytype/4.phtml
type StockProfit struct {
	ID                 uint    `gorm:"primary_key"`
	Code               string  `gorm:"index"`
	GrossTradingIncome float64 // 营业总收入
	TotalOperatingCost float64 // 营业总成本
	NetProfit          float64 // 净利润 = 利润 + 营业外收入 - 营业外支出 - 所得税
	Date               string
}

// 利润表
func (StockProfit) TableName() string {
	return "magic_stock_profit"
}

// 财务报告(Financial Statements) 资产负债表(The Balance Sheet) 现金流量报表(The Income Statement)
