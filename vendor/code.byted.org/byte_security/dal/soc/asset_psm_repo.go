package soc

import (
	"code.byted.org/gopkg/gorm"
)

type PSMRepo struct {
	gorm.Model
	PsmId  string `json:"psm_id"`
	RepoId string `json:"repo_id"`
}

func (PSMRepo) TableName() string {
	return "byte_security_asset_psm_repo"
}
