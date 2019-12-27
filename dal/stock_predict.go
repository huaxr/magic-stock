// @Time:       2019/12/2 下午4:00

package dal

type Predict struct {
	ID           int    `gorm:"primary_key" json:"id,omitempty"`
	Code         string `json:"code,omitempty"`
	Name         string `json:"name,omitempty"`
	Date         string `json:"date,omitempty"`
	Condition    string `sql:"type:text" json:"condition,omitempty"` // 分析指标打标
	BadCondition string `sql:"type:text" json:"bad_condition,omitempty"`
	Finance      string `sql:"type:text" json:"finance,omitempty"` // 财务指标打标
	FundCount    int    `json:"fund_count,omitempty"`              // 一共几只基金持有
	SMCount      int    `json:"sm_count,omitempty"`                // 私募数量
	Score        int    `json:"score,omitempty"`                   // 得分
	//GMCount     int     `json:"gm_count"`   // 公募数量
	Price   float64 `json:"price,omitempty"`   // 昨日收盘价格
	Percent float64 `json:"percent,omitempty"` // 昨日涨跌
}

func (Predict) TableName() string {
	return "magic_stock_predict"
}
