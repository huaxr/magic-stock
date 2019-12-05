// @Contact:    huaxinrui
// @Time:       2019/9/26 下午6:14

package d_auth

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/platform_api/core"
	"code.byted.org/byte_security/platform_api/core/store"
)

type CasbinDaoIF interface {
	Create(Casbin *auth.CasbinRule) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*auth.CasbinRule, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.CasbinRule, error)
	Count(where string, args []interface{}) (int, error)
}

var CasbinDao CasbinDaoIF

func init() {
	tmp := new(MysqlCasbin)
	tmp.Store = core.Backend.Store
	CasbinDao = tmp
}

type MysqlCasbin struct {
	Store store.StoreIF
}

func InitMysqlCasbin(i store.StoreIF) *MysqlCasbin {
	a := new(MysqlCasbin)
	a.Store = i
	return a
}

func (m *MysqlCasbin) Create(Casbin *auth.CasbinRule) error {
	return m.Store.GetDB().Model(Casbin).Create(Casbin).Error
}

func (m *MysqlCasbin) Delete(id int) error {
	return m.Store.GetDB().Delete(&auth.CasbinRule{}, "id = ?", id).Error
}

func (m *MysqlCasbin) Update(id int, ma map[string]interface{}) error {
	return m.Store.GetDB().Update(&auth.CasbinRule{}, ma).Error
}

func (m *MysqlCasbin) Query(where string, args []interface{}) (*auth.CasbinRule, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = auth.CasbinRule{}
	query_obj.Where = where
	query_obj.Args = args
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*auth.CasbinRule), nil
}

func (m *MysqlCasbin) Count(where string, args []interface{}) (int, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = auth.CasbinRule{}
	query_obj.Where = where
	query_obj.Args = args
	return m.Store.Count(query_obj)
}

func (m *MysqlCasbin) QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.CasbinRule, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = []auth.CasbinRule{}
	query_obj.Where = where
	query_obj.Args = args
	query_obj.Limit = limit
	query_obj.Offset = offset
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*[]auth.CasbinRule), nil
}
