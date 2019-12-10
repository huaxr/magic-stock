// @Time:       2019/12/2 下午4:00

package dal

type Predict struct {
	ID          int `gorm:"primary_key"`
	Code        string
	Name        string
	Date        string
	Condition   string  `sql:"type:text"`
	RealPercent float64 // 真实的涨跌情况， 用来判断准确性
	FundCount   int     // 一共几只基金持有
	SMCount     int     // 私募数量
	GMCount     int     // 公募数量
}

func (Predict) TableName() string {
	return "magic_stock_predict"
}
