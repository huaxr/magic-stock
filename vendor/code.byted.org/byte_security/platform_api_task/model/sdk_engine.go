package model

type ProductAndPolicy struct {
	Product      string          `json:"product"`
	Key          []string        `json:"key"`
	RelationList []*Relation     `json:"relation_list"`
	Policy       []*PolicyDetail `json:"policy"`
	State        int             `json:"state"`
}

type Relation struct {
	Sign      []string `json:"group"`
	PolicyIds []int    `json:"policy_ids"`
	GID       int      `json:"gid"`
	State     int      `json:"state"` //启用0，灰度1，禁用2
}

type PolicyDetail struct {
	PID  int      `json:"pid"`
	VID  int      `json:"vid"`
	Vars []string `json:"vars"`
}

type GetPolicy struct {
	PID        int        `json:"pid"`
	VID        int        `json:"vid"`
	PolicyName string     `json:"policy_name"`
	Type       int        `json:"type"`
	Priority   int        `json:"priority"`
	Version    string     `json:"version"`
	Filter     string     `json:"filter"`
	RuleList   []*Rule    `json:"rule_list"`
	Vars       []*Variate `json:"vars"`
}

type Rule struct {
	RuleId   int    `json:"rule_id"`
	RuleName string `json:"rule_name"`
	Express  string `json:"express"`
	Status   int    `json:"status"`
}

type Variate struct {
	VarName  string   `json:"var_name"`
	VarType  string   `json:"var_type"`  //变量值类型
	SortType int      `json:"sort_type"` //变量类型，1字面值，2列表，3表达式
	Value    []string `json:"value"`
}
