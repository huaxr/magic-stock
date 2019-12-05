// @Contact:    huaxinrui
// @Time:       2019/9/26 下午9:25

package s_auth

import (
	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/dal/soc"
	"code.byted.org/byte_security/platform_api/core"
	"code.byted.org/byte_security/platform_api/dao/d_auth"
	"log"
)

type GroupServiceIF interface {
	Create(event *auth.Group) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*auth.Group, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.Group, error)
	Count(where string, args []interface{}) (int, error)
	GetGroupByAsset(tb, asset string) string
}

var GroupServiceGlobal GroupServiceIF

func init() {
	tmp := new(GroupService)
	tmp.dao = d_auth.GroupDao
	GroupServiceGlobal = tmp
}

type GroupService struct {
	dao d_auth.GroupDaoIF
}

func InitGroupService(i d_auth.GroupDaoIF) *GroupService {
	a := new(GroupService)
	a.dao = i
	return a
}

func (m *GroupService) Create(app *auth.Group) error {
	return m.dao.Create(app)
}

func (m *GroupService) Delete(id int) error {
	return m.dao.Delete(id)
}

func (m *GroupService) Update(id int, ma map[string]interface{}) error {
	return m.dao.Update(id, ma)
}

func (m *GroupService) Query(where string, args []interface{}) (*auth.Group, error) {
	return m.dao.Query(where, args)
}

func (m *GroupService) QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.Group, error) {
	return m.dao.QueryAll(where, args, offset, limit)
}

func (m *GroupService) Count(where string, args []interface{}) (int, error) {
	return m.dao.Count(where, args)
}

func (m *GroupService) GetGroupByAsset(tb, asset string) string {
	var GROUP string
	var err error
	switch tb {
	case "byte_security_asset_host":
		var host soc.Host
		err = core.Backend.Store.GetDB().Model(&soc.Host{}).Where("ip = ?", asset).Find(&host).Error
		if err != nil {
			log.Println(1, err, asset)
			return ""
		}
		GROUP = host.GroupId

	case "byte_security_asset_domain":
		var domain soc.Domain
		err := core.Backend.Store.GetDB().Model(&soc.Domain{}).Where("name = ?", asset).Find(&domain).Error
		if err != nil {
			log.Println(1, err, asset)
			return ""
		}
		GROUP = domain.GroupId

	case "byte_security_asset_product":
		var product soc.Product
		err = core.Backend.Store.GetDB().Model(&soc.Product{}).Where("name = ?", asset).Find(&product).Error
		if err != nil {
			log.Println(err)
			return ""
		}
		GROUP = product.Group
	}
	return GROUP
}
