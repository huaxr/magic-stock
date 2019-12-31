package nac

import "time"

type NatRule struct {
	ID 			int  		`gorm:"primary_key"`
	BuzName 	string
	DeployType 	int
	SrcIP 		string  	`gorm:"column:src_ip"`
	SrcPSM 		string		`gorm:"column:src_psm"`
	ServiceID 	int 		`gorm:"column:service_id"`
	ServiceName string
	DestIP 		string  	`gorm:"column:dest_ip"`
	DestDomain 	string
	AccessType 	int
	Comments 	string 		`gorm:"type:text"`
	Creator 	string
	IsEnabled 	int
	IsActive 	int
	RuleClass 	int
	RuleType 	int
	IsRoot 		int
	CreateTime 	time.Time
	UpdateTime 	time.Time
}

func (NatRule) TableName() string {
	return "byte_security_nac_nat_rule"
}