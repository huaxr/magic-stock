// @Time:       2019/12/1 下午3:11

package control

import (
	"errors"
	"magic/stock/core/store"
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
	// 获取预测数据列表 post 请求
	GetPredict(c *gin.Context)
	// 获取股票详情
	GetDetail(c *gin.Context)
	GetFunds(c *gin.Context)
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
	//where, args := check.ParamParse.ParseParamsToWhereArgs(c, []string{"condition"}, true)
	//log.Println(where, args)
	var post model.GetPredicts
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(400, gin.H{"error_code": 1, "err_msg": err.Error(), "data": nil})
		return
	}
	var predicts []dal.Predict
	tmp := store.MysqlClient.GetDB().Model(&dal.Predict{})
	for _, i := range post.Predicts {
		tmp = tmp.Where("`condition` regexp ?", i)
	}
	tmp.Limit(limit).Offset(offset).Find(&predicts)
	d.Response(c, predicts, nil)
}

func (d *PredictControl) GetDetail(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		d.Response(c, nil, errors.New("证券代码为空"))
	}
	var TicketHistory []dal.TicketHistory
	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ?", code).Limit(50).Order("date desc").Find(&TicketHistory)

	var TicketHistoryWeekly []dal.TicketHistoryWeekly
	store.MysqlClient.GetDB().Model(&dal.TicketHistoryWeekly{}).Where("code = ?", code).Limit(20).Order("date desc").Find(&TicketHistoryWeekly)

	var Stockholder []dal.Stockholder
	store.MysqlClient.GetDB().Model(&dal.Stockholder{}).Where("code = ?", code).Find(&Stockholder)

	var Stock dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", code).Find(&Stock)

	var Predict dal.Predict
	store.MysqlClient.GetDB().Model(&dal.Predict{}).Where("code = ?", code).Find(&Predict)

	var StockCashFlow []dal.StockCashFlow
	store.MysqlClient.GetDB().Model(&dal.StockCashFlow{}).Where("code = ?", code).Find(&StockCashFlow)

	var StockLiabilities []dal.StockLiabilities
	store.MysqlClient.GetDB().Model(&dal.StockLiabilities{}).Where("code = ?", code).Find(&StockLiabilities)

	var StockProfit []dal.StockProfit
	store.MysqlClient.GetDB().Model(&dal.StockProfit{}).Where("code = ?", code).Find(&StockProfit)

	var response model.StockDetail
	response.TicketHistory = TicketHistory
	response.Stockholder = Stockholder
	response.Stock = Stock
	response.Predict = Predict
	response.StockCashFlow = StockCashFlow
	response.StockLiabilities = StockLiabilities
	response.StockProfit = StockProfit
	response.TicketHistoryWeekly = TicketHistoryWeekly
	d.Response(c, response, nil)
}

func (d *PredictControl) GetFunds(c *gin.Context) {
	offset, limit := check.ParamParse.GetPagination(c)
	code := c.DefaultQuery("code", "")
	if code == "" {
		d.Response(c, nil, errors.New("证券代码为空"))
	}
	var FundHoldRank []dal.FundHoldRank
	var Funds []model.StockFund
	store.MysqlClient.GetDB().Model(&dal.FundHoldRank{}).Where("code = ?", code).Offset(offset).Limit(limit).Find(&FundHoldRank)
	for _, i := range FundHoldRank {
		var fund dal.FundRank
		store.MysqlClient.GetDB().Model(&dal.FundRank{}).Where("fund_code = ?", i.FundCode).Find(&fund)
		Funds = append(Funds, model.StockFund{FundInfo: fund, Percent: i.Percent})
	}
	d.Response(c, Funds, nil)
}
