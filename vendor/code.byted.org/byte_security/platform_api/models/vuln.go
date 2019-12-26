// @Contact:    huaxinrui
// @Time:       2019/9/26 下午3:47

package models

import "time"

type VulnResList struct {
	Id           int        `json:"id"`
	Level        string     `json:"level"`
	Submitter    string     `json:"submitter"`
	Title        string     `json:"title"`
	HandUp       bool       `json:"hand_up"`
	Status       string     `json:"status"`
	ContactOwner string     `json:"contact_owner"`
	Name         string     `json:"type_name"`
	SubName      string     `json:"type_sub_name"`
	Ticket_id    int        `json:"ticket_id"`
	EventID      int        `json:"event_id"`
	CreatedTime  time.Time  `json:"created"`
	LimitTime    *time.Time `json:"limited"`

	ProductName    string `json:"product_name"`
	AssetTypeKey   string `json:"asset_type_key"`
	AssetTypeValue string `json:"asset_type_value"`

	BusinessId int `json:"business_id"`
}

type VulnResDetail struct {
	Id           int        `json:"id"`
	ProductName  string     `json:"product_name"`
	Code         string     `json:"code"`
	Suggestion   string     `json:"suggestion"`
	Detail       string     `json:"detail"`
	Level        string     `json:"level"`
	Submitter    string     `json:"submitter"`
	Title        string     `json:"title"`
	HandUp       bool       `json:"hand_up"`
	Status       string     `json:"status"`
	ContactOwner string     `json:"contact_owner"`
	Name         string     `json:"type_name"`
	SubName      string     `json:"type_sub_name"`
	EventId      int        `json:"event_id"`
	Ticket_id    int        `json:"ticket_id"`
	CreatedTime  time.Time  `json:"created"`
	LimitTime    *time.Time `json:"limit_time"`

	Psm            string `json:"psm"`
	TestType       string `json:"test_type"`
	From           string `json:"from"`
	EnclosedFile   string `json:"enclosed_file"`
	AssetTypeKey   string `json:"asset_type_key"`
	AssetTypeValue string `json:"asset_type_value"`
	AppId          int    `json:"app_id"`
}

type VTName struct {
	Name string `json:"name"`
}

type VTNameSub struct {
	SubName string `json:"sub_name"`
}

type SubmitPostData struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	TestType     string `json:"test_type"`
	AppId        int    `json:"app_id"`
	VULNType     string `json:"vuln_type"`
	VULNSubType  string `json:"vuln_sub_type"`
	Level        string `json:"level"` // 对应等级
	From         string `json:"from"`
	Submitter    string `json:"submitter"`
	RelatedOwner string `json:"related_owner"`
	Contact      string `json:"contact"`
	KeepSecret   bool   `json:"keep_secret"`
	Detail       string `json:"detail"`
	Suggestion   string `json:"suggestion"`
	EnclosedFile string `json:"enclosed_file"`

	AssetTypeKey   string `json:"asset_type_key"`
	AssetTypeValue string `json:"asset_type_value"`
}

type UpdatePostData struct {
	TestType     string `json:"test_type"`
	AppId        int    `json:"app_id"`
	VULNType     string `json:"vuln_type"`
	VULNSubType  string `json:"vuln_sub_type"`
	Level        string `json:"level"` // 对应等级
	From         string `json:"from"`
	RelatedOwner string `json:"related_owner"`
	Contact      string `json:"contact"`
	KeepSecret   bool   `json:"keep_secret"`
	Detail       string `json:"detail"`
	Suggestion   string `json:"suggestion"`
	EnclosedFile string `json:"enclosed_file"`

	AssetTypeKey   string `json:"asset_type_key"`
	AssetTypeValue string `json:"asset_type_value"`
}

type SRCData struct {
	Id                int    `json:"id"`
	App               string `json:"app"`
	Title             string `json:"title"`
	Detail            string `json:"detail"`
	Type              string `json:"type"`
	Platform          string `json:"platform"`
	EnclosedFile      string `json:"enclosed_file"`
	Level             string `json:"level"`
	NickName          string `json:"nick_name"`
	Contact           string `json:"contact"`
	RelatedAssetType  string `json:"related_asset_type"`
	RelatedAssetValue string `json:"related_asset_value"`
}

type SRCresult struct {
	Status string `json:"status"`
	Data   []SRCData
}
