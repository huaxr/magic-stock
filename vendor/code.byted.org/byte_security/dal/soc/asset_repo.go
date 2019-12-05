// @Contact:    huaxinrui
// @Time:       2019/8/2 下午5:13

// 代码仓库
package soc

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type Repo struct {
	gorm.Model
	Name     string      `json:"name"`
	Owner    string      `json:"owner"`
	RepoName string      `json:"repo_name"`
	GroupId  string      `json:"group_id"`
	HttpUrl  string      `json:"http_url"`
	Desc     string      `json:"desc"`
	Lang     string      `json:"lang"`
	Extra    common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (Repo) New() interface{} {
	return &Repo{}
}

func (Repo) TableName() string {
	return "byte_security_asset_repo"
}

func (repo Repo) GetAssetKey() string {
	return "repo_name"
}

func (repo Repo) GetAssetValue() string {
	return repo.RepoName
}
