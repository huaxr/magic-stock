// @Time:       2020/1/5 下午5:51

package dal

import "github.com/jinzhu/gorm"

// 用户的需求
type UserSelect struct {
	gorm.Model
	UserId int     `json:"user_id"`
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

func (UserSelect) TableName() string {
	return "magic_stock_core_user_select"
}
