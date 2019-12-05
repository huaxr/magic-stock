// @Contact:    huaxinrui
// @Time:       2019/8/2 下午4:33

package soc

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

// 域名和psm多对多
// psm 和 code一对一
// website 和 domain 多对一

type Website struct {
	gorm.Model
	Domain   Domain `gorm:"ForeignKey:DomainId"`
	DomainId int    // 属于哪个domain
	//Scm common.JSON `sql:"type:json" json:"object,omitempty"`  // 第三方库列表
	Wappalyzer common.JSON `sql:"type:json" json:"wappalyzer"` // web 指纹
	Url        string
	Owner      string
	GroupId    string
}

func (Website) TableName() string {
	return "byte_security_asset_website"
}

func (Website) New() interface{} {
	return &Website{}
}

func (web Website) GetAssetKey() string {
	return "url"
}

func (web Website) GetAssetValue() string {
	return web.Url
}
