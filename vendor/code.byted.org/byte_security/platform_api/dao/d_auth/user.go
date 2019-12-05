// @Contact:    huaxinrui
// @Time:       2019/9/26 下午6:14

package d_auth

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/platform_api/core"
	"code.byted.org/byte_security/platform_api/core/store"
)

type UserDaoIF interface {
	Create(User *auth.User) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*auth.User, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.User, error)
	Count(where string, args []interface{}) (int, error)
}

var UserDao UserDaoIF

func init() {
	tmp := new(MysqlUser)
	tmp.Store = core.Backend.Store
	UserDao = tmp
}

type MysqlUser struct {
	Store store.StoreIF
}

func InitMysqlUser(i store.StoreIF) *MysqlUser {
	a := new(MysqlUser)
	a.Store = i
	return a
}

func (m *MysqlUser) Create(User *auth.User) error {
	return m.Store.GetDB().Model(User).Create(User).Error
}

func (m *MysqlUser) Delete(id int) error {
	return m.Store.GetDB().Delete(&auth.User{}, "id = ?", id).Error
}

func (m *MysqlUser) Update(id int, ma map[string]interface{}) error {
	return m.Store.GetDB().Update(&auth.User{}, ma).Error
}

func (m *MysqlUser) Query(where string, args []interface{}) (*auth.User, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = auth.User{}
	query_obj.Where = where
	query_obj.Args = args
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*auth.User), nil
}

func (m *MysqlUser) Count(where string, args []interface{}) (int, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = auth.User{}
	query_obj.Where = where
	query_obj.Args = args
	return m.Store.Count(query_obj)
}

func (m *MysqlUser) QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.User, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = []auth.User{}
	query_obj.Where = where
	query_obj.Args = args
	query_obj.Limit = limit
	query_obj.Offset = offset
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*[]auth.User), nil
}
