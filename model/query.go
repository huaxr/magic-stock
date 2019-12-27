// @Time:       2019/11/27 下午8:25

package model

type NewQuery struct {
	Type       interface{}
	Where      interface{}
	Args       []interface{}
	SelectOnly string
	Distinct   string
	Limit      int
	Offset     int
}
