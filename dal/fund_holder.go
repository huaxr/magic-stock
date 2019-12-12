// @Time:       2019/12/2 下午4:24

package dal

// 基金持仓
type FundHoldRank struct {
	ID        uint    `gorm:"primary_key" json:"id"`
	Type      string  `json:"type"` // 基金类型 lof ...
	FundCode  string  `gorm:"index" json:"fund_code"`
	FundSName string  `json:"fund_s_name"`
	Code      string  `gorm:"index" json:"code"`
	Name      string  `json:"name"`
	Percent   float64 `json:"percent"`
	Time      string  `json:"time"`
}

func (FundHoldRank) TableName() string {
	return "magic_stock_fund_hold"
}
