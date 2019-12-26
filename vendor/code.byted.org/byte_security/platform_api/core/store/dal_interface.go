package store

import (
	"fmt"
	"reflect"
)

var DB ModelInterface = &Mysql{}

type ModelInterface interface {
	NewDefaultQueryParam() *QueryParam
	QueryList(param *QueryParam) (ret interface{}, total int, err error)
	QueryJoinList(param *QueryParam, result interface{}) (int, error)
	QueryCount(param *QueryParam) (total int, err error)
	QueryOne(param *QueryParam) (ret interface{}, err error)
	CreateObject(value interface{}) (err error)
	UpdateObject(model interface{}, value interface{}) (err error)
	DeleteObjects(model, ids interface{}, field ...string) (err error)
	DeleteObjectsCascade(model, relateModel, ids interface{}, field string) error
	DeleteAndCreate(model interface{}, field string, id interface{}, value ...interface{}) error
}

type QueryParam struct {
	Table   interface{}   `json:"table"`
	Select  string        `json:"select"`
	Preload string        `json:"preload"`
	Join    string        `json:"join"`
	Where   string        `json:"where"`
	Args    []interface{} `json:"args"`
	Offset  int           `json:"offset"`
	Limit   int           `json:"limit"`
	OrderBy string        `json:"order_by"`
}

func (param *QueryParam) AppendWhere(where string, args []interface{}) {
	if len(where) <= 0 {
		return
	}
	if len(param.Where) > 0 {
		param.Where += " and " + where
	} else {
		param.Where = where
	}
	param.Args = append(param.Args, args...)
}

func (c *Mysql) NewDefaultQueryParam() *QueryParam {
	return &QueryParam{
		Select:  "",
		Where:   "",
		Args:    nil,
		Offset:  0,
		Limit:   20,
		OrderBy: "created_at desc",
	}
}

func (c *Mysql) QueryList(param *QueryParam) (interface{}, int, error) {
	var total int
	ret := reflect.New(reflect.TypeOf(param.Table)).Interface()
	model := reflect.New(reflect.TypeOf(param.Table).Elem()).Interface()
	query := c.db.Model(model)
	if len(param.Where) > 0 {
		query = query.Where(param.Where, param.Args...)
	}
	err := query.Count(&total).Error
	if err != nil {
		return ret, total, err
	}
	if len(param.Select) > 0 {
		query = query.Select(param.Select)
	}
	if len(param.Preload) > 0 {
		query = query.Preload(param.Preload)
	}
	if param.Limit > 0 {
		query = query.Offset(param.Offset).Limit(param.Limit)
	}
	err = query.Order(param.OrderBy).Find(ret).Error
	return ret, total, err
}

func (c *Mysql) QueryJoinList(param *QueryParam, result interface{}) (int, error) {
	var total int
	model := reflect.New(reflect.TypeOf(param.Table).Elem()).Interface()
	query := c.db.Model(model)
	if len(param.Where) > 0 {
		query = query.Where(param.Where, param.Args...)
	}
	err := query.Count(&total).Error
	if err != nil {
		return total, err
	}
	if len(param.Select) > 0 {
		query = query.Select(param.Select)
	}
	if len(param.Join) > 0 {
		query = query.Joins(param.Join)
	}
	if len(param.Preload) > 0 {
		query = query.Preload(param.Preload)
	}
	if param.Limit > 0 {
		query = query.Offset(param.Offset).Limit(param.Limit)
	}
	err = query.Order(param.OrderBy).Scan(result).Error
	return total, err
}

func (c *Mysql) QueryCount(param *QueryParam) (total int, err error) {
	model := reflect.New(reflect.TypeOf(param.Table)).Interface()
	query := c.db.Model(model)
	if len(param.Where) > 0 {
		query = query.Where(param.Where, param.Args...)
	}
	err = query.Count(&total).Error
	return
}

func (c *Mysql) QueryOne(param *QueryParam) (interface{}, error) {
	model := reflect.New(reflect.TypeOf(param.Table)).Interface()
	query := c.db.Model(model)
	if len(param.Select) > 0 {
		query = query.Select(param.Select)
	}
	if len(param.Where) > 0 {
		query = query.Where(param.Where, param.Args...)
	}
	if len(param.Preload) > 0 {
		query = query.Preload(param.Preload)
	}
	err := query.Order(param.OrderBy).Find(model).Error
	return model, err
}

func (c *Mysql) CreateObject(value interface{}) (err error) {
	err = c.db.Create(value).Error
	return
}

func (c *Mysql) UpdateObject(model, value interface{}) (err error) {
	err = c.db.Model(model).Updates(value).Error
	return
}

func (c *Mysql) DeleteObjects(model, ids interface{}, field ...string) (err error) {
	if len(field) > 0 {
		err = c.db.Delete(model, fmt.Sprintf("%s in (?)", field[0]), ids).Error
	} else {
		err = c.db.Delete(model, "id in (?)", ids).Error
	}
	return
}

func (c *Mysql) DeleteObjectsCascade(model, relateModel, ids interface{}, field string) error {
	tx := c.db.Begin()
	if tx.Where(fmt.Sprintf("%s in (?)", field), ids).Delete(relateModel); tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	if tx.Where("id in (?)", ids).Delete(model); tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	return tx.Commit().Error
}

func (c *Mysql) DeleteAndCreate(model interface{}, field string, id interface{}, value ...interface{}) error {
	tx := c.db.Begin()
	if tx.Where(fmt.Sprintf("%s = ?", field), id).Delete(model); tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	for _, item := range value {
		err := tx.Create(item).Error
		if err != nil {
			tx.Rollback()
			return tx.Error
		}
	}
	return tx.Commit().Error
}
