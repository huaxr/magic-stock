// @Time:       2019/12/10 下午2:42

package dal

// 条件表
type Conditions struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Type string `json:"type"` // 条件类型
	Name string `json:"name"` // 条件名称
}

func (Conditions) TableName() string {
	return "magic_stock_conditions"
}
