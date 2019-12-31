// 资产漏洞报告

package soc

import (
	"time"

	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type Vulnerability struct {
	gorm.Model
	Source          string //漏洞来源
	Code            string
	Level           string //漏洞等级
	Submitter       string //提交者
	SubmitterLeader string
	Title           string `gorm:"size:255"`  //漏洞标题
	Detail          string `gorm:"type:text"` //漏洞详情

	ProductID  int //所属产品
	AppID      int
	BusinessId int `gorm:"index"` // 所属业务线

	TypeId         int
	Status         string //状态 pending ignore, done, process
	TicketId       int    //工单id
	EventId        int    //事件id
	SysWorkflowId  int
	ContactOwner   string //联系人
	EnclosedFile   string //附件
	Suggestion     string `gorm:"type:text"` //修复建议
	HandUp         bool   // 是否挂起 表示正在处理
	TestType       string
	Docs           string // 知识库id列表1,2,3
	KeepSecret     bool
	AssetTypeKey   string // 所属资产，如domain, psm 等
	AssetTypeValue string
	CreatedTime    time.Time // 创建时间
	LimitTime      time.Time
	Extra          common.JSON `sql:"type:json" json:"object,omitempty"`
}

func (Vulnerability) TableName() string {
	return "byte_security_asset_vuln"
}
