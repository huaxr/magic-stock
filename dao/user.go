// @Contact:    huaxinrui
// @Time:       2019/9/26 下午6:14

package dao

import (
	"magic/stock/core/store"
	"magic/stock/dal"
)

type UserDaoIF interface {
	Create(User *dal.User) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*dal.User, error)
	QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.User, error)
	Count(where string, args []interface{}) (int, error)
}

var UserDao UserDaoIF

func init() {
	tmp := new(MysqlUser)
	tmp.Store = store.MysqlClient
	UserDao = tmp
}

type MysqlUser struct {
	Store store.StoreIF
}

func (m *MysqlUser) Create(User *dal.User) error {
	return m.Store.GetDB().Model(User).Create(User).Error
}

func (m *MysqlUser) Delete(id int) error {
	return m.Store.GetDB().Delete(&dal.User{}, "id = ?", id).Error
}

func (m *MysqlUser) Update(id int, ma map[string]interface{}) error {
	return m.Store.GetDB().Update(&dal.User{}, ma).Error
}

func (m *MysqlUser) Query(where string, args []interface{}) (*dal.User, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = dal.User{}
	query_obj.Where = where
	query_obj.Args = args
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*dal.User), nil
}

func (m *MysqlUser) Count(where string, args []interface{}) (int, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = dal.User{}
	query_obj.Where = where
	query_obj.Args = args
	return m.Store.Count(query_obj)
}

func (m *MysqlUser) QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.User, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = []dal.User{}
	query_obj.Where = where
	query_obj.Args = args
	query_obj.Limit = limit
	query_obj.Offset = offset
	query_obj.SelectOnly = select_only
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*[]dal.User), nil
}
