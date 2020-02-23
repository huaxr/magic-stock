// @Time:       2020/2/23 上午10:18

package dal

// 融资融券

type RzRq struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Code string `gorm:"index" json:"code"`
	Name string `json:"name"`

	// 融资
	Balance float64 `json:"balance"`  // 融资余额 元
	Buy     float64 `json:"buy"`      // 买入额
	PayBack float64 `json:"pay_back"` // 偿还额

	// 融券
	YlYe         float64 `json:"yl_ye"`          //余量金额(元)
	AllCount     float64 `json:"all_count"`      // 余量 股
	SaleCount    float64 `json:"sale_count"`     // 卖出量(股)
	PayBackCount float64 `json:"pay_back_count"` // 偿还量(股)
	RqYe         float64 `json:"rq_ye"`          //融券余额(元)
	Date         string  `gorm:"index" json:"date"`
}

func (RzRq) TableName() string {
	return "magic_stock_rzrq"
}
