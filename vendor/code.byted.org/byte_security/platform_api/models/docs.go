package models

import "time"

// 知识库文章
type ArticleDetail struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`     // 作者
	AbilityID int       `json:"ability_id"` // 能力
	TagID     int       `json:"tag_id"`     // 标签
	CatalogID int       `json:"catalog_id"` // 目录
	Contactor string    `json:"contactor"`  // 联系人
	Title     string    `json:"title"`      // 标题
	Desc      string    `json:"desc"`       // 简介
	Theory    string    `json:"theory"`     // 原理
	Testing   string    `json:"testing"`    // 测试方法
	Repair    string    `json:"repair"`     // 修复方法
	Link      string    `json:"link"`       // 相关链接
	PostedAt  time.Time `json:"posted_at"`  // 发布时间
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

type ArticleList struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`     // 作者
	AbilityID int       `json:"ability_id"` // 能力
	TagID     int       `json:"tag_id"`     // 标签
	CatalogID int       `json:"catalog_id"` // 目录
	Title     string    `json:"title"`      // 标题
	Desc      string    `json:"desc"`       // 简介
	PostedAt  time.Time `json:"posted_at"`  // 发布时间
}

type Catalog struct {
	Name string `json:"name"`
	PID  int    `json:"pid"` // 父目录 ID
}
