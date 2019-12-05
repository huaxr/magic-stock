// @Contact:    huaxinrui
// @Time:       2019/9/26 下午6:14

package dao

import (
	"magic/stock/core/store"
	"magic/stock/dal"
)

type PredictDaoIF interface {
	Create(Predict *dal.Predict) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*dal.Predict, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]dal.Predict, error)
	Count(where string, args []interface{}) (int, error)
}

var PredictDao PredictDaoIF

func init() {
	tmp := new(MysqlPredict)
	tmp.Store = store.MysqlClient
	PredictDao = tmp
}

type MysqlPredict struct {
	Store store.StoreIF
}

func (m *MysqlPredict) Create(Predict *dal.Predict) error {
	return m.Store.GetDB().Model(Predict).Create(Predict).Error
}

func (m *MysqlPredict) Delete(id int) error {
	return m.Store.GetDB().Delete(&dal.Predict{}, "id = ?", id).Error
}

func (m *MysqlPredict) Update(id int, ma map[string]interface{}) error {
	return m.Store.GetDB().Update(&dal.Predict{}, ma).Error
}

func (m *MysqlPredict) Query(where string, args []interface{}) (*dal.Predict, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = dal.Predict{}
	query_obj.Where = where
	query_obj.Args = args
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*dal.Predict), nil
}

func (m *MysqlPredict) Count(where string, args []interface{}) (int, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = dal.Predict{}
	query_obj.Where = where
	query_obj.Args = args
	return m.Store.Count(query_obj)
}

func (m *MysqlPredict) QueryAll(where string, args []interface{}, offset, limit int) (*[]dal.Predict, error) {
	query_obj := m.Store.NewQuery()
	query_obj.Type = []dal.Predict{}
	query_obj.Where = where
	query_obj.Args = args
	query_obj.Limit = limit
	query_obj.Offset = offset
	result, err := m.Store.Query(query_obj)
	if err != nil {
		return nil, err
	}
	return result.(*[]dal.Predict), nil
}
