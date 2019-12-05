// @Contact:    huaxinrui
// @Time:       2019/9/26 下午6:14

package d_auth

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/platform_api/core"
	"code.byted.org/byte_security/platform_api/core/store"
)

type TokenDaoIF interface {
	Create(Token *auth.Token) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*auth.Token, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.Token, error)
	Count(where string, args []interface{}) (int, error)
}

var TokenDao TokenDaoIF

func init() {
	tmp := new(MysqlToken)
	tmp.Store = core.Backend.Store
	TokenDao = tmp
}

type MysqlToken struct {
	Store store.StoreIF
}

func InitMysqlToken(i store.StoreIF) *MysqlToken {
	a := new(MysqlToken)
	a.Store = i
	return a
}

func (m *MysqlToken) Create(Token *auth.Token) error {
	return m.Store.GetDB().Model(Token).Create(Token).Error
}

func (m *MysqlToken) Delete(id int) error {
	return m.Store.GetDB().Delete(&auth.Token{}, "id = ?", id).Error
}

func (m *MysqlToken) Update(id int, ma map[string]interface{}) error {
	return m.Store.GetDB().Update(&auth.Token{}, ma).Error
}

func (m *MysqlToken) Query(where string, args []interface{}) (*auth.Token, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = auth.Token{}
	query_obj.Where = where
	query_obj.Args = args
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*auth.Token), nil
}

func (m *MysqlToken) Count(where string, args []interface{}) (int, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = auth.Token{}
	query_obj.Where = where
	query_obj.Args = args
	return m.Store.Count(query_obj)
}

func (m *MysqlToken) QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.Token, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = []auth.Token{}
	query_obj.Where = where
	query_obj.Args = args
	query_obj.Limit = limit
	query_obj.Offset = offset
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*[]auth.Token), nil
}
