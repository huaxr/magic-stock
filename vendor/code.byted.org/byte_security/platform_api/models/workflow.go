package models

import (
	"code.byted.org/byte_security/dal/common"
	"time"
)

type CreateTicketRequestAuto struct {
	Title  string                 `json:"title"`
	Desc   string                 `json:"desc"`
	Type   string                 `json:"type"`
	Users  string                 `json:"users"`
	Risk   int                    `json:"risk"`
	Detail map[string]interface{} `json:"detail"`
}

type CreateTicketRequestHand struct {
	Detail map[string]interface{} `json:"detail"`
	FlowId int                    `json:"flow_id"`
	Title  string                 `json:"title"`
	Desc   string                 `json:"desc"`
}

type CreateTicketRequestHIDS struct {
	Token  string   `json:"token"`
	Title  string   `json:"title"`
	Desc   string   `json:"desc"`
	Type   string   `json:"type"`
	Users  []string `json:"users"` // 当为hids时候， 从host找到tag的负责人列表
	Risk   int      `json:"risk"`  // 风险值  1,2,3,4  低位~严重
	Detail struct {
		Type      string    `json:"type"`
		TimeStamp time.Time `json:"timestamp"`
		Detail    struct {
			RowLog  map[string]interface{} `json:"raw_log"`
			Host    string                 `json:"host"`
			LogType string                 `json:"log_type"`
			RuleMap map[string]string      `json:"rule_map"`
			GroupId int                    `json:"group_id"`
			Risk    int                    `json:"risk"` // 风险值  1,2,3,4  低位~严重
		} `json:"detail"`
	} // 工单详情
}

type CancelTicketRequest struct {
	TicketId int `json:"ticket_id"`
}

type GetTicketRequest struct {
	TicketDetailId int `json:"ticket_id"`
}

type CommentTicketRequest struct {
	TicketId int    `json:"ticket_id"`
	Content  string `json:"content"`
}

type NodeOptionRequest struct {
	NodeId   int    `json:"node_id"`
	LastNode bool   `json:"last_node"`
	Option   string `json:"option"`
	UserName string `json:"user_name"` // 当 option 为 change_user时候
}

type AckOptionRequest struct {
	GroupId      int           `json:"group_id"`
	ExpireMinute time.Duration `json:"expire"`
	Type         string        `json:"type"` // 类型 e.g.  HIDS, WAF
}

type AddSub struct {
	TicketId int    `json:"ticket_id"`
	UserName string `json:"user_name"`
}

// result
type NodeTicket struct {
	Id       int    `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Priority int    `json:"priority"`
	State    string `json:"state"`
	//UserId int `json:"user_id"`
	//UserName string `json:"user_name"`
	//Users []map[string]int  // [{"huaxinrui":1, ...}]
	Users []string `json:"users"`
}

type ResultTicket struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  string `json:"type"`
	//Level int `json:"level"`
	UserId       int          `json:"user_id"` // 创建者
	UserName     string       `json:"user_name"`
	State        string       `json:"state"`
	TicketSource string       `json:"ticket_source"`
	Detail       interface{}  `json:"detail"`
	Extra        common.JSON  `json:"extra"`
	NodeInfo     []NodeTicket `json:"node_info"` // 节点详情
	CreatedAt    time.Time    `json:"created_at"`
	ProductName  string       `json:"product_name"`
	ProductOwner string       `json:"product_owner"`
}

type ResultMyTicket struct {
	Id           int                    `json:"id"`
	CreatedAt    time.Time              `json:"created_at"`
	Title        string                 `json:"title"`
	Desc         string                 `json:"desc"`
	Type         string                 `json:"type"`
	UserName     string                 `json:"user_name"`
	UserId       int                    `json:"user_id"`
	State        string                 `json:"state"` //('all') ('pending', '处理中'),  ('complete', '完成'),  ('cancel', '取消'), # 发起人可以取消 ('rejected', '拒绝')
	Detail       common.JSON            `json:"detail"`
	Extra        common.JSON            `json:"extra"`
	TicketSource string                 `json:"ticket_source"`
	Process      map[string]interface{} `json:"process"` // 处理进度
}

type AssetTicketResult struct {
	Id        int       `json:"id"`
	Asset     string    `json:"asset"` // 关联资产
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	State     string    `json:"state"`
}

type TicketProcessRes struct {
	State    string `json:"state"`
	Priority int    `json:"priority"`
	UserName string `json:"user_name"`
}

type RiskResult struct {
	Asset       string                 `json:"asset"`
	AssetInfo   map[string]interface{} `json:"asset_info"`
	Level       int                    `json:"level"`
	TicketCount int                    `json:"ticket_count"`
	EventCount  int                    `json:"event_count"`
	VulnCount int `json:"vuln_count"`
}
