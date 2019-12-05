// @Contact:    huaxinrui
// @Time:       2019/7/30 下午2:46

package rule_engine

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type Task struct {
	gorm.Model
	ProductId      int    //必选项，执行此产品的产品名称
	TaskName       string //必填项，定义下发任务的名称
	TaskID         string `gorm:"index"` //十六位唯一ID
	Owner          string //必填项，定义下发任务的人员名称，与lark与SSO名称一致
	DataSource     int    //必填，定义数据输入用户源头，0为kafka,后续可扩展
	TaskType       int
	TaskStopType   int         `json:"task_stop_type"`                        //必填项，0为消费数量级停止,1时间条件停止
	TaskStopDetail int         `json:"task_stop_detail"`                      //必选项，0为消费数量多少条后停止，1为多少秒后停止
	TaskDetail     common.JSON `sql:"type:json" json:"task_detail,omitempty"` //必填，请序列化PolicyTaskDetail结构体为json后存入
	TaskFilter     common.JSON `sql:"type:json" json:"task_filter,omitempty"`
	TaskStatus     int         //必填，0为达到条件后停止，1为正在运行，2为用户手动停止
}

func (Task) TableName() string {
	return "byte_security_ruleengine_task"
}

func (t *Task) BeforeSave() (err error) {
	if t.TaskID == "" {
		t.TaskID = "123"
	}
	return nil
}
