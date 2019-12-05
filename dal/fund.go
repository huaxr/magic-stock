// @Time:       2019/12/2 下午4:21

package dal

// 基金排行
type FundRank struct {
	ID         uint `gorm:"primary_key"`
	Type       string
	FundCode   string `gorm:"index"`
	FundSName  string
	LastMonth  float64
	LastWeek   float64
	Last3Month float64
	Last6Month float64
	SinceBase  float64 // 成立以来涨跌幅
	ThisYear   float64
	LastYear   float64
	Last2Year  float64
	Last3Year  float64
	Time       string
}

func (FundRank) TableName() string {
	return "magic_stock_fund_rank"
}
