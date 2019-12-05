// @Contact:    huaxinrui
// @Time:       2019/7/5 下午5:27

package dynamic

import (
	"fmt"
	d "huaxinrui/stock/dao"
	"testing"
)

func Init() {
	orm := d.ORM{Store: "root:root@(127.0.0.1:3316)/stock?charset=utf8mb4&parseTime=True&loc=Local"}
	orm.InitDal(&d.Code{}, &d.TicketHistory{}, &d.Predict{}, &d.PredictDebug{}, &d.TicketHistoryWeekly{}, &d.Stockholder{}, &d.FundRank{}, &d.FundHoldRank{})
}

// 获取符合条件的股票数据
func TestGetTicketByConditions(t *testing.T) {
	Init()
	var code []d.Code
	d.Backend.DB.Model(&d.Code{}).Find(&code)

	for _, i := range code {
		//GetTicketByConditions(i.Code)
		GetTicketPredict(i.Code, i.Name, 0, true) // offset 表示获得上几个交易日的数据, 1 表示上个交易日
	}
}

// 每天收盘执行一次, 收集所有股票的当天股价和成交量  http://hq.sinajs.cn/list=sz000001
func TestGetAllTicketTodayDetail(t *testing.T) {
	Init()
	var code []d.Code
	d.Backend.DB.Model(&d.Code{}).Where("id >= 0").Find(&code)
	for _, i := range code {
		GetAllTicketTodayDetail(i.Code, i.Name, "2019-10-11", "2019-10-14") // last_day 用来获取昨日的收盘价， 方便计算出今日的涨跌幅
	}
}

// 计算预测结果并把结果保存到 predict 记录
func TestCalcThePredictResult(t *testing.T) {
	Init()
	var p2 []d.PredictDebug
	var p []d.PredictDebug
	last_day, today := "2019-09-23", "2019-09-24"

	d.Backend.DB.Model(&d.PredictDebug{}).Where("date = ?", last_day).Find(&p) // 获取前交易日的预测
	for _, i := range p {
		var r d.TicketHistory
		d.Backend.DB.Model(&d.TicketHistory{}).Where("date = ? and code = ?", today, i.Code).Find(&r) // 将今天实际的percent更新
		i.RealPercent = r.Percent
		d.Backend.DB.Save(&i)
	}

	d.Backend.DB.Model(&d.PredictDebug{}).Where("date = ?", last_day).Find(&p2)
	for _, i := range p2 {
		x := d.Predict{Code: i.Code, Name: i.Name, Date: i.Date, Condition: i.Condition, RealPercent: i.RealPercent, Attrs: i.Attrs, Holder: i.Holder, FundCount: i.FundCount}
		d.Backend.DB.Save(&x)
	}
}

// 把今日的收盘价 加入到周线均价的表中
func TestAddTodayShouToWeek(t *testing.T) {
	Init()
	var code []d.Code
	d.Backend.DB.Model(&d.Code{}).Where("id >= 0").Find(&code)
	for _, i := range code {
		AddTodayShouToWeek(i.Code, "2019-09-30", "", "2019-10-11") // 再次用的时候把 2019-06-28 全部删掉 用来计算均价用
	}
}

func TestGetRank(t *testing.T) {
	Init()
	GetRanks()
}

func TestGetHigh(t *testing.T) {
	Init()
	date := []string{"2019-10-10", "2019-10-09", "2019-09-30"}
	for _, i := range date {
		GetHigh(i)
	}
}

// 获取今日所有股票前4个小时的数据
func TestGetAllTicketHourK(t *testing.T) {
	Init()
	var code []d.Code
	d.Backend.DB.Model(&d.Code{}).Where("id >= 0").Find(&code)
	for _, i := range code {
		GetAllTicketHourK(i.Code, i.Name, "2019-07-12")
	}
}

// 获取分价
func TestGetPriceDivide(t *testing.T) {
	//for i := -900; i<=0; i += 30 {
	//	GetPrice("601019", i, i+30, 1500)
	//	i += 1
	//	time.Sleep(3 * time.Second)
	//}
	GetPriceDivide("601019", -55, 0, 1000)
}

// 龙虎版 机构top30
func TestGetDepartmentLongHuRank(t *testing.T) {
	Init()
	GetDepartmentLongHuRank()
}

// 龙虎版 最近五天 股票
func TestGetTicketLongHuRank(t *testing.T) {
	Init()
	GetTicketLongHuRank()
}

// 龙虎版 今日 详情
func TestGetTicketLongHuRandDetail(t *testing.T) {
	Init()
	var ticket []d.LongHuRankTicket
	d.Backend.DB.Model(&d.LongHuRankTicket{}).Find(&ticket)
	for _, i := range ticket {
		GetTicketLongHuRandDetail(i.Code, i.Name, "2019-07-09")
	}
}

// 每天收盘执行一次 计算当天的MACD
func TestCalcMacdEveryDay(t *testing.T) {
	Init()
	var code []d.Code
	d.Backend.DB.Model(&d.Code{}).Where("id >= 0").Find(&code)
	for _, i := range code {
		CalcMacdEveryDay(i.Code, "2019-06-25", "2019-06-26")
	}
}

func TestX(t *testing.T) {
	Init()
	//var history []d.TicketHistoryTmp
	//d.Backend.DB.Model(&d.TicketHistoryTmp{}).Where("id >= 0").Find(&history)
	//for _, i := range history{
	//
	//	x := d.TicketHistory{Code:i.Code,
	//		Name:i.Name,
	//		Date:i.Date,
	//		Kai:i.Kai, High:i.High, Shou:i.Shou, Low:i.Low, TotalCount:i.TotalCount, TotalMoney:i.TotalMoney,
	//		Percent:i.Percent}
	//	d.Backend.DB.Save(&x)
	//}
	//var tf d.Code
	type res struct {
		Code string
		Name string
	}
	var xx res
	d.Backend.DB.Table("code").Where("code = ?", "000001").Select("code, name").Find(&xx)
	//var xx []map[string]interface{}
	//d.Backend.DB.ScanRows(x, &xx)
	fmt.Println(xx)
}
