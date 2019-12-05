package models

import (
	"time"

	task "code.byted.org/byte_security/platform_api_task/model"
)

type BLIds struct {
	Ids []int `json:"ids"`
}

type BLTaskOpt struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	RuleIds   string `json:"rule_ids"`
	Enabled   string `json:"enabled"`
	ExeMethod string `json:"exe_method"`
	Schedule  string `json:"schedule"`
}

type EnabledBLTaskOpt struct {
	Ids     []int  `json:"ids"`
	Enabled string `json:"enabled"`
}

type BLTaskResult struct {
	ErrorCode int               `json:"error_code"`
	ErrorMsg  string            `json:"error_msg"`
	Data      task.TaskResponse `json:"data"`
}

type BLTaskDetail struct {
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Data      interface{} `json:"data"`
}

type ProductAsset struct {
	ProductName string `json:"product_name"`
	AssetTable  string `json:"asset_table"`
}

type BaseLineEvent struct {
	ID        int       `json:"id"`
	Level     int       `json:"level"`
	Key       int       `json:"key"`
	TimeStamp time.Time `json:"time_stamp"`
	Name      string    `json:"name"`
	TicketID  int       `json:"ticket_id"`
}

type BaseLineRule struct {
	Level    int    `json:"level"`
	RuleName string `json:"rule_name"`
}

type AssetCompliance struct {
	TaskName  string         `json:"task_name"`
	TimeStamp time.Time      `json:"time_stamp"`
	TaskID    uint           `json:"task_id"`
	ID        int            `json:"id"`
	TicketID  int            `json:"ticket_id"`
	Rules     []BaseLineRule `json:"rules"`
}

// because of sth., duplicate this struct
type Task struct {
	ID                     uint       `json:"id"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at"`
	TaskName               string     `gorm:"not null"`
	DataSourceID           int
	TaskType               int `gorm:"not null"`
	TaskFilter             string
	TaskDetail             string `gorm:"not null"`
	TaskServiceType        int    `gorm:"not null"`
	TaskActivate           bool   `gorm:"not null"`
	ConsumeMaxData         int
	HitMaxData             int
	ExpireDate             int64
	TriggerType            int
	TriggerFrequencyType   int
	TriggerDetail          string
	LastRunningChildTaskId int `gorm:"not null"`
	TaskVersion            int
}
