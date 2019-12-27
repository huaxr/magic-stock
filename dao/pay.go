// @Contact:    huaxinrui
// @Time:       2019/9/26 下午6:14

package dao

import (
	"magic/stock/core/store"
	"magic/stock/dal"
)

type PayDaoIF interface {
	Create(Pay *dal.Pay) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*dal.Pay, error)
	QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.Pay, error)
	Count(where string, args []interface{}) (int, error)
}

var PayDao PayDaoIF

func init() {
	tmp := new(MysqlPay)
	tmp.Store = store.MysqlClient
	PayDao = tmp
}

type MysqlPay struct {
	Store store.StoreIF
}

func (m *MysqlPay) Create(Pay *dal.Pay) error {
	return m.Store.GetDB().Model(Pay).Create(Pay).Error
}

func (m *MysqlPay) Delete(id int) error {
	return m.Store.GetDB().Delete(&dal.Pay{}, "id = ?", id).Error
}

func (m *MysqlPay) Update(id int, ma map[string]interface{}) error {
	return m.Store.GetDB().Update(&dal.Pay{}, ma).Error
}

func (m *MysqlPay) Query(where string, args []interface{}) (*dal.Pay, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = dal.Pay{}
	query_obj.Where = where
	query_obj.Args = args
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*dal.Pay), nil
}

func (m *MysqlPay) Count(where string, args []interface{}) (int, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = dal.Pay{}
	query_obj.Where = where
	query_obj.Args = args
	return m.Store.Count(query_obj)
}

func (m *MysqlPay) QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.Pay, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = []dal.Pay{}
	query_obj.Where = where
	query_obj.Args = args
	query_obj.Limit = limit
	query_obj.Offset = offset
	query_obj.SelectOnly = select_only
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*[]dal.Pay), nil
}
