package model

type RulesTaskDetail struct {
	Rules []*Rule    //训练的规则
	Vars  []*Variate //如果规则表达式expr中包含变量表达式，则以map的方式拉取获得
}
