// @Time:       2019/11/27 下午7:53

package dal

// 价格档位
type Price struct {
	ID       uint   `gorm:"primary_key"`
	TypeDesc string `json:"type_desc"` //  类型描述
	Spend    int    `json:"spend"`
	Month    int    `json:"month"` // 几个月
}

func (Price) TableName() string {
	return "magic_stock_core_price"
}
