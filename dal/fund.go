// @Time:       2019/12/2 下午4:21

package dal

// 基金排行
type FundRank struct {
	ID         uint    `gorm:"primary_key" json:"id"`
	Type       string  `json:"type"`
	FundCode   string  `gorm:"index" json:"fund_code"`
	FundSName  string  `json:"fund_s_name"`
	LastMonth  float64 `json:"last_month"`
	LastWeek   float64 `json:"last_week"`
	Last3Month float64 `json:"last_3_month"`
	Last6Month float64 `json:"last_6_month"`
	SinceBase  float64 `json:"since_base"` // 成立以来涨跌幅
	ThisYear   float64 `json:"this_year"`
	LastYear   float64 `json:"last_year"`
	Last2Year  float64 `json:"last_2_year"`
	Last3Year  float64 `json:"last_3_year"`
	Time       string  `json:"time"`
}

func (FundRank) TableName() string {
	return "magic_stock_fund_rank"
}
