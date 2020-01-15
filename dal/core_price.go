// @Time:       2019/11/27 下午7:53

package dal

// 价格档位
type Price struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Type     string `json:"type"`      //  类型 member data query
	TypeDesc string `json:"type_desc"` //  类型描述
	Spend    int    `json:"spend"`
	Count    int    `json:"count"` // 几个月 几次 等
}

func (Price) TableName() string {
	return "magic_stock_core_price"
}
