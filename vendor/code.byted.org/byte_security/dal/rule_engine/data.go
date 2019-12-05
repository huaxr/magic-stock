// @Contact:    huaxinrui
// @Time:       2019/7/30 下午2:53

package rule_engine

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type TaskDataSource struct {
	gorm.Model
	Data   common.JSON `sql:"type:json" json:"data,omitempty"`
	TaskId int
}

func (TaskDataSource) TableName() string {
	return "byte_security_ruleengine_taskdata"

}
