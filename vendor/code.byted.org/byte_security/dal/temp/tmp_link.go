// @Contact:    huaxinrui
// @Time:       2019/10/30 下午3:59

package temp

import (
	"code.byted.org/gopkg/gorm"
	"time"
)

type TmpLink struct {
	gorm.Model
	LinkType string // 临时链接类型
	RelatedId int  // 关联ID
	Url string   // 关联的url md5(id, type)
 	Created time.Time
	Owner string
	DisableTime time.Time // 失效时间
}

func (TmpLink) TableName() string {
	return "byte_security_temp_link"
}
