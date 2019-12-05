package policy

import (
	"code.byted.org/gopkg/gorm"
	"time"
)


type Task struct {
	gorm.Model
	//StrategyVersion StrategyVersion `gorm:"ForeignKey:StrategyVersionId"`
	StrategyVersionId int
	Name string
	StartTime time.Time
	Count int
	Status string
	Deleted bool
}


func (Task) TableName() string {
	return "byte_security_policy_task"
}
