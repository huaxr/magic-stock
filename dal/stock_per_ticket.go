// @Time:       2019/12/2 下午4:15

package dal

// 每股指标
type StockPerTicket struct {
	ID             uint    `gorm:"primary_key" json:"id"`
	Code           string  `gorm:"index" json:"code"`
	Tanboshouyi    float64 `json:"tanboshouyi"`    // 摊薄每股收益(元)
	Jiaquanshouyi  float64 `json:"jiaquanshouyi"`  // 加权每股收益(元)
	Jinzichanfront float64 `json:"jinzichanfront"` // 每股净资产_调整前(元)

	Shouyiafter        float64 `json:"shouyiafter"`        // 每股收益_调整后(元)
	Jinzichanafter     float64 `json:"jinzichanafter"`     // 每股净资产_调整后(元)
	Jingyingxianjinliu float64 `json:"jingyingxianjinliu"` // 每股经营性现金流(元)
	Gubengongjijin     float64 `json:"gubengongjijin"`     // 每股资本公积金(元)
	Weifenpeilirun     float64 `json:"weifenpeilirun"`     // 每股未分配利润(元)
	Date               string  `json:"date"`
}

func (StockPerTicket) TableName() string {
	return "magic_stock_per_ticket"
}
