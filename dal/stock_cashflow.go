// @Time:       2019/12/9 下午2:17

package dal

// 现金流量表
type StockCashFlow struct {
	ID                  uint    `gorm:"primary_key"`
	Code                string  `gorm:"index"`
	ManageCashFlow      float64 // 经营活动产生的现金流量净额 (销售商品、提供劳务收到的现金 + 收到的税费返还 +收到的其他与经营活动有关的现金-支付的各项税费-经营活动现金流出小计....)
	InvestCashFlow      float64 //投资活动产生的现金流量净额 (收回投资所收到的现金 + 取得投资收益所收到的现金 - 投资所支付的现金...)
	FundraisingCashFlow float64 // 筹资活动产生的现金流量净额 (取得借款收到的现金 + 发行债券收到的现金- 偿还债务支付的现金 - 支付其他与筹资活动有关的现金)
	CashRemain          float64 //期末现金及现金等价物余额
	Date                string
}

// 资产负债表
func (StockCashFlow) TableName() string {
	return "magic_stock_cashflow"
}
