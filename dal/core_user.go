// @Time:       2019/11/27 下午7:53

package dal

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	OpenId           string    `json:"open_id"`            // openid
	UserName         string    `json:"user_name"`          // 姓名
	Avatar           string    `json:"avatar"`             // 头像
	Address          string    `json:"address"`            // 住址
	Sex              int       `json:"sex"`                // 性别  1==male
	City             string    `json:"city"`               // 城市
	Province         string    `json:"province"`           // 省份
	Country          string    `json:"country"`            // 国家
	Phone            string    `json:"phone"`              // 电话 (下单前必填)
	Exp              int       `json:"exp"`                // 经验值
	QueryLeft        int       `json:"query_left"`         // 剩余查询次数 （没加入一个新用户加5次？）
	MemberExpireTime time.Time `json:"member_expire_time"` // 会员失效时间
	ShareToken       string    `json:"token"`              // 分享token
}

func (User) TableName() string {
	return "magic_stock_core_user"
}
