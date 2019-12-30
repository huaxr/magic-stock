// @Time:       2019/12/30 上午11:09

package dal

type StockFund struct {
	ID                 uint   `gorm:"primary_key" json:"id"`
	Code               string `gorm:"index" json:"code"` // 股票代码
	FundCode           string `gorm:"index" json:"fund_code"`
	FundName           string `json:"fund_name"`
	Count              string `json:"count"`                // 持股数量(股)
	PercentLiutong     string `json:"percent_liutong"`      // 占本基金所持流通股比例(%)
	Change             string `json:"change"`               // 持股变化
	Price              string `json:"price"`                // 持股市值(万元)
	PercentJingzhi     string `json:"percent_jingzhi"`      // 占净值比例(%)
	PercentSignalStock string `json:"percent_signal_stock"` // 占个股流通市值比例(%)
	Date               string `json:"date"`                 // 日期
}

func (StockFund) TableName() string {
	return "magic_stock_stock_fund"
}
