// @Contact:    huaxinrui
// @Time:       2019/9/23 上午11:28

package store

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/dal/event"
	"code.byted.org/byte_security/dal/policy"
	"code.byted.org/byte_security/dal/rule_engine"
	"code.byted.org/byte_security/dal/seal"
	"code.byted.org/byte_security/dal/soc"
	"code.byted.org/byte_security/dal/temp"
	"code.byted.org/byte_security/dal/workflow"
	"code.byted.org/byte_security/platform_api/models"
	"code.byted.org/gopkg/gorm"
	_ "code.byted.org/gopkg/mysql-driver"
)

type MI interface {
	GetDB() *gorm.DB
	NewTransaction() *Transaction
	QueryJson(table, filed, k string, v interface{}) (*gorm.DB, error)
}

type Mysql struct {
	StorageClient
	db   *gorm.DB
	pool *sync.Pool
}

func InitStore(path string, debug bool) *Mysql {
	e := new(Mysql)
	db, err := gorm.Open("mysql2", path)
	if err != nil {
		panic(err)
	}
	db.DB().SetConnMaxLifetime(60 * time.Second)
	db.DB().SetMaxOpenConns(30)
	e.typ = "mysql"
	e.mutex = new(sync.Mutex)
	e.debug = debug
	db.LogMode(debug)
	e.db = db
	e.pool = &sync.Pool{
		New: func() interface{} {
			return new(models.NewQuery)
		},
	}
	go db.AutoMigrate(&auth.Group{}, &auth.User{}, &auth.AuthApply{}, &auth.CasbinRule{}, &auth.Policy{}, &auth.Role{},
		&auth.Token{}, &auth.GroupProduct{}, &auth.UserProduct{}, &auth.UserRole{}, &auth.UserBusiness{},
		&temp.TmpLink{}, &temp.TmpRecord{}, // 临时链接表
		&soc.App{}, &soc.Domain{}, &soc.AssetOwner{}, &soc.Website{}, &soc.DomainPSM{}, &soc.PSM{}, &soc.Repo{}, &soc.PSMHost{},
		&soc.Host{}, &soc.IDC{}, &soc.NetWork{}, &soc.Product{}, &soc.M2MOuter{}, &soc.PSMRepo{}, &soc.Vulnerability{},
		&soc.VulnerabilityType{}, &soc.CrawlerSeebug{}, &soc.Business{},
		&workflow.Ticket{}, &workflow.WorkFlow{}, &workflow.Record{}, &workflow.NodeUser{},
		&workflow.Node{}, &workflow.NodeTemplate{}, &workflow.Subscribe{}, &workflow.TypeList{},
		&policy.RuleExpress{}, &policy.Field{}, &policy.Group{}, &policy.GroupKey{}, &policy.Location{}, &policy.GroupAndStrategyOnline{},
		&policy.ExpressAndStrategyVersionOffline{}, &policy.GroupAndLocation{}, &policy.Product{},
		&policy.Rule{}, &policy.StrategyOffline{}, &policy.StrategyOnline{}, &policy.StrategyVersionOffline{},
		&policy.StrategyVersionOnline{}, &policy.Task{}, &policy.VariateOffline{}, &policy.VariateObjOffline{},
		&policy.VariateOnline{}, &policy.VariateObjOnline{},
		&seal.Device{}, &seal.User{},
		&event.HResults{},
		&rule_engine.Task{}, &rule_engine.TaskDataResult{}, &rule_engine.TaskDataSource{})
	return e
}

// normal query
func (m *Mysql) Query(query *models.NewQuery) (interface{}, error) {
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

func (e *Mysql) Count(query *models.NewQuery) (int, error) {
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

func (e *Mysql) NewQuery() *models.NewQuery {
	q := e.pool.Get().(*models.NewQuery)
	defer e.pool.Put(q)
	// bug report: the pointer object will use limit offset
	// so extra query params should clear before return and put back to pool
	q.Distinct = ""
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

func (e *Mysql) Safe() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
}
