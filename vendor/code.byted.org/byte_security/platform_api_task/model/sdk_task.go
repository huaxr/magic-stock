package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type BaseTaskResponseModel struct {
	TaskID  int `json:"task_id"`
	Version int `json:"version"`
	//任务自动停止条件
	ExpireDate      int64 `json:"expire_date"` //任务过期时间
	TaskContentInfo []*BaseTaskContent
}

type SourceInfo struct {
	DataSource int               `json:"data_source"` //数据来源类型，0，kafka，1，用户输入
	Field      map[string]string `json:"field"`
}

type BaseTaskContent struct {
	SourceInfo *SourceInfo

	TaskStatus TASK_STATUS `json:"task_status"` //0，条件停止，1，正在运行，2手动终止，3，下发失败

	MaxData int `json:"max_data"` //最大消费额

	//任务过滤条件，满足过滤条件的数据才会进行相关决策
	TaskFilterType TASK_FILTER_TYPE `json:"task_filter_type"` //数据过滤类型
	TaskFilter     []byte           `json:"task_filter"`      //数据过滤条件

	//任务详情
	TaskType TASK_TYPE `json:"type"`   //0 变量 ，1规则，2策略，3变量集，4规则集
	Detail   []byte    `json:"detail"` //VarTaskDetail、RuleTaskDetail、PolicyTaskDetail
}

type KeyTaskFilter struct {
	GroupKeys  []string `json:"group_keys"`
	GroupValue []string `json:"group_value"`
	//resultIsStore	int`json:"result_is_store"` 预定义字段，过滤（策略和规则）存储的结果，0为存储全部结果，1为存储命中结果
}

type TaskTrigger struct {
	Type          int            `json:"type"` //0。小时级，1，天级，2，周级，3月级
	FrequencyType int            `json:"frequency_type"`
	TriggerDetail *TriggerDetail `json:"trigger_detail"`
}

type TriggerDetail struct {
	TriggerMinute int   `json:"trigger_minute"`
	TriggerHours  []int `json:"trigger_hours"`
	TriggerWeeks  []int `json:"trigger_weeks"`
	TriggerDays   []int `json:"trigger_days"`
}

func (taskTrigger *TaskTrigger) GetTaskTriggerCron() (string, error) {
	var specCron string
	var minute, hours, week, day string
	if taskTrigger.Type != 1 {
		err := errors.New("非周期性任务")
		return "", err
	}
	switch taskTrigger.FrequencyType {
	case 0:
		minute = strconv.Itoa(taskTrigger.TriggerDetail.TriggerMinute)
		hours = strings.Replace(strings.Trim(fmt.Sprint(taskTrigger.TriggerDetail.TriggerHours), "[]"), " ", ",", -1)
		week = "*"
		day = "*"
	case 1:
		minute = strconv.Itoa(taskTrigger.TriggerDetail.TriggerMinute)
		hours = strconv.Itoa(taskTrigger.TriggerDetail.TriggerHours[0])
		week = "*"
		day = "*"
	case 2:
		minute = strconv.Itoa(taskTrigger.TriggerDetail.TriggerMinute)
		hours = strconv.Itoa(taskTrigger.TriggerDetail.TriggerHours[0])
		week = strings.Replace(strings.Trim(fmt.Sprint(taskTrigger.TriggerDetail.TriggerWeeks), "[]"), " ", ",", -1)
		day = "*"
	case 3:
		minute = strconv.Itoa(taskTrigger.TriggerDetail.TriggerMinute)
		hours = strconv.Itoa(taskTrigger.TriggerDetail.TriggerHours[0])
		week = "*"
		day = strings.Replace(strings.Trim(fmt.Sprint(taskTrigger.TriggerDetail.TriggerDays), "[]"), " ", ",", -1)
	default:
		err := errors.New("不可解析的周期类型")
		return " ", err
	}
	specCron = "0" + " " + minute + " " + hours + " " + day + " * " + week
	return specCron, nil
}
