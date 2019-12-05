// @Contact:    huaxinrui
// @Time:       2019/9/26 下午9:25

package s_auth

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/platform_api/dao/d_auth"
)

type UserServiceIF interface {
	Create(event *auth.User) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*auth.User, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.User, error)
	Count(where string, args []interface{}) (int, error)
}

var UserServiceGlobal UserServiceIF

func init() {
	tmp := new(UserService)
	tmp.dao = d_auth.UserDao
	UserServiceGlobal = tmp
}

type UserService struct {
	dao d_auth.UserDaoIF
}

func (m *UserService) Create(app *auth.User) error {
	return m.dao.Create(app)
}

func (m *UserService) Delete(id int) error {
	return m.dao.Delete(id)
}

func (m *UserService) Update(id int, ma map[string]interface{}) error {
	return m.dao.Update(id, ma)
}

func (m *UserService) Query(where string, args []interface{}) (*auth.User, error) {
	return m.dao.Query(where, args)
}

func (m *UserService) QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.User, error) {
	return m.dao.QueryAll(where, args, offset, limit)
}

func (m *UserService) Count(where string, args []interface{}) (int, error) {
	return m.dao.Count(where, args)
}
