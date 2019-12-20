// @Time:       2019/11/27 下午7:53

package dal

// 提需求
type Demand struct {
	ID      uint   `gorm:"primary_key"`
	Content string `json:"content"` // 需求内容
	Status  string `json:"status"`  // 处理状态
}

func (Demand) TableName() string {
	return "magic_stock_core_demand"
}
