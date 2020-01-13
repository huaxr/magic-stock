// @Time:       2020/1/13 下午5:56

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

// 概念股
type StockConcept struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

func (StockConcept) TableName() string {
	return "magic_stock_concept"
}

// 股票其它标签
type StockLabels struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

func (StockLabels) TableName() string {
	return "magic_stock_labels"
}

// 条件表
type Conditions struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Type  string `json:"type"`  // 条件类型
	Name  string `json:"name"`  // 条件名称
	Flag  string `json:"flag"`  // 互斥标记
	Score int    `json:"score"` // 分数
}

func (Conditions) TableName() string {
	return "magic_stock_conditions"
}

// 高级条件表
type HighConditions struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Type  string `json:"type"`  // 类型
	Tag   string `json:"tag"`   // 标签
	Name  string `json:"name"`  // 条件名称
	Field string `json:"field"` // 条件的字段名
	Desc  string `json:"desc"`  // 描述
}

func (HighConditions) TableName() string {
	return "magic_stock_conditions_high"
}
