package model

//拉去任务详情
type RulesTaskDetailForm struct {
	ProductName        string `json:"product_name"`
	ProductDescription string `json:"product_description"`
	RulesList          []*RuleForm
}

type RuleForm struct {
	RuleId    int    `json:"rule_id"`
	RuleName  string `json:"rule_name"`
	RuleExpr  string `json:"rule_expr"`
	RiskLevel int    `json:"risk_level"`
}
