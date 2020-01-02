// @Time:       2019/11/27 下午7:53

package dal

import (
	"github.com/jinzhu/gorm"
)

// 支付记录
type Pay struct {
	gorm.Model
	UserId     int    `json:"user_id"`     // 支付用户
	Spend      int    `json:"spend"`       // 支付总金额
	Type       string `json:"type"`        // 支付类型 年费至尊
	Month      int    `json:"month"`       // 月数
	PaySuccess bool   `json:"pay_success"` // 支付成功
	OrderId    string `json:"order_id"`    // 订单号
}

func (Pay) TableName() string {
	return "magic_stock_core_pay"
}
