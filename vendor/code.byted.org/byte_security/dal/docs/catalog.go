package docs

import (
	"code.byted.org/gopkg/gorm"
)

// 知识库目录
type Catalog struct {
	gorm.Model

	Name string `json:"name"`
	PID  *int   `json:"pid"` // 父目录 ID
}

func (Catalog) TableName() string {
	return "byte_security_docs_catalog"
}
