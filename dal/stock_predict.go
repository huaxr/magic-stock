// @Time:       2019/12/2 下午4:00

package dal

type Predict struct {
	ID              int     `gorm:"primary_key" json:"id"`
	Code            string  `json:"code"` //`json:"code,omitempty"`
	Name            string  `json:"name"`
	Date            string  `json:"date"`
	Condition       string  `sql:"type:text" json:"condition"` // 分析指标打标
	BadCondition    string  `sql:"type:text" json:"bad_condition"`
	Condition_      string  `sql:"type:text" json:"condition_"`
	BadCondition_   string  `sql:"type:text" json:"bad_condition_"`
	Finance         string  `sql:"type:text" json:"finance"` // 财务指标打标
	FundCount       int     `json:"fund_count"`              // 一共几只基金持有
	SMCount         int     `json:"sm_count"`                // 私募数量
	FenghongCount   int     `json:"fenghong_count"`          // 分红次数
	SongguCount     int     `json:"songgu_count"`            // 送股次数
	ZhuangzengCount int     `json:"zhuangzeng_count"`        // 转增次数
	PeiguCount      int     `json:"peigu_count"`             // 配股次数
	ZengfaCount     int     `json:"zengfa_count"`            // 增发次数
	SubcompCount    int     `json:"subcomp_count"`           // 参股公司的数量
	Score           int     `json:"score"`                   // 得分
	Price           float64 `json:"price"`                   // 昨日收盘价格
	Percent         float64 `json:"percent"`                 // 昨日涨跌
}

func (Predict) TableName() string {
	return "magic_stock_predict"
}
