// @Contact:    huaxinrui
// @Time:       2019/9/26 下午9:25

package s_auth

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/platform_api/dao/d_auth"
)

type TokenServiceIF interface {
	Create(event *auth.Token) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*auth.Token, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.Token, error)
	Count(where string, args []interface{}) (int, error)
}

var TokenServiceGlobal TokenServiceIF

func init() {
	tmp := new(TokenService)
	tmp.dao = d_auth.TokenDao
	TokenServiceGlobal = tmp
}

type TokenService struct {
	dao d_auth.TokenDaoIF
}

func InitTokenService(i d_auth.TokenDaoIF) *TokenService {
	a := new(TokenService)
	a.dao = i
	return a
}

func (m *TokenService) Create(app *auth.Token) error {
	return m.dao.Create(app)
}

func (m *TokenService) Delete(id int) error {
	return m.dao.Delete(id)
}

func (m *TokenService) Update(id int, ma map[string]interface{}) error {
	return m.dao.Update(id, ma)
}

func (m *TokenService) Query(where string, args []interface{}) (*auth.Token, error) {
	return m.dao.Query(where, args)
}

func (m *TokenService) QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.Token, error) {
	return m.dao.QueryAll(where, args, offset, limit)
}
func (m *TokenService) Count(where string, args []interface{}) (int, error) {
	return m.dao.Count(where, args)
}
