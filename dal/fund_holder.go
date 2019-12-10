// @Time:       2019/12/2 下午4:24

package dal

// 基金持仓
type FundHoldRank struct {
	ID        uint   `gorm:"primary_key"`
	Type      string // 基金类型 lof ...
	FundCode  string `gorm:"index"`
	FundSName string
	Code      string `gorm:"index"`
	Name      string
	Percent   float64
	Time      string
}

func (FundHoldRank) TableName() string {
	return "magic_stock_fund_hold"
}
