// @Time:       2019/12/1 下午4:19

package adapter

import (
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/dao"
	"time"
)

type PayServiceIF interface {
	Create(event *dal.Pay) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*dal.Pay, error)
	QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.Pay, error)
	Count(where string, args []interface{}) (int, error)
	// 更新支付标志位, 发送短信, 生成真实订单数据
	UpdatePaySuccessAndGenerateIndent(order_id string)
}

var PayServiceGlobal PayServiceIF

func init() {
	tmp := new(PayService)
	tmp.dao = dao.PayDao
	PayServiceGlobal = tmp
}

type PayService struct {
	dao dao.PayDaoIF
}

func (m *PayService) Create(Pay *dal.Pay) error {
	return m.dao.Create(Pay)
}

func (m *PayService) Delete(id int) error {
	return m.dao.Delete(id)
}

func (m *PayService) Update(id int, ma map[string]interface{}) error {
	return m.dao.Update(id, ma)
}

func (m *PayService) Query(where string, args []interface{}) (*dal.Pay, error) {
	return m.dao.Query(where, args)
}

func (m *PayService) QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.Pay, error) {
	return m.dao.QueryAll(where, args, offset, limit, select_only)
}

func (m *PayService) Count(where string, args []interface{}) (int, error) {
	return m.dao.Count(where, args)
}

func (m *PayService) UpdatePaySuccessAndGenerateIndent(order_id string) {
	payment, err := m.Query("order_id = ? and pay_success = ?", []interface{}{order_id, false})
	if err != nil {
		log.Println("该订单回调不存在")
		return
	}
	payment.PaySuccess = true
	store.MysqlClient.GetDB().Save(payment)

	// update user obj 续充和开通
	user, _ := UserServiceGlobal.Query("id = ?", []interface{}{payment.UserId})
	now := time.Now()
	if now.Before(user.MemberExpireTime) {
		user.MemberExpireTime = now.AddDate(0, 1*payment.Month, 0)
	} else {
		user.MemberExpireTime = user.MemberExpireTime.AddDate(0, 1*payment.Month, 0)

	}
	store.MysqlClient.GetDB().Save(user)
}
