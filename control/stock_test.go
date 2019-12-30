// @Time:       2019/12/15 下午7:35

package control

import (
	"log"
	"magic/stock/core/store"
	"strings"
	"testing"
)

type Location struct {
	Location string
}

// 获取地区
func TestGetLocations(t *testing.T) {
	rows, _ := store.MysqlClient.GetDB().Raw("select distinct(location) from magic_stock_code order by convert(location using gbk)").Rows()
	var locations []string
	for rows.Next() {
		var l Location
		// ScanRows scan a row into user
		store.MysqlClient.GetDB().ScanRows(rows, &l)
		if l.Location != "" {
			locations = append(locations, l.Location)
		}
	}
	log.Println(strings.Join(locations, ","))
}

type Belong struct {
	Belong string
}

// 获取belong
func TestGetBelongs(t *testing.T) {
	rows, _ := store.MysqlClient.GetDB().Raw("select distinct(belong) from magic_stock_code order by convert(belong using gbk)").Rows()
	var belongs []string
	for rows.Next() {
		var l Belong
		// ScanRows scan a row into user
		store.MysqlClient.GetDB().ScanRows(rows, &l)
		if l.Belong != "" {
			belongs = append(belongs, l.Belong)
		}
	}
	log.Println(strings.Join(belongs, ","))
}

type Concept struct {
	Name string
}

// 获取concept
func TestGetConcepts(t *testing.T) {
	rows, _ := store.MysqlClient.GetDB().Raw("select distinct(name) from magic_stock_concept order by convert(name using gbk)").Rows()
	var belongs []string
	for rows.Next() {
		var l Concept
		// ScanRows scan a row into user
		store.MysqlClient.GetDB().ScanRows(rows, &l)
		if l.Name != "" {
			belongs = append(belongs, l.Name)
		}
	}
	log.Println(strings.Join(belongs, ","))
}

// 获取labels
func TestGetLabels(t *testing.T) {
	rows, _ := store.MysqlClient.GetDB().Raw("select distinct(name) from magic_stock_labels order by convert(name using gbk)").Rows()
	var belongs []string
	for rows.Next() {
		var l Concept
		// ScanRows scan a row into user
		store.MysqlClient.GetDB().ScanRows(rows, &l)
		if l.Name != "" && l.Name != "昨日涨停" {
			belongs = append(belongs, l.Name)
		}
	}
	log.Println(strings.Join(belongs, ","))
}
