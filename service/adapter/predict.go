// @Time:       2019/12/1 下午3:23

package adapter

import (
	"magic/stock/dal"
	"magic/stock/dao"
)

type PredictServiceIF interface {
	Create(event *dal.Predict) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*dal.Predict, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]dal.Predict, error)
	Count(where string, args []interface{}) (int, error)
}

var PredictServiceGlobal PredictServiceIF

func init() {
	tmp := new(PredictService)
	tmp.dao = dao.PredictDao
	PredictServiceGlobal = tmp
}

type PredictService struct {
	dao dao.PredictDaoIF
}

func (m *PredictService) Create(Predict *dal.Predict) error {
	return m.dao.Create(Predict)
}

func (m *PredictService) Delete(id int) error {
	return m.dao.Delete(id)
}

func (m *PredictService) Update(id int, ma map[string]interface{}) error {
	return m.dao.Update(id, ma)
}

func (m *PredictService) Query(where string, args []interface{}) (*dal.Predict, error) {
	return m.dao.Query(where, args)
}

func (m *PredictService) QueryAll(where string, args []interface{}, offset, limit int) (*[]dal.Predict, error) {
	return m.dao.QueryAll(where, args, offset, limit)
}

func (m *PredictService) Count(where string, args []interface{}) (int, error) {
	return m.dao.Count(where, args)
}
