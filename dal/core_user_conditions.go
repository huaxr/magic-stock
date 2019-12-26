// @Time:       2019/11/27 下午7:53

package dal

// 用户的条件列表
type UserConditions struct {
	ID         uint   `gorm:"primary_key"`
	Name       string `json:"name"` // 用户为条件定义的名字
	UserId     int    `json:"user_id"`
	Conditions JSON   `sql:"type:text" json:"conditions"`
}

func (UserConditions) TableName() string {
	return "magic_stock_core_user_conditions"
}
