// @Time:       2019/12/2 下午4:15

package dal

type Stockholder struct {
	ID         uint   `gorm:"primary_key"`
	Code       string `gorm:"index"`
	Name       string
	HolderName string
	Count      string
	Percent    string
	Change     string
}

func (Stockholder) TableName() string {
	return "magic_stock_holder"
}
