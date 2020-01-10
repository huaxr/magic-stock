// @Time:       2020/1/5 下午5:51

package dal

import "github.com/jinzhu/gorm"

// 用户的需求
type UserDemands struct {
	gorm.Model
	UserId   int    `json:"user_id"`
	Type     string `json:"type"`                     // 平台建议, 开发需求, 好评留言
	Content  string `sql:"type:text" json:"content"`  // 用户需求内容
	Response string `sql:"type:text" json:"response"` // 管理员回复
}

func (UserDemands) TableName() string {
	return "magic_stock_core_user_demand"
}
