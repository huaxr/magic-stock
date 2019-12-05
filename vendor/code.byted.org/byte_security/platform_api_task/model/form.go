package model

import (
	"errors"
)

type TaskRequest struct {
	TaskID        int              `json:"task_id" form:"task_id"`
	TaskIds       []int            `json:"task_ids" form:"task_ids"`
	TaskName      string           `json:"task_name"  form:"task_name"`
	DataSourceId  int              `json:"data_source_id"`
	ServiceType   int              `json:"service_type" form:"service_type"`
	Status        int              `json:"status"`
	AllActiveType bool             `json:"all_active_type" form:"all_active_type" `
	TaskActive    bool             `json:"task_active" form:"task_active"`
	TaskType      int              `json:"task_type"`
	ContentDetail []byte           `json:"content_detail"`
	TaskFilter    *KeyTaskFilter   `json:"task_filter"`
	MaxData       int              `json:"max_data"`
	ExpireDate    int64            `json:"expire_date"`
	TaskDetail    map[string][]int `json:"task_detail"`
	TaskTrigger   *TaskTrigger     `json:"task_trigger"`
	Page          int              `json:"page" form:"page"`
	PageSize      int              `json:"page_size" form:"page_size"`
}

type TaskResponse struct {
	TaskData []*TaskData `json:"task_data"`
	Total    int         `json:"total"`
}

type TaskData struct {
	TaskID          int              `json:"task_id"`
	TaskName        string           `json:"task_name"`
	DataSourceId    int              `json:"data_source_id"`
	ServiceType     int              `json:"service_type" form:"service_type"`
	Status          int              `json:"status"`
	TaskActive      bool             `json:"task_active" form:"task_active"`
	TaskType        int              `json:"task_type"`
	ContentDetail   []byte           `json:"content_detail"`
	TaskFilter      *KeyTaskFilter   `json:"task_filter"`
	MaxData         int              `json:"max_data"`
	ExpireDate      int64            `json:"expire_date"`
	TaskDetail      map[string][]int `json:"task_detail"`
	TaskTrigger     *TaskTrigger     `json:"task_trigger"`
	TaskVersion     int              `json:"task_version"`
	TaskLastChildId int              `json:"task_last_child_id"`
}


//新建表单表单验证
func (t *TaskRequest) NewCheck() error {
	//非空逻辑判断
	if !t.parametersNotNull(t.TaskName, t.TaskType, t.ServiceType, t.TaskDetail) {
		return errors.New("参数必填项存在空值")
	}
	return nil
}

//开启新任务表单验证
func (t *TaskRequest) StartCheck() error {
	//非空逻辑判断
	if !t.parametersNotNull(t.TaskID) {
		return errors.New("任务ID值为空")
	}
	return nil
}

//更新任务表单验证
func (t *TaskRequest) UpdateCheck() error {
	//非空逻辑判断
	if !t.parametersNotNull(t.TaskID) {
		return errors.New("参数值存在为空")
	}
	return nil
}

//删除任务表单验证
func (t *TaskRequest) DeleteCheck() error {
	//非空逻辑判断
	if !t.parametersNotNull(t.TaskIds) {
		return errors.New("参数值存在为空")
	}
	return nil
}

//查询所有当前类型表单验证
func (t *TaskRequest) ListAllCheck() error {
	//非空逻辑判断
	if !t.parametersNotNull(t.ServiceType) {
		return errors.New("参数值存在为空")
	}
	return nil
}

//查询最后一条运行完成的任务
func (t *TaskRequest) LastFinishCheck() error {
	//非空逻辑判断
	if !t.parametersNotNull(t.TaskID) {
		return errors.New("参数值存在为空")
	}
	return nil
}

//批量启停表单验证
func (t *TaskRequest) ActiveChangeBatchCheck() error {
	//非空逻辑判断
	if !t.parametersNotNull(t.TaskIds) {
		return errors.New("参数值存在为空")
	}
	return nil
}

func (t *TaskRequest) parametersNotNull(args ...interface{}) bool {
	for i := range args {
		switch args[i].(type) {
		case string:
			if args[i] == "" {
				return false
			}
		case int:
			if args[i] == 0 {
				return false
			}
		case map[string][]int:
			if args[i] == nil {
				return false
			}
		case []byte:
			if args[i] == nil {
				return false
			}
		case []int:
			if len(args[i].([]int)) == 0 {
				return false
			}

		}
	}
	return true
}

type BasicResponse struct {
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Data      interface{} `json:"data"`
}

func (b *BasicResponse) Error(errMsg string) *BasicResponse {
	b.ErrorCode = 1
	b.ErrorMsg = errMsg
	b.Data = nil
	return b
}

func (b *BasicResponse) Success(data interface{}) *BasicResponse {
	b.ErrorCode = 0
	b.ErrorMsg = ""
	b.Data = data
	return b
}

type EngineResponseForm struct {
	Status int    `json:"status"`
	Data   []byte `json:"data"`
	ErrMsg string `json:"err_msg"`
}
