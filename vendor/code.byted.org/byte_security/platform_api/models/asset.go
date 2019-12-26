// @Contact:    huaxinrui
// @Time:       2019/8/28 下午12:12

package models

import (
	"time"

	dal_common "code.byted.org/byte_security/dal/common"
	"code.byted.org/byte_security/dal/soc"
)

type ResultDomainDetail struct {
	Icp     dal_common.JSON `json:"icp"`
	Whois   dal_common.JSON `json:"whois"`
	Cert    dal_common.JSON `json:"cert"`
	Dns     dal_common.JSON `json:"dns"`
	Cluster dal_common.JSON `json:"cluster"`
}

type ResultHostDetail struct {
	Created time.Time     `json:"created"`
	Updated time.Time     `json:"updated"`
	User    string        `json:"user"`
	Dpkg    string        `json:"dpkg"`
	Crontab string        `json:"crontab"`
	Eth     string        `json:"eth"`
	Ports   string        `json:"ports"`
	OuterIp []soc.NetWork `json:"outer_ip"`

	Id           int    `json:"id"`
	Ip           string `json:"ip"`
	HostName     string `json:"host_name"`
	Os           string `json:"os"`
	Kernel       string `json:"kernel"`
	HoneyPot     int    `json:"honeypot"`
	HoneyPotPort int    `json:"honeypot_port"`
	State        int    `json:"state"`
	AgentVersion string `json:"agent_version"`
	Idc          string `json:"idc"`
}

type WebsiterDetail struct {
	Wappalyzer dal_common.JSON `json:"wappalyzer"`
	Url        string          `json:"url"`
	Owner      string          `json:"owner"`
}

type PsmDetail struct {
	Psm         string          `json:"psm"`
	Owner       string          `json:"owner"`
	Group       string          `json:"group"`
	ScmRepos    dal_common.JSON `json:"scm_repos"`
	Subscribers dal_common.JSON `json:"subscribers"`
	Extra       dal_common.JSON `json:"extra"`
	NodeID      int             `json:"node_id"`
	Path        string          `json:"path"`
	CreatedTime time.Time       `json:"created_time"`
	UpdatedTime time.Time       `json:"updated_time"`
}

type RepoDetail struct {
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	RepoName string `json:"repo_name"`
	GroupId  string `json:"group_id"`
	HttpUrl  string `json:"http_url"`
	Desc     string `json:"desc"`
	Lang     string `json:"lang"`
}

type HostList struct {
	Level        int    `json:"level"`
	RelateTicket int    `json:"related_ticket"`
	Ip           string `json:"ip"`
	AgentVersion string `json:"agent_version"`
	Owner        string `json:"owner"`
	State        string `json:"state"`
	Idc          string `json:"idc"`
}

type DomainList struct {
	Level        int    `json:"level"`
	RelateTicket int    `json:"related_ticket"`
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	Visible      string `json:"visible"`
	Operator     string `json:"operator"`
	Vendor       string `json:"vendor"`
	Usage        string `json:"usage"`
	AccType      string `json:"acc_type"` // inner api public rd
	Scope        string `json:"scope"`    // 可见范围
}

type RepoList struct {
	Level        int    `json:"level"`
	RelateTicket int    `json:"related_ticket"`
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	Lang         string `json:"lang"`
	Desc         string `json:"desc"`
	HttpUrl      string `json:"http_url"`
}

type AssetGroupAndOwner struct {
	GroupName string `json:"group_name"`
	Owner     string `json:"owner"`
}

type RelatedPsm struct {
	Psm         string    `json:"psm"`
	Owner       string    `json:"owner"`
	Created     time.Time `json:"created"`
	Level       int       `json:"level"`
	TicketCount int       `json:"ticket_count"`
	EventCount  int       `json:"event_count"`
	VulnCount   int       `json:"vuln_count"`
}

type RelatedDomain struct {
	Name        string `json:"name"` // 域名
	Owner       string `json:"owner"`
	AccType     string `json:"acc_type"` // inner api public rd
	Scope       string `json:"scope"`    // 可见范围
	Vendor      string `json:"vendor"`   // cdn
	Level       int    `json:"level"`
	TicketCount int    `json:"ticket_count"`
	EventCount  int    `json:"event_count"`
	VulnCount   int    `json:"vuln_count"`
}

type RelatedRepo struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Desc        string `json:"desc"`
	Lang        string `json:"lang"`
	Level       int    `json:"level"`
	TicketCount int    `json:"ticket_count"`
	EventCount  int    `json:"event_count"`
	VulnCount   int    `json:"vuln_count"`
}

type RelatedHost struct {
	Ip          string `json:"ip"`
	Owner       string `json:"owner"`
	Level       int    `json:"level"`
	TicketCount int    `json:"ticket_count"`
	EventCount  int    `json:"event_count"`
	VulnCount   int    `json:"vuln_count"`
}

type RelatedWebsite struct {
	Wappalyzer dal_common.JSON `sql:"type:json" json:"wappalyzer"` // web 指纹
	Url        string
	Owner      string
	GroupId    string
}

type RelatedResult struct {
	Psm     []RelatedPsm     `json:"psm"`
	Repo    []RelatedRepo    `json:"repo"`
	Domain  []RelatedDomain  `json:"domain"`
	Host    []RelatedHost    `json:"host"`
	Website []RelatedWebsite `json:"website"`
	Total   int              `json:"total"`
}

type DomainNameTimeResult struct {
	Domain    string `json:"domain"`
	CreatedAt string `json:"created_at"`
}
