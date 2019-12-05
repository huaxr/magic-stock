// @Contact:    huaxinrui
// @Time:       2019/8/8 上午10:28

package models

import (
	"time"

	dal_common "code.byted.org/byte_security/dal/common"
)

type EventOption struct {
	Title  string `json:"title"`
	Ids    []int  `json:"ids"`
	Option string `json:"option"`
	Token  string `json:"token"`
	User   string `json:"user"` // lark bot option user
}

type LarkUserInfo struct {
	EmployeeID string `json:"employee_id"`
	OpenID     string `json:"open_id"`
	TenantKey  string `json:"tenant_key"`
	UserID     string `json:"user_id"`
}

type EventOptionFromLark struct {
	AdditionalParameter LarkUserInfo `json:"additional_parameter"`
	CustomerParameter   EventOption  `json:"customer_parameter"`
}

// result
type ResultEventDetail struct {
	//Detail  dal_common.JSON `json:"detail"`
	//RawLog  dal_common.JSON `json:"raw_log"`
	RawLog    *EventLog       `json:"raw_log"`
	RuleMap   dal_common.JSON `json:"rule_map"`
	Users     string          `json:"users"`
	AssetId   int             `json:"asset_id"`
	Uuid      string          `json:"uuid"`
	Name      string          `json:"name"`
	Asset     string          `json:"asset"`
	Type      string          `json:"type"`
	Level     int             `json:"level"`
	State     string          `json:"state"`
	TicketId  int             `json:"ticket_id"`
	Handler   string          `json:"handler"` // 处理人
	CreatedAt time.Time       `json:"created_at"`
	TimeStamp time.Time       `json:"timestamp"`
	AssetType string          `json:"asset_type"`
}

type ResultEventVulnDetail struct {
	Detail map[string]interface{} `json:"detail"`
	Type   string
	State  string

	Users     string    `json:"users"`
	AssetId   int       `json:"asset_id"`
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Asset     string    `json:"asset"`
	Level     int       `json:"level"`
	TicketId  int       `json:"ticket_id"`
	Handler   string    `json:"handler"` // 处理人
	CreatedAt time.Time `json:"created_at"`
	TimeStamp time.Time `json:"timestamp"`

	AssetType string `json:"asset_type"` // 标识时间详情中的资产类型  psm
}

type ResultEventList struct {
	Id       uint      `json:"id"`
	Created  time.Time `json:"created"`
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Uuid     string    `json:"uuid"` // 新增uuid
	Key      int       `json:"key"`
	Asset    string    `json:"asset"`
	State    string    `json:"state"`
	Level    int       `json:"level"`
	TicketId uint      `json:"ticket_id"`
	Handler  string    `json:"handler"`
}

// es logs
type EventLog struct {
	AssetDetail map[string]interface{} `json:"asset_detail"`
	RawLog      map[string]interface{} `json:"raw_log"`
	Type        string                 `json:"type"`
	Uuid        string                 `json:"uuid"`
	Users       string                 `json:"users"`
	GroupId     int                    `json:"group_id"`
	Risk        int                    `json:"risk"`
	Asset       string                 `json:"asset"`
	Id          int                    `json:"id"`
	RuleMap     map[string]string      `json:"rule_map"`
	Timestamp   time.Time              `json:"timestamp"`
}

type CreateEventRequestAuto struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	//Detail map[string]interface{} `json:"detail"`
	Details SubmitPostData `json:"details"`
}

type AssetEventResult struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	TimeStamp time.Time `json:"time_stamp"`
	State     string    `json:"state"`
	TicketId  int       `json:"ticket_id"`
}

type EventASSET struct {
	TicketLevel int `json:"ticket_level"`
	TicketCount int `json:"ticket_count"`
	EventCount  int `json:"event_count"`
	VulnCount   int `json:"vuln_count"`
}

type EventResponse struct {
	Asset       EventASSET  `json:"asset"`
	AssetDetail interface{} `json:"asset_detail"`
	AssetType   string      `json:"asset_type"`
	AssetValue  string      `json:"asset_value"`
}
