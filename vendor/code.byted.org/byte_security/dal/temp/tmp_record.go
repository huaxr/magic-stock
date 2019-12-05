// @Contact:    huaxinrui
// @Time:       2019/10/30 下午4:11

package temp

import (
	"code.byted.org/gopkg/gorm"
)

type TmpRecord struct {
	gorm.Model
	LinkId int  // 关联id
	User string // 访问人
	Success bool
}

func (TmpRecord) TableName() string {
	return "byte_security_temp_record"
}
