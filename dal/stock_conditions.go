// @Time:       2019/12/10 下午2:42

package dal

// 条件表
type Conditions struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Type  string `json:"type"`  // 条件类型
	Name  string `json:"name"`  // 条件名称
	Flag  string `json:"flag"`  // 互斥标记
	Score int    `json:"score"` // 分数
}

func (Conditions) TableName() string {
	return "magic_stock_conditions"
}

// 高级条件表
type HighConditions struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Type  string `json:"type"`  // 类型
	Tag   string `json:"tag"`   // 标签
	Name  string `json:"name"`  // 条件名称
	Field string `json:"field"` // 条件的字段名
	Desc  string `json:"desc"`  // 描述
}

func (HighConditions) TableName() string {
	return "magic_stock_conditions_high"
}
