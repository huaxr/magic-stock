// @Time:       2019/12/9 上午11:32

package dal

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
