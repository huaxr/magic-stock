// @Time:       2019/12/1 下午3:11

package control

import (
	"log"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/service/adapter"
	"magic/stock/service/check"

	"github.com/gin-gonic/gin"
)

type PredictIF interface {
	Query(where string, args []interface{}) (*dal.Predict, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]dal.Predict, error)
	Exist(where string, args []interface{}) bool
	GetPredict(c *gin.Context)
	Response(c *gin.Context, data interface{}, err error)
}

var PredictControlGlobal PredictIF

func init() {
	tmp := new(PredictControl)
	tmp.service = adapter.PredictServiceGlobal
	tmp.response = new(model.HttpResponse)
	PredictControlGlobal = tmp
}

type PredictControl struct {
	service  adapter.PredictServiceIF
	response *model.HttpResponse
}

func (u *PredictControl) Query(where string, args []interface{}) (*dal.Predict, error) {
	return u.service.Query(where, args)
}

func (u *PredictControl) QueryAll(where string, args []interface{}, offset, limit int) (*[]dal.Predict, error) {
	return u.service.QueryAll(where, args, offset, limit)
}

func (u *PredictControl) Exist(where string, args []interface{}) bool {
	c, _ := u.service.Count(where, args)
	return c > 0
}

func (d *PredictControl) Response(c *gin.Context, data interface{}, err error) {
	c.AbortWithStatusJSON(200, d.response.Response(data, err))
}

func (d *PredictControl) GetPredict(c *gin.Context) {
	offset, limit := check.ParamParse.GetPagination(c)
	where, args := check.ParamParse.ParseParamsToWhereArgs(c, []string{"condition"}, true)
	log.Println(where, args)
	pres, _ := d.QueryAll(where, args, offset, limit)
	d.Response(c, pres, nil)
}
