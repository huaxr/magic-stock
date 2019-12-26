package docs

import (
	"time"

	"code.byted.org/gopkg/gorm"
)

// 知识库内容
type Article struct {
	gorm.Model

	Author    string    `json:"author"`     // 作者
	AbilityID int       `json:"ability_id"` // 能力
	TagID     int       `json:"tag_id"`     // 标签
	Catalog   Catalog   `gorm:"ForeignKey:CatalogID"  json:"catalog"`
	CatalogID int       `json:"catalog_id"`               // 目录
	Contactor string    `json:"contactor"`                // 联系人
	Title     string    `json:"title"`                    // 标题
	Desc      string    `json:"desc"`                     // 简介
	Theory    string    `gorm:"type:text" json:"theory"`  // 原理
	Testing   string    `gorm:"type:text" json:"testing"` // 测试方法
	Repair    string    `gorm:"type:text" json:"repair"`  // 修复方法
	Link      string    `json:"link"`                     // 相关链接
	PostedAt  time.Time `json:"posted_at"`                // 发布时间
}

func (Article) TableName() string {
	return "byte_security_docs_article"
}
