package nac

type NatTicket struct {
	ID         int `gorm:"primary_key"`
	NatRuleID  int `gorm:"column:nat_rule_id"`
	WorkflowID int `gorm:"workflow_id"`
}
