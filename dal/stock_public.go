// @Time:       2019/12/2 下午4:15

package dal

// 公告新闻
type StockPublicNews struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Code  string `json:"code"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Time  string `json:"time"`
	Type  string `json:"type"`
}

func (StockPublicNews) TableName() string {
	return "magic_stock_public_news"
}

// 季度报告
type StockPublicReports struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Code  string `json:"code"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Time  string `json:"time"`
	Type  string `json:"type"`
}

func (StockPublicReports) TableName() string {
	return "magic_stock_public_reports"
}
