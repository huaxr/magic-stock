// @Time:       2019/12/1 下午3:11

package control

import (
	"errors"
	"fmt"
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/service/adapter"
	"magic/stock/service/check"
	"magic/stock/utils"
	"strings"

	"gopkg.in/fatih/set.v0"

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
	// 获取概念列表
	GetConcepts(c *gin.Context)
	GetLabels(c *gin.Context)
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

var (
	PredictControlGlobal PredictIF
	OrderLimit           = []string{"score", "percent", "price", "fund_count", "sm_count"}
)

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

func (d *PredictControl) getMinMax(da map[string]float64) (float64, float64) {
	min, ok_min := da["min"]
	if !ok_min {
		min = -9999
	}
	max, ok_max := da["max"]
	if !ok_max {
		max = 9999
	}
	return min, max
}

func (d *PredictControl) ParseStockPerTicket(param map[string]float64, field string) set.Interface {
	tmp := set.New(set.ThreadSafe)
	if len(param) > 0 {
		min, max := d.getMinMax(param)
		var codes []Codes
		store.MysqlClient.GetDB().Model(&dal.StockPerTicket{}).Select("code").Where(fmt.Sprintf("%s >= ? and %s <= ?", field, field), min, max).Scan(&codes)
		for _, i := range codes {
			tmp.Add(i.Code)
		}
	}
	return tmp
}

func (d *PredictControl) ParseLastDayRange(param map[string]float64, date string, field string) set.Interface {
	tmp := set.New(set.ThreadSafe)
	if len(param) > 0 {
		min, max := d.getMinMax(param)
		var codes []Codes
		store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Select("code").Where(fmt.Sprintf("date = ? and %s >= ? and %s <= ?", field, field), date, min, max).Scan(&codes)
		for _, i := range codes {
			tmp.Add(i.Code)
		}
	}
	return tmp
}

type Codes struct {
	Code string
}

func (d *PredictControl) PredictList(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	var post model.GetPredicts
	err := c.BindJSON(&post)
	if err != nil {
		d.Response(c, nil, err)
		return
	}

	if !authentication.Member {
		if authentication.QueryLeft == 0 {
			d.Response(c, nil, errors.New("查询次数不足"))
			return
		} else {
			user_obj, _ := UserControlGlobal.Query("id = ?", []interface{}{authentication.Uid})
			left := user_obj.QueryLeft - 1
			exp := user_obj.Exp + 1
			user_obj.QueryLeft = left
			user_obj.Exp = exp
			err := store.MysqlClient.GetDB().Save(&user_obj).Error
			log.Println("查询次数剩余", authentication.User, authentication.QueryLeft, left, err)
		}
	}
	offset, limit := check.ParamParse.GetPagination(c)
	// 如果用户提交查询并保存查询结果
	if post.Save {
		err := adapter.UserServiceGlobal.SaveUserConditions(&post, authentication)
		if err != nil {
			log.Println("保存用户查询数据失败", err)
		}
	}
	var where_belongs, where_locations, where_concepts []string
	var args_belongs, args_locationgs, args_concepts []interface{}

	belong_set := set.New(set.ThreadSafe)
	location_set := set.New(set.ThreadSafe)
	concept_set := set.New(set.ThreadSafe)
	per_ticket_set1 := set.New(set.ThreadSafe)
	per_ticket_set2 := set.New(set.ThreadSafe)
	per_ticket_set3 := set.New(set.ThreadSafe)
	per_ticket_set4 := set.New(set.ThreadSafe)
	per_ticket_set5 := set.New(set.ThreadSafe)
	per_ticket_set6 := set.New(set.ThreadSafe)

	last_day_set1 := set.New(set.ThreadSafe)
	last_day_set2 := set.New(set.ThreadSafe)
	last_day_set3 := set.New(set.ThreadSafe)
	last_day_set4 := set.New(set.ThreadSafe)

	ability_set1 := set.New(set.ThreadSafe)
	ability_set2 := set.New(set.ThreadSafe)
	ability_set3 := set.New(set.ThreadSafe)
	ability_set4 := set.New(set.ThreadSafe)
	ability_set5 := set.New(set.ThreadSafe)
	ability_set6 := set.New(set.ThreadSafe)
	ability_set7 := set.New(set.ThreadSafe)
	ability_set8 := set.New(set.ThreadSafe)
	ability_set9 := set.New(set.ThreadSafe)
	ability_set10 := set.New(set.ThreadSafe)
	ability_set11 := set.New(set.ThreadSafe)
	ability_set12 := set.New(set.ThreadSafe)
	ability_set13 := set.New(set.ThreadSafe)
	ability_set14 := set.New(set.ThreadSafe)
	ability_set15 := set.New(set.ThreadSafe)
	ability_set16 := set.New(set.ThreadSafe)
	ability_set17 := set.New(set.ThreadSafe)

	if len(post.Query.Belongs) > 0 {
		var codes []Codes
		for _, i := range post.Query.Belongs {
			where_belongs = append(where_belongs, "belong = ?")
			args_belongs = append(args_belongs, i)
		}
		where_str := strings.Join(where_belongs, " OR ")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_belongs...).Scan(&codes)
		for _, i := range codes {
			belong_set.Add(i.Code)
		}
	}

	if len(post.Query.Locations) > 0 {
		var codes []Codes
		for _, i := range post.Query.Locations {
			where_locations = append(where_locations, "location = ?")
			args_locationgs = append(args_locationgs, i)
		}
		where_str := strings.Join(where_locations, " OR ")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_locationgs...).Scan(&codes)
		for _, i := range codes {
			location_set.Add(i.Code)
		}
	}

	if len(post.Query.Concepts) > 0 || len(post.Query.Labels) > 0 {
		var codes []Codes
		arrays := append(post.Query.Concepts, post.Query.Labels...)
		for _, i := range arrays {
			where_concepts = append(where_concepts, "concept like ?")
			args_concepts = append(args_concepts, "%"+i+"%")
		}
		where_str := strings.Join(where_concepts, " OR ")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_concepts...).Scan(&codes)
		for _, i := range codes {
			concept_set.Add(i.Code)
		}
	}

	per_ticket_set1 = d.ParseStockPerTicket(post.Query.PerTickets.Shouyiafter, "shouyiafter")
	per_ticket_set2 = d.ParseStockPerTicket(post.Query.PerTickets.Jiaquanshouyi, "jiaquanshouyi")
	per_ticket_set3 = d.ParseStockPerTicket(post.Query.PerTickets.Jinzichanafter, "jinzichanafter")
	per_ticket_set4 = d.ParseStockPerTicket(post.Query.PerTickets.Jingyingxianjinliu, "jingyingxianjinliu")
	per_ticket_set5 = d.ParseStockPerTicket(post.Query.PerTickets.Gubengongjijin, "gubengongjijin")
	per_ticket_set6 = d.ParseStockPerTicket(post.Query.PerTickets.Weifenpeilirun, "weifenpeilirun")

	last_day_set1 = d.ParseLastDayRange(post.Query.LastDayRange.LastPercent, post.Date, "percent")
	last_day_set2 = d.ParseLastDayRange(post.Query.LastDayRange.LastAmplitude, post.Date, "amplitude")
	last_day_set3 = d.ParseLastDayRange(post.Query.LastDayRange.LastTurnoverrate, post.Date, "turnover_rate")
	last_day_set4 = d.ParseLastDayRange(post.Query.LastDayRange.LastPrice, post.Date, "shou")

	// 盈利能力
	ability_set1 = d.ParseStockPerTicket(post.Query.YlAbility.YlZongzichanlirunlv, "yl_zongzichanlirunlv")
	ability_set2 = d.ParseStockPerTicket(post.Query.YlAbility.YlZhuyingyewulirunlv, "yl_zhuyingyewulirunlv")
	ability_set3 = d.ParseStockPerTicket(post.Query.YlAbility.YlZongzichanjinglirunlv, "yl_zongzichanjinglirunlv")
	ability_set4 = d.ParseStockPerTicket(post.Query.YlAbility.YlYingyelirunlv, "yl_yingyelirunlv")
	ability_set5 = d.ParseStockPerTicket(post.Query.YlAbility.YlXiaoshoujinglilv, "yl_xiaoshoujinglilv")
	ability_set6 = d.ParseStockPerTicket(post.Query.YlAbility.YlGubenbaochoulv, "yl_gubenbaochoulv")
	ability_set7 = d.ParseStockPerTicket(post.Query.YlAbility.YlJingzichanbaochoulv, "yl_jingzichanbaochoulv")
	ability_set8 = d.ParseStockPerTicket(post.Query.YlAbility.YlZichanbaochoulv, "yl_zichanbaochoulv")
	// 成长能力
	ability_set9 = d.ParseStockPerTicket(post.Query.CzAbility.CzZhuyingyewushouruzengzhanglv, "cz_zhuyingyewushouruzengzhanglv")
	ability_set10 = d.ParseStockPerTicket(post.Query.CzAbility.CzJinglirunzengzhanglv, "cz_jinglirunzengzhanglv")
	ability_set11 = d.ParseStockPerTicket(post.Query.CzAbility.CzJingzichanzengzhanglv, "cz_jingzichanzengzhanglv")
	ability_set12 = d.ParseStockPerTicket(post.Query.CzAbility.CzZongzichanzengzhanglv, "cz_zongzichanzengzhanglv")
	// 运营能力
	ability_set13 = d.ParseStockPerTicket(post.Query.YyAbility.YyYingshouzhangkuanzhouzhuanlv, "yy_yingshouzhangkuanzhouzhuanlv")
	ability_set14 = d.ParseStockPerTicket(post.Query.YyAbility.YyCunhuozhouzhuanglv, "yy_cunhuozhouzhuanglv")
	ability_set15 = d.ParseStockPerTicket(post.Query.YyAbility.YyLiudongzichanzhouzhuanglv, "yy_liudongzichanzhouzhuanglv")
	ability_set16 = d.ParseStockPerTicket(post.Query.YyAbility.YyZongzichanzhouzhuanglv, "yy_zongzichanzhouzhuanglv")
	ability_set17 = d.ParseStockPerTicket(post.Query.YyAbility.YyGudongquanyizhouzhuanglv, "yy_gudongquanyizhouzhuanglv")

	all_sets := []set.Interface{belong_set, location_set, concept_set,
		per_ticket_set1, per_ticket_set2, per_ticket_set3, per_ticket_set4, per_ticket_set5, per_ticket_set6,
		last_day_set1, last_day_set2, last_day_set3, last_day_set4,
		ability_set1, ability_set2, ability_set3, ability_set4, ability_set5, ability_set6, ability_set7, ability_set8,
		ability_set9, ability_set10, ability_set11, ability_set12, ability_set13, ability_set14, ability_set15, ability_set16, ability_set17}

	used_sets := []set.Interface{}
	for _, i := range all_sets {
		if len(i.List()) > 0 {
			used_sets = append(used_sets, i)
		}
	}
	var coders []interface{}
	if len(used_sets) == 0 {
		coders = nil
	} else if len(used_sets) == 1 {
		coders = used_sets[0].List()
	} else if len(used_sets) == 2 {
		coders = set.Intersection(used_sets[0], used_sets[1]).List()
	} else if len(used_sets) > 2 {
		coders = set.Intersection(used_sets[0], used_sets[1], used_sets[2:]...).List()
	}
	var predicts []dal.Predict
	var total int
	tmp := store.MysqlClient.GetDB().Model(&dal.Predict{}).Where("date = ?", post.Date)
	for _, i := range post.Query.Predicts {
		tmp = tmp.Where("`condition` regexp ?", i)
	}
	tmp = tmp.Where("code IN (?)", coders)
	tmp.Count(&total)

	log.Println(fmt.Sprintf("一共筛选(%d个), 总数(%d个)", len(coders), total))

	if !utils.ContainsString(OrderLimit, post.Order) {
		post.Order = "score"
	}
	tmp.Order(fmt.Sprintf("%s desc", post.Order)).Limit(limit).Offset(offset).Find(&predicts)

	var response []model.PredictListResponse
	for _, i := range predicts {
		var coder dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", i.Code).Find(&coder)
		x := model.PredictListResponse{Name: i.Name, Code: i.Code, Price: i.Price, Percent: i.Percent, Location: coder.Location,
			Form: coder.OrganizationalForm, Belong: coder.Belong, FundCount: i.FundCount, SimuCount: i.SMCount, Conditions: i.Condition, BadConditions: i.BadCondition, Finance: i.Finance,
			Date: i.Date, Score: i.Score}
		response = append(response, x)
	}
	d.Response(c, map[string]interface{}{"result": response, "total": total}, nil)
}

func (d *PredictControl) GetDetail(c *gin.Context) {
	date := c.DefaultQuery("date", "")
	code := c.DefaultQuery("code", "")
	if code == "" || date == "" {
		d.Response(c, nil, errors.New("证券代码/日期为空"))
		return
	}
	var TicketHistory []dal.TicketHistory
	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ? and date <= ?", code, date).Limit(70).Order("date desc").Find(&TicketHistory)

	var TicketHistoryWeekly []dal.TicketHistoryWeekly
	store.MysqlClient.GetDB().Model(&dal.TicketHistoryWeekly{}).Where("code = ? and date <= ?", code, date).Limit(40).Order("date desc").Find(&TicketHistoryWeekly)

	var Stockholder []dal.Stockholder
	store.MysqlClient.GetDB().Model(&dal.Stockholder{}).Where("code = ?", code).Find(&Stockholder)

	var Stock dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", code).Find(&Stock)

	var Predict dal.Predict
	store.MysqlClient.GetDB().Model(&dal.Predict{}).Where("code = ? and date = ?", code, date).Find(&Predict)

	var StockCashFlow []dal.StockCashFlow
	store.MysqlClient.GetDB().Model(&dal.StockCashFlow{}).Where("code = ?", code).Find(&StockCashFlow)

	var StockLiabilities []dal.StockLiabilities
	store.MysqlClient.GetDB().Model(&dal.StockLiabilities{}).Where("code = ?", code).Find(&StockLiabilities)

	var StockProfit []dal.StockProfit
	store.MysqlClient.GetDB().Model(&dal.StockProfit{}).Where("code = ?", code).Find(&StockProfit)

	var PerTickets dal.StockPerTicket
	store.MysqlClient.GetDB().Model(&dal.StockPerTicket{}).Where("code = ?", code).Find(&PerTickets)

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
		return
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
		return
	}
	var FundHoldRank []dal.FundHoldRank
	store.MysqlClient.GetDB().Model(&dal.FundHoldRank{}).Where("fund_code = ?", code).Find(&FundHoldRank)
	d.Response(c, FundHoldRank, nil)
}

func (d *PredictControl) TopHolderHold(c *gin.Context) {
	holder := c.DefaultQuery("holder_name", "")
	if holder == "" {
		d.Response(c, nil, errors.New("查询用户为空"))
		return
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

type Concepts struct {
	Name string `json:"name"`
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

func (d *PredictControl) GetConcepts(c *gin.Context) {
	var x []Concepts
	var result []string
	store.MysqlClient.GetDB().Model(&dal.StockConcept{}).Select("name").Order("name desc").Scan(&x)
	for _, i := range x {
		if i.Name != "" {
			result = append(result, i.Name)
		}
	}
	d.Response(c, result, nil)
}

func (d *PredictControl) GetLabels(c *gin.Context) {
	var x []Concepts
	var result []string
	store.MysqlClient.GetDB().Model(&dal.StockLabels{}).Select("name").Order("name desc").Scan(&x)
	for _, i := range x {
		if i.Name != "" {
			result = append(result, i.Name)
		}
	}
	d.Response(c, result, nil)
}
