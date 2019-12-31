// @Contact:    huaxinrui
// @Time:       2019/7/9 上午11:25

package event

import (
	"time"

	"code.byted.org/byte_security/dal/common"
	"code.byted.org/gopkg/gorm"
)

type HResults struct {
	gorm.Model
	Asset     string `gorm:"index"` // 资产 如 ip， domain
	Uuid      string `gorm:"index"`
	Name      string // rule_str 改成name
	Type      string // 事件类型 HIDS VULN
	Level     int    // 事件等级
	Key       int
	TimeStamp time.Time
	//Detail common.JSON   `sql:"type:json" json:"object,omitempty"`
	RawLog  common.JSON `sql:"type:json" json:"raw_log,omitempty"` // TODO delete this field and query from es by uid
	RuleMap common.JSON `sql:"type:json" json:"rule_map,omitempty"`
	Docs    string      `json:"docs"` // 对应的知识库
	//Users  common.JSON   `sql:"type:json" json:"object,omitempty"`
	Users     string
	State     string `gorm:"index"` // pending ignore, done, distort, hangup（误报）
	Handler   string // 处理人
	TicketId  int    `gorm:"index"` // 对应工单
	Psm       string `gorm:"index"` // 对应psm
	Group     string // 事件所属组
	AssetType string `gorm:"index"` // 资产类型  host domain product ...
}

func (HResults) TableName() string {
	return "byte_security_event_result"
}
