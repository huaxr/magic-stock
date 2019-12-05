// @Contact:    huaxinrui
// @Time:       2019/9/24 上午11:34

package models

type NewQuery struct {
	Type     interface{}
	Where    interface{}
	Args     []interface{}
	Distinct string
	Limit    int
	Offset   int
}
