// @Time:       2019/11/27 下午7:53

package dal

import (
	"github.com/jinzhu/gorm"
)

// 用户拉活列表
type UserShare struct {
	gorm.Model      // 加入时间等
	ShareUserId int `json:"share_user_id"`
	BeShareId   int `json:"be_share_id"`
}

func (UserShare) TableName() string {
	return "magic_stock_core_user_share"
}
