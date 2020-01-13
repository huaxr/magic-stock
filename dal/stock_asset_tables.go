// @Time:       2019/12/9 下午2:17

package dal

// 现金流量表
type StockCashFlow struct {
	ID                  uint    `gorm:"primary_key" json:"id"`
	Code                string  `gorm:"index" json:"code"`
	ManageCashFlow      float64 `json:"manage_cash_flow"`      // 经营活动产生的现金流量净额 (销售商品、提供劳务收到的现金 + 收到的税费返还 +收到的其他与经营活动有关的现金-支付的各项税费-经营活动现金流出小计....)
	InvestCashFlow      float64 `json:"invest_cash_flow"`      //投资活动产生的现金流量净额 (收回投资所收到的现金 + 取得投资收益所收到的现金 - 投资所支付的现金...)
	FundraisingCashFlow float64 `json:"fundraising_cash_flow"` // 筹资活动产生的现金流量净额 (取得借款收到的现金 + 发行债券收到的现金- 偿还债务支付的现金 - 支付其他与筹资活动有关的现金)
	CashRemain          float64 `json:"cash_remain"`           //期末现金及现金等价物余额
	Date                string  `json:"date"`
}

// 资产负债表
func (StockCashFlow) TableName() string {
	return "magic_stock_cashflow"
}

type StockLiabilities struct {
	ID                    uint    `gorm:"primary_key" json:"id"`
	Code                  string  `gorm:"index" json:"code"`
	CurrentAssets         float64 `json:"current_assets"`          // 流动资产合计 （活动资金+交易性金融资产+应收票据+应收款账+存活货..）
	NotCurrentAssets      float64 `json:"not_current_assets"`      // 非流动资产合计 （长期应收款+在建工程+固定资产+商誉+开发支出...）
	TotalAssets           float64 `json:"total_assets"`            // 资产总计
	CurrentLiabilities    float64 `json:"current_liabilities"`     // 流动负债合计  (短期借贷+交易性金融负债+应付票据+应付税费，利息，股利...)
	NotCurrentLiabilities float64 `json:"not_current_liabilities"` // 非流动负债合计 (长期借贷， 应付债券+长期职工薪酬+长期延递收益...)
	TotalLiabilities      float64 `json:"total_liabilities"`       // 负债合计
	Date                  string  `json:"date"`
}

// 资产负债表
func (StockLiabilities) TableName() string {
	return "magic_stock_liabilities"
}

// http://vip.stock.finance.sina.com.cn/corp/go.php/vFD_ProfitStatement/stockid/600536/ctrl/part/displaytype/4.phtml
type StockProfit struct {
	ID                 uint    `gorm:"primary_key" json:"id"`
	Code               string  `gorm:"index" json:"code"`
	GrossTradingIncome float64 `json:"gross_trading_income"` // 营业总收入
	TotalOperatingCost float64 `json:"total_operating_cost"` // 营业总成本
	NetProfit          float64 `json:"net_profit"`           // 净利润 = 利润 + 营业外收入 - 营业外支出 - 所得税
	Date               string  `json:"date"`
}

// 利润表
func (StockProfit) TableName() string {
	return "magic_stock_profit"
}

// 财务报告(Financial Statements) 资产负债表(The Balance Sheet) 现金流量报表(The Income Statement)
