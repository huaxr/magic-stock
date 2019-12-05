// @Contact:    huaxinrui
// @Time:       2019/9/26 下午6:14

package d_auth

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/platform_api/core"
	"code.byted.org/byte_security/platform_api/core/store"
)

type GroupDaoIF interface {
	Create(Group *auth.Group) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*auth.Group, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.Group, error)
	Count(where string, args []interface{}) (int, error)
}

var GroupDao GroupDaoIF

func init() {
	tmp := new(MysqlGroup)
	tmp.Store = core.Backend.Store
	GroupDao = tmp
}

type MysqlGroup struct {
	Store store.StoreIF
}

func InitMysqlGroup(i store.StoreIF) *MysqlGroup {
	a := new(MysqlGroup)
	a.Store = i
	return a
}

func (m *MysqlGroup) Create(Group *auth.Group) error {
	return m.Store.GetDB().Model(Group).Create(Group).Error
}

func (m *MysqlGroup) Delete(id int) error {
	return m.Store.GetDB().Delete(&auth.Group{}, "id = ?", id).Error
}

func (m *MysqlGroup) Update(id int, ma map[string]interface{}) error {
	return m.Store.GetDB().Update(&auth.Group{}, ma).Error
}

func (m *MysqlGroup) Query(where string, args []interface{}) (*auth.Group, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = auth.Group{}
	query_obj.Where = where
	query_obj.Args = args
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*auth.Group), nil
}

func (m *MysqlGroup) QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.Group, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = []auth.Group{}
	query_obj.Where = where
	query_obj.Args = args
	query_obj.Limit = limit
	query_obj.Offset = offset
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*[]auth.Group), nil
}

func (m *MysqlGroup) Count(where string, args []interface{}) (int, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = auth.Group{}
	query_obj.Where = where
	query_obj.Args = args
	return m.Store.Count(query_obj)
}
