// @Contact:    huaxinrui
// @Time:       2019/9/23 上午11:28

package store

import (
	"errors"
	"fmt"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/service/conf"
	"magic/stock/utils"
	"reflect"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	MysqlClient StoreIF
)

type MI interface {
	GetDB() *gorm.DB
	GetOnlineDB() *gorm.DB
	NewTransaction() *Transaction
	QueryJson(table, filed, k string, v interface{}) (*gorm.DB, error)
}

type Mysql struct {
	StorageClient
	db       *gorm.DB
	dbonline *gorm.DB
	pool     *sync.Pool
}

func init() {
	e := new(Mysql)
	db, err := gorm.Open("mysql", conf.Config.Store)
	if err != nil {
		panic(err)
	}
	db.DB().SetConnMaxLifetime(60 * time.Second)
	db.DB().SetMaxOpenConns(30)

	if utils.TellEnv() == "loc" {
		e.dbonline, err = gorm.Open("mysql", "root:Icode_787518771@tcp(154.8.228.39:3306)/stock?charset=utf8mb4&parseTime=True&loc=Local")
		if err != nil {
			panic(err)
		}
		e.dbonline.DB().SetConnMaxLifetime(60 * time.Second)
		e.dbonline.DB().SetMaxOpenConns(30)
		go e.dbonline.AutoMigrate(&dal.Code{}, &dal.Stockholder{},
			&dal.Predict{}, &dal.TicketHistoryWeekly{}, dal.TicketHistory{}, &dal.User{}, &dal.Pay{},
			&dal.StockCashFlow{}, &dal.StockProfit{}, &dal.StockLiabilities{}, &dal.Conditions{}, &dal.UserConditions{},
			&dal.StockConcept{}, &dal.StockLabels{}, &dal.StockPerTicket{}, &dal.StockFund{}, &dal.StockFengHong{}, &dal.StockPeiGu{},
			&dal.StockZengFa{}, &dal.Price{}, &dal.UserShare{}, &dal.HighConditions{}, &dal.UserDemands{}, &dal.StockSubCompany{}, &dal.HistoryALL1{}, &dal.HistoryALL2{},
			&dal.TicketHistoryWeeklyALL{}, &dal.TicketHistoryMonth{}, &dal.TicketHistoryMonthAll{}, &dal.StockPublicNews{}, &dal.StockPublicReports{})
	}

	e.typ = "mysql"
	e.mutex = new(sync.Mutex)
	e.debug = false
	db.LogMode(e.debug)
	e.db = db
	e.pool = &sync.Pool{
		New: func() interface{} {
			return new(model.NewQuery)
		},
	}
	go db.AutoMigrate(&dal.Code{}, &dal.Stockholder{},
		&dal.Predict{}, &dal.TicketHistoryWeekly{}, dal.TicketHistory{}, &dal.User{}, &dal.Pay{},
		&dal.StockCashFlow{}, &dal.StockProfit{}, &dal.StockLiabilities{}, &dal.Conditions{}, &dal.UserConditions{},
		&dal.StockConcept{}, &dal.StockLabels{}, &dal.StockPerTicket{}, &dal.StockFund{}, &dal.StockFengHong{}, &dal.StockPeiGu{},
		&dal.StockZengFa{}, &dal.Price{}, &dal.UserShare{}, &dal.HighConditions{}, &dal.UserDemands{}, &dal.StockSubCompany{}, &dal.HistoryALL1{}, &dal.HistoryALL2{},
		&dal.TicketHistoryWeeklyALL{}, &dal.TicketHistoryMonth{}, &dal.TicketHistoryMonthAll{}, &dal.StockPublicNews{}, &dal.StockPublicReports{})
	MysqlClient = e
}

// normal query
func (m *Mysql) Query(query *model.NewQuery) (interface{}, error) {
	var a interface{}
	var err error
	var tmp *gorm.DB
	t := reflect.TypeOf(query.Type)
	if strings.HasPrefix(t.String(), "[]") {
		a = reflect.New(reflect.TypeOf(query.Type)).Interface()
		b := reflect.New(reflect.TypeOf(query.Type).Elem()).Interface()
		tmp = m.db.Model(b)
	} else {
		a = reflect.New(reflect.TypeOf(query.Type)).Interface()
		tmp = m.db.Model(a)
	}
	if query.Where != nil {
		tmp = tmp.Where(query.Where, query.Args...)
	}
	if query.SelectOnly != "" {
		tmp = tmp.Select(query.SelectOnly)
	}
	if query.Limit != 0 {
		tmp = tmp.Limit(query.Limit)
	}
	if query.Offset != 0 {
		tmp = tmp.Offset(query.Offset)
	}
	if query.Distinct != "" {
		tmp = tmp.Group(query.Distinct)
	}
	err = tmp.Order("id desc").Find(a).Error
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (e *Mysql) Count(query *model.NewQuery) (int, error) {
	var err error
	var count int
	var tmp *gorm.DB
	a := reflect.New(reflect.TypeOf(query.Type)).Interface()
	if query.Where != nil {
		tmp = e.db.Model(a).Where(query.Where, query.Args...)
	} else {
		tmp = e.db.Model(a)
	}
	if query.Distinct != "" {
		tmp = tmp.Group(query.Distinct)
	}
	err = tmp.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (e *Mysql) GetType() string {
	return e.typ
}

func (e *Mysql) GetDB() *gorm.DB {
	return e.db
}

func (e *Mysql) GetOnlineDB() *gorm.DB {
	return e.dbonline
}

func (e *Mysql) NewQuery() *model.NewQuery {
	q := e.pool.Get().(*model.NewQuery)
	defer e.pool.Put(q)
	// bug report: the pointer object will use limit offset
	// so extra query params should clear before return and put back to pool
	q.Distinct = ""
	q.SelectOnly = ""
	q.Limit = 0
	q.Offset = 0
	return q
}

func (e *Mysql) begin() *gorm.DB {
	return e.db.Begin()
}

func (e *Mysql) NewTransaction() *Transaction {
	return &Transaction{TX: e.begin()}
}

func (e *Mysql) judgeSql(filed ...string) error {
	sql := "'\";%="
	for _, i := range filed {
		valid_filed := strings.ContainsAny(i, sql)
		if valid_filed {
			return errors.New("Field sql injection")
		}
	}
	return nil
}

// query json option
// bug report: k can never be simple int argument as suffix, like: 22ss
func (e *Mysql) QueryJson(table, field, k string, v interface{}) (*gorm.DB, error) {
	err := e.judgeSql(field, k)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return e.db.Table(table).Where(fmt.Sprintf("%s->'$.%s' != ?", field, k), ""), nil
	}
	return e.db.Table(table).Where(fmt.Sprintf("%s->'$.%s' = ?", field, k), v.(string)), nil
}
