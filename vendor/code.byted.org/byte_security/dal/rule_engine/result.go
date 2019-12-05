// @Contact:    huaxinrui
// @Time:       2019/7/30 下午2:56

package rule_engine

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type TaskDataResult struct {
	gorm.Model
	TaskId int
	DataId int
	Result common.JSON `sql:"type:json" json:"result"`
	Error  string      //错误信息字段
}

func (TaskDataResult) TableName() string {
	return "byte_security_ruleengine_taskresult"
}
