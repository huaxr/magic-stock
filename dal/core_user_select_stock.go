// @Time:       2020/1/5 下午5:51

package dal

import "time"

// 用户的自选
type UserSelect struct {
	ID          uint      `gorm:"primary_key"`
	UserId      int       `json:"user_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	CreatedTime time.Time `json:"created_time"`
}

func (UserSelect) TableName() string {
	return "magic_stock_core_user_select"
}
