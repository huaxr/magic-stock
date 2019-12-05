// @Contact:    huaxinrui
// @Time:       2019/7/30 下午3:33

package rule_engine

type VarTaskDetail struct {
	VarName  string
	VarExpr  string
	VarsExpr map[string]string
}

type RuleTaskDetail struct {
	RuleName string
	RuleExpr string
	VarsExpr map[string]string //如果规则表达式expr中包含变量表达式，则以map的方式拉取获得
}

type PolicyTaskDetail struct {
	GroupKeys  []string
	GroupValue []string
	Pid        int    //策略ID
	PolicyName string //策略名称
	Type       int    //策略类型
	Priority   int    //优先级信息
	Rule       []*RuleOffline
	Vars       []*VarsOffline
}

type VarsOffline struct {
	Name     string   //变量名称
	Type     string   //变量类型
	SortType int      //分类类型
	Value    []string //value类型的值
}

type RuleOffline struct {
	RuleName string //规则名称
	Express  string //表达式
}
