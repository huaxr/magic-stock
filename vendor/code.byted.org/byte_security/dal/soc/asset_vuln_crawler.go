// @Contact:    huaxinrui
// @Time:       2019/11/5 下午2:41

package soc

import "time"

type CrawlerSeebug struct {
	ID        uint `gorm:"primary_key"`
	Title string
	Serial string
	Time string
	Component string
	Cve string
	Type string
	AffectVersion string
	Level string
	UpTime time.Time
	SeebugUrl string
	AffectUrls string `sql:"type:text"`// 影响的资产
	Flag bool // AffectUrls 是否为空
}

func (CrawlerSeebug) TableName() string {
	return "byte_security_asset_vuln_crawler"
}
