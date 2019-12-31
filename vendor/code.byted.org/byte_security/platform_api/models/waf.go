package models

import (
	"code.byted.org/byte_security/dal/common"
	"code.byted.org/byte_security/dal/policy"
)

/*============================线下策略==================================*/
type StrategyCreate struct {
	ProductId    int    `json:"product_id"`
	StrategyName string `json:"strategy_name"`
	//StrategyVersion string              `json:"strategy_version"`
	Priority int    `json:"priority"`
	Desc     string `json:"desc"`
	//Rules           []map[string]string `json:"rules"`
}

type StrategyEdit struct {
	StrategyId int    `json:"strategy_id"`
	Priority   int    `json:"priority"`
	Desc       string `json:"desc"`
	Filter     string `json:"filter"`
}

type StrategyAddVersion struct {
	ProductId         int    `json:"product_id"`
	StrategyOfflineID int    `json:"strategy_id"`
	VersionName       string `json:"version_name"`
	Filter            string `json:"filter"`
}

type StrategyVersionAddExpress struct {
	StrategyVersionOfflineId int `json:"strategy_version_offline_id"`
	//Rules             []struct {
	//	Name    string `json:"name"`
	//	Express string `json:"express"`
	//	Status  int    `json:"status"`
	//} `json:"rules"`
	Expresses []struct {
		ExpressId int `json:"express_id"`
		Status    int `json:"status"`
		Priority  int `json:"priority"`
	} `json:"rules"`
}

type StrategyVersionExpressEnabled struct {
	StrategyOfflineVersionID int `json:"strategy_version_offline_id"`
}

type StrategyVersionDelExpress struct {
	StrategyVersionOfflineId int   `json:"version_id"`
	ExpressId                []int `json:"ids"`
}

type StrategyVersionEditExpress struct {
	StrategyVersionOfflineId int `json:"version_id"`
	ExpressId                int `json:"id"`
	Status                   int `json:"status"`
	Priority                 int `json:"priority"`
}

type StrategyDelete struct {
	StrategyOfflineID []int `json:"ids"`
}

type StrategyVersionDelete struct {
	StrategyVersionID []int `json:"ids"`
}

/*============================线上策略==================================*/

type StrategyCopy struct {
	StrategyOfflineID        int    `json:"strategy_id"`
	StrategyOfflineVersionID int    `json:"strategy_version_id"`
	Name                     string `json:"name"`
}

type StrategyActive struct {
	StrategyOnlineID  int `json:"strategy_id"`
	StrategyVersionID int `json:"strategy_version_id"`
}

/*====================================规则表达式============================*/

type RuleExpressCreate struct {
	Name      string `json:"rule_name"`
	Desc      string `json:"rule_desc"`
	Express   string `json:"express"`
	ProductId int    `json:"product_id"`
}

type RuleExpressEdit struct {
	ProductId int    `json:"product_id"`
	ExpressId int    `json:"id"`
	Name      string `json:"name"`
	Express   string `json:"express"`
	Desc      string `json:"desc"`
}

type RuleExpressDelete struct {
	ExpressId []int `json:"ids"`
}

type StrategyVersionOption struct {
	Option            string `json:"option"`
	StrategyVersion   string `json:"strategy_version"`
	StrategyVersionId uint   `json:"strategy_version_id"`
	Rules             []struct {
		Name    string `json:"name"`
		Express string `json:"express"`
		Status  int    `json:"status"`
	} `json:"rules"`
}

type StrategyTask struct {
	Option            string `json:"option"`
	TaskId            uint   `json:"task_id"`
	StrategyVersionId uint   `json:"strategy_version_id"`
	TaskName          string `json:"task_name"`
	Count             int    `json:"count"`
}

type ProductCreate struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Key    string `json:"key"`
	Status int    `json:"status"`

	Code     string `json:"code"`      // 代号
	Doc      string `json:"doc"`       // 文档
	Tag      string `json:"tag"`       //（数据安全/代码安全/）
	Type     string `json:"type"`      // 能力 0决策，1基线，2第三方扫描，3第三方监控，4第三方加固
	FieldsId int    `json:"fields_id"` // 关联字段集（跟分析有关的能力）
}

type ProductStatusChange struct {
	ProductID string `json:"product_id"`
	Status    int    `json:"status"`
}

type GroupCreate struct {
	ProductId int                 `json:"product_id"`
	GroupKeys []map[string]string `json:"group_keys"` // [{"key":"psm", "value":"s.s.s"}, ... ]
	Desc      string              `json:"desc"`
	Locations []int               `json:"locations"`
	Strategys []int               `json:"strategys"`
	State     int                 `json:"state"`
}

type GroupDelete struct {
	GroupIds []int `json:"group_ids"`
}

type GroupUpdate struct {
	GroupId   int   `json:"group_id"`
	Locations []int `json:"locations"` // 修改区域（多选）
	State     int   `json:"state"`
	Strategys []int `json:"strategys"` // 修改策略
}

/*=============================变量 =========================*/
type VariateCreate struct {
	ProductID int      `json:"product_id"`
	Name      string   `json:"name"`
	Desc      string   `json:"desc"`
	SortType  int      `json:"sort_type"`
	IsStore   int      `json:"is_store"`
	Objs      []string `json:"objs"`
}

type VariateEdit struct {
	ProductId int      `json:"product_id"`
	VariateId int      `json:"id"`
	Desc      string   `json:"desc"`
	SortType  int      `json:"sort_type"`
	Objs      []string `json:"objs"`
}

type VariateDelete struct {
	ProductID int   `json:"product_id"`
	VariateId []int `json:"ids"`
}

/*======================字段====================*/
type FieldCreate struct {
	ProductID    int `json:"product_id"`
	NameWithType []struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		MappingField string `json:"mapping_field"`
		MappingName  string `json:"mapping_name"`
	} `json:"name_type"`
}

type FieldDelete struct {
	FieldID []int `json:"ids"`
}

/*=====================地区=====================*/

type LocationCreate struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Desc  string `json:"desc"`
}

// result
type ResultGroupList struct {
	ID        int                             `json:"Id"`
	ProductID int                             `json:"ProductID"`
	Desc      string                          `json:"Desc"`
	State     int                             `json:"State"`
	GroupKeys []policy.GroupKey               `json:"GroupKeys"`
	Strategys []policy.GroupAndStrategyOnline `json:"Strategys"`
	Locations []policy.GroupAndLocation       `json:"Locations"`
}

type ResultStrategyOffline struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	VersionCount int    `json:"count"`
	Filter       string `json:"filter"`
	Priority     int    `json:"priority"`
}

type ResultStrategyVersion struct {
	Id          int    `json:"id"`
	VersionName string `json:"name"`
	RuleCount   int    `json:"rule_count"`
	ProductId   int    `json:"product_id"`
	Filter      string `json:"filter"`
}

type ResultStrategyVersionDetail struct {
	ID       int    `json:"id"`
	Status   int    `json:"status"`
	Priority int    `json:"priority"`
	Express  string `json:"express"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
}

type ResultStrategyOnLine struct {
	Id           int    `json:"Id"`
	Name         string `json:"Name"`
	Desc         string `json:"Desc"`
	VersionCount int    `json:"VersionCount"`
	GroupCount   int    `json:"GroupCount"`
	VersionName  string `json:"VersionName"` // 生肖版本名称
	Filter       string `json:"filter"`
}

type ResultStrategyOnLineVersion struct {
	Id          int    `json:"id"`
	VersionName string `json:"name"`
	RuleCount   int    `json:"rule_count"`
	Active      bool   `json:"active"`
	ProductId   int    `json:"product_id"`
	Filter      string `json:"filter"`
}

type ResultVariate struct {
	Id       int         `json:"ID"`
	Name     string      `json:"Name"`
	Desc     string      `json:"Desc"`
	Type     string      `json:"Type"`
	SortType int         `json:"SortType"`
	Value    interface{} `json:"Value"`
}

type TaskCreate struct {
	ProductId int    //必选项，执行此产品的产品名称
	TaskName  string //必填项，定义下发任务的名称

	Owner          string //必填项，定义下发任务的人员名称，与lark与SSO名称一致
	DataSource     int    //必填，定义数据输入用户源头，0为kafka,后续可扩展
	TaskType       int
	TaskStopType   int         `json:"task_stop_type"`                   //必填项，0为消费数量级停止,1时间条件停止
	TaskStopDetail int         `json:"task_stop_detail"`                 //必选项，0为消费数量多少条后停止，1为多少秒后停止
	TaskDetail     common.JSON `sql:"type:json" json:"object,omitempty"` //必填，请序列化PolicyTaskDetail结构体为json后存入
	TaskFilter     common.JSON `sql:"type:json" json:"object,omitempty"`
	TaskStatus     int         //必填，0为达到条件后停止，1为正在运行，2为用户手动停止
}

type CheckVariate struct {
	VarID    int    `json:"id"`
	TestData string `json:"data"`
}

type CheckExpress struct {
	ExpressID int    `json:"id"`
	TestData  string `json:"data"`
}

type CheckStrategyOffline struct {
	StrategyOfflineID int    `json:"id"`
	TestData          string `json:"data"`
}

// 能力
type DataSourceDetail struct {
	ID           int               `json:"id" form:"id"`
	Name         string            `json:"name" form:"name"`
	Code         string            `json:"code" form:"name"`
	ConnInfo     map[string]string `json:"conn_info"form:"conn_info"`
	SourceType   int               `json:"source_type" form:"source_type"`
	IsEvent      bool              `json:"is_event" form:"is_event"`
	RwType       int               `json:"rw_type" form:"rw_type"`
	FieldsId     int               `json:"fields_id" form:"fields_id" `
	FieldsDetail FieldsResponse    `json:"fields_detail" form:"fields_detail"`
}
type FieldsResponse struct {
	ID            int                   `json:"id" form:"id"`
	Name          string                `json:"name" `
	Code          string                `json:"code" `
	Desc          string                `json:"desc"`
	FieldDetail   []*FieldDetail        `json:"field_detail" form:"field_detail"`
	DataSourceUse []*DataSourceAbstract `json:"data_source_use"`
}
type DataSourceAbstract struct {
	ID   int    `json:"id"`
	Name string `json:"name" `
	Code string `json:"code" `
}
type FieldDetail struct {
	ID        int             `json:"id" form:"id"`
	Name      string          `json:"name" form:"name"`
	Type      string          `json:"type" form:"type"`
	Desc      string          `json:"desc"`
	SignField []*SignResponse `json:"sign_field"`
}
type SignResponse struct {
	ID       int    `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Code     string `json:"code" form:"name"`
	IsUnique bool   `json:"is_unique" form:"name"`
	Type     int    `json:"type" form:"name"`
	Scope    []int  `json:"scope" form:"scope"`
}
type DataSourceDetailResponse struct {
	Data      DataSourceDetail `json:"data"`
	ErrorCode int              `json:"error_code"`
	ErrorMsg  string           `json:"error_msg"`
}
type FieldsDetailResponse struct {
	Data      map[string]string `json:"data"`
	ErrorCode int               `json:"error_code"`
	ErrorMsg  string            `json:"error_msg"`
}
