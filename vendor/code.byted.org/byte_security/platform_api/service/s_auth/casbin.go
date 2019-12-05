// @Contact:    huaxinrui
// @Time:       2019/9/26 下午9:25

package s_auth

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/platform_api/dao/d_auth"
)

type CasbinServiceIF interface {
	Create(event *auth.CasbinRule) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*auth.CasbinRule, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.CasbinRule, error)
	Count(where string, args []interface{}) (int, error)
}

var CasbinServiceGlobal CasbinServiceIF

func init() {
	tmp := new(CasbinRuleService)
	tmp.dao = d_auth.CasbinDao
	CasbinServiceGlobal = tmp
}

type CasbinRuleService struct {
	dao d_auth.CasbinDaoIF
}

func InitEventService(i d_auth.CasbinDaoIF) *CasbinRuleService {
	a := new(CasbinRuleService)
	a.dao = i
	return a
}

func (m *CasbinRuleService) Create(app *auth.CasbinRule) error {
	return m.dao.Create(app)
}

func (m *CasbinRuleService) Delete(id int) error {
	return m.dao.Delete(id)
}

func (m *CasbinRuleService) Update(id int, ma map[string]interface{}) error {
	return m.dao.Update(id, ma)
}

func (m *CasbinRuleService) Query(where string, args []interface{}) (*auth.CasbinRule, error) {
	return m.dao.Query(where, args)
}

func (m *CasbinRuleService) Count(where string, args []interface{}) (int, error) {
	return m.dao.Count(where, args)
}

func (m *CasbinRuleService) QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.CasbinRule, error) {
	return m.dao.QueryAll(where, args, offset, limit)
}
