// @Time:       2019/12/2 下午4:15

package dal

type Stockholder struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	Code       string `gorm:"index" json:"code"`
	Name       string `json:"name"`
	HolderName string `json:"holder_name"`
	Count      string `json:"count"`
	Percent    string `json:"percent"`
	Change     string `json:"change"`
}

func (Stockholder) TableName() string {
	return "magic_stock_holder"
}
