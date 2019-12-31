package models

import "time"

// 新建NAT规则
type NatRuleCreate struct {
	BuzName     string `json:"buz_name"`
	DeployType  int    `json:"deploy_type"`
	SrcIP       string `json:"src_ip"`
	SrcPSM      string `json:"src_psm"`
	ServiceID   int    `json:"service_id"`
	ServiceName string `json:"service_name"`
	DestIP      string `json:"dest_ip"`
	DestDomain  string `json:"dest_domain"`
	AccessType  int    `json:"access_type"`
	Comments    string `json:"comments"`
	RuleClass   int    `json:"rule_class"`
	RuleType    int    `json:"rule_type"`
	Creator     string `json:"creator"`
	IsRoot      int    `json:"is_root"`
	FlowID      int    `json:"flow_id"`
}

// 更新NAT规则
type NatRuleUpdate struct {
	ID          int    `json:"id"`
	BuzName     string `json:"buz_name"`
	ServiceID   int    `json:"service_id"`
	ServiceName string `json:"service_name"`
	DestIP      string `json:"dest_ip"`
	DestDomain  string `json:"dest_domain"`
	AccessType  int    `json:"access_type"`
	Comments    string `json:"comments"`
}

type NatRuleSearch struct {
	BuzName    string    `json:"buz_name"`
	RuleType   int       `json:"rule_type"`
	Creator    string    `json:"creator"`
	IsActive   int       `json:"is_active"`
	Comments   string    `json:"comments"`
	CreateTime time.Time `json:"create_time"`
}

type NatRuleOpt struct {
	ID int    `json:"id"`
	OP string `json:"op"`
}

type NatRuleAuth struct {
	OPT int `json:"opt"`
}
