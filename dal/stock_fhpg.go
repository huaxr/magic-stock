// @Time:       2019/12/30 下午7:55

package dal

// 分红配股表
type StockFengHong struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	Code       string `gorm:"index" json:"code"`
	Date       string `json:"date"`
	SongGu     string `json:"song_gu"`     // 每10股送股个数
	ZhuangZeng string `json:"zhuang_zeng"` // 转增
	PaiXi      string `json:"pai_xi"`      // 派息税前 元
	Process    string `json:"process"`     // 进度
}

func (StockFengHong) TableName() string {
	return "magic_stock_fenghong"
}

// 配股
type StockPeiGu struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	Code   string `gorm:"index" json:"code"`
	Date   string `json:"date"`
	Count  string `json:"count"`  // 配股方案(每10股配股股数)
	Price  string `json:"price"`  // 配股价格(元)
	Number string `json:"number"` // 基准股本(股)
}

func (StockPeiGu) TableName() string {
	return "magic_stock_peigu"
}

// 增发
type StockZengFa struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Code      string `gorm:"index" json:"code"`
	Date      string `json:"date"`
	Way       string `json:"way"`        // 发行方式
	Price     string `json:"price"`      // 发行价格
	AllPrice  string `json:"all_price"`  // 实际公司募集资金总额
	CostPrice string `json:"cost_price"` // 发行费用总额：
	AllCount  string `json:"all_count"`  // 实际发行数量
}

func (StockZengFa) TableName() string {
	return "magic_stock_zengfa"
}
