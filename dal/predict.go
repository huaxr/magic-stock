// @Time:       2019/12/2 下午4:00

package dal

type Predict struct {
	ID          int     `gorm:"primary_key" json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Date        string  `json:"date"`
	Condition   string  `sql:"type:text" json:"condition"`
	RealPercent float64 `json:"real_percent"` // 真实的涨跌情况， 用来判断准确性
	FundCount   int     `json:"fund_count"`   // 一共几只基金持有
	SMCount     int     `json:"sm_count"`     // 私募数量
	//GMCount     int     `json:"gm_count"`     // 公募数量
}

func (Predict) TableName() string {
	return "magic_stock_predict"
}
