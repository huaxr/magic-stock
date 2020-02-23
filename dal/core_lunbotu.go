// @Time:       2020/2/23 下午3:57

package dal

// 轮播图
type LunBoTu struct {
	ID      uint   `gorm:"primary_key"`
	Pic     string `json:"pic"`
	Content string `json:"content"`
	Link    string `json:"link"`
	Disable bool   `json:"disable"`
}

func (LunBoTu) TableName() string {
	return "magic_stock_core_lunbotu"
}
