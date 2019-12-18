// @Time:       2019/12/1 下午3:11

package control

import (
	"errors"
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/service/adapter"
	"magic/stock/service/check"
	"strings"

	"github.com/gin-gonic/gin"
)

type PredictIF interface {
	Query(where string, args []interface{}) (*dal.Predict, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]dal.Predict, error)
	Exist(where string, args []interface{}) bool
	// 获取预测数据列表 post 请求
	PredictList(c *gin.Context)
	// 获取预测的时间点
	GetPredictDates(c *gin.Context)
	GetConditions(c *gin.Context)
	// 获取股票详情
	GetDetail(c *gin.Context)
	GetFunds(c *gin.Context)
	// 通过机构code查询机构持仓
	FundHold(c *gin.Context)
	// 通过名称查询流通股东可能存在的其它持仓
	TopHolderHold(c *gin.Context)
	// 获取所有行业的列表
	GetBelongs(c *gin.Context)
	// 获取所有地区的列表
	GetLocations(c *gin.Context)
	// 获取所有的组织形式列表
	GetOrganizationalForms(c *gin.Context)
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

type Codes struct {
	Code string
}

func (d *PredictControl) PredictList(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	offset, limit := check.ParamParse.GetPagination(c)
	var post model.GetPredicts
	err := c.BindJSON(&post)
	if err != nil {
		d.Response(c, nil, err)
	}
	// 如果用户提交查询并保存查询结果
	if post.Save {
		err := adapter.UserServiceGlobal.SaveUserConditions(&post, authentication)
		if err != nil {
			log.Println("保存用户查询数据失败", err)
		}
	}
	var where_belongs, where_locations, where_forms []string
	var args_belongs, args_locationgs, args_forms []interface{}
	var belongs, locations, forms []string

	if len(post.Query.Belongs) > 0 {
		var codes []Codes
		for _, i := range post.Query.Belongs {
			where_belongs = append(where_belongs, "belong = ?")
			args_belongs = append(args_belongs, i)
		}
		where_str := strings.Join(where_belongs, "OR")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_belongs...).Scan(&codes)
		for _, i := range codes {
			belongs = append(belongs, i.Code)
		}
	}

	if len(post.Query.Locations) > 0 {
		var codes []Codes
		for _, i := range post.Query.Locations {
			where_locations = append(where_locations, "location = ?")
			args_locationgs = append(args_locationgs, i)
		}
		where_str := strings.Join(where_locations, "OR")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_locationgs...).Scan(&codes)
		for _, i := range codes {
			locations = append(locations, i.Code)
		}
	}

	if len(post.Query.OrganizationalForm) > 0 {
		var codes []Codes
		for _, i := range post.Query.OrganizationalForm {
			where_forms = append(where_forms, "organizational_form = ?")
			args_forms = append(args_forms, i)
		}
		where_str := strings.Join(where_forms, "OR")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_forms...).Scan(&codes)
		for _, i := range codes {
			forms = append(forms, i.Code)
		}
	}

	var predicts []dal.Predict
	tmp := store.MysqlClient.GetDB().Model(&dal.Predict{}).Where("date = ?", post.Date)

	if len(belongs) > 0 {
		tmp = tmp.Where("code IN (?)", belongs)
	}
	if len(locations) > 0 {
		tmp = tmp.Where("code IN (?)", locations)
	}
	if len(forms) > 0 {
		tmp = tmp.Where("code IN (?)", forms)
	}
	for _, i := range post.Query.Predicts {
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
	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ?", code).Limit(70).Order("date desc").Find(&TicketHistory)

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

func (d *PredictControl) FundHold(c *gin.Context) {
	code := c.DefaultQuery("fund_code", "")
	if code == "" {
		d.Response(c, nil, errors.New("机构代码为空"))
	}
	var FundHoldRank []dal.FundHoldRank
	store.MysqlClient.GetDB().Model(&dal.FundHoldRank{}).Where("fund_code = ?", code).Find(&FundHoldRank)
	d.Response(c, FundHoldRank, nil)
}

func (d *PredictControl) TopHolderHold(c *gin.Context) {
	holder := c.DefaultQuery("holder_name", "")
	if holder == "" {
		d.Response(c, nil, errors.New("查询用户为空"))
	}
	var Stockholder []dal.Stockholder
	store.MysqlClient.GetDB().Model(&dal.Stockholder{}).Where("holder_name = ?", holder).Find(&Stockholder)
	d.Response(c, Stockholder, nil)
}

type PredictDate struct {
	Date string `json:"date"`
}

type Belongs struct {
	Belong string `json:"date"`
}

type Locations struct {
	Location string `json:"location"`
}

type OrganizationalForms struct {
	OrganizationalForm string `json:"organizational_form"`
}

func (d *PredictControl) GetPredictDates(c *gin.Context) {
	var x []PredictDate
	store.MysqlClient.GetDB().Model(&dal.Predict{}).Select("distinct(date) as date").Order("date desc").Scan(&x)
	d.Response(c, x, nil)
}

func (d *PredictControl) GetBelongs(c *gin.Context) {
	var x []Belongs
	var result []string
	store.MysqlClient.GetDB().Model(&dal.Code{}).Select("distinct(belong) as belong").Order("belong").Scan(&x)
	for _, i := range x {
		if i.Belong != "" {
			result = append(result, i.Belong)
		}
	}
	d.Response(c, result, nil)
}

func (d *PredictControl) GetLocations(c *gin.Context) {
	var x []Locations
	var result []string
	store.MysqlClient.GetDB().Model(&dal.Code{}).Select("distinct(location) as location").Order("location").Scan(&x)
	for _, i := range x {
		if i.Location != "" {
			result = append(result, i.Location)
		}
	}
	d.Response(c, result, nil)
}

func (d *PredictControl) GetOrganizationalForms(c *gin.Context) {
	var x []OrganizationalForms
	var result []string
	store.MysqlClient.GetDB().Model(&dal.Code{}).Select("distinct(organizational_form) as organizational_form").Order("organizational_form").Scan(&x)
	for _, i := range x {
		if i.OrganizationalForm != "" {
			result = append(result, i.OrganizationalForm)
		}
	}
	d.Response(c, result, nil)
}

func (d *PredictControl) GetConditions(c *gin.Context) {
	var x []dal.Conditions
	response := map[string][]string{}
	store.MysqlClient.GetDB().Model(&dal.Conditions{}).Find(&x)
	for _, i := range x {
		response[i.Type] = append(response[i.Type], i.Name)
	}
	d.Response(c, response, nil)
}
