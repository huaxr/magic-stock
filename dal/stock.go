// @Time:       2019/12/2 下午2:28

package dal

type Code struct {
	ID      uint   `gorm:"primary_key"`
	Code    string `gorm:"not null;unique"`
	Name    string
	Belong  string // 所属行业板块
	Concept string // 所属概念板块
}

func (Code) TableName() string {
	return "magic_stock_code"
}
