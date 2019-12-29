// @Time:       2019/12/2 下午4:00

package dal

type Predict struct {
	ID           int    `gorm:"primary_key" json:"id"`
	Code         string `json:"code"` //`json:"code,omitempty"`
	Name         string `json:"name"`
	Date         string `json:"date"`
	Condition    string `sql:"type:text" json:"condition"` // 分析指标打标
	BadCondition string `sql:"type:text" json:"bad_condition"`
	Finance      string `sql:"type:text" json:"finance"` // 财务指标打标

	FundCount int `json:"fund_count"` // 一共几只基金持有
	SMCount   int `json:"sm_count"`   // 私募数量
	Score     int `json:"score"`      // 得分
	//GMCount     int     `json:"gm_count"`   // 公募数量
	Price   float64 `json:"price"`   // 昨日收盘价格
	Percent float64 `json:"percent"` // 昨日涨跌
}

func (Predict) TableName() string {
	return "magic_stock_predict"
}
