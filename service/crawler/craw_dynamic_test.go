// @Time:       2019/12/3 下午1:51

package crawler

import (
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/utils"
	"sync"
	"testing"
	"time"
)

var wg, wg2 sync.WaitGroup //定义一个同步等待的组

var (
	today_str    = "2020-03-20"
	last_day_str = "2020-03-19" // 上一个交易日数据 可以计算量比用

	delete_day_week  = "2020-03-19" // 要删除的周线日线 当为周一的时候置空
	delete_day_month = "2020-03-19" // 当为一月的开始时候置空

	week_begin  = "2020-03-16" // 本周的开始 周一
	month_begin = "2020-03-02" // 本月的开始 一号 基本可不变
)

// 获取今日的所有股票 周 月线， 分析结果并自动加入线上
func TestGetAllTicketTodayDetail(t *testing.T) {
	start := time.Now()
	today := today_str
	wg.Add(2)
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id < 2000").Find(&code)
		for _, i := range code {
		RE:
			err := CrawlerGlobal.GetAllTicketTodayDetail(i.Code, i.Name, today, last_day_str, false)
			if err != nil {
				log.Println("爬虫错误， 休眠10秒继续...", i.Name)
				time.Sleep(10 * time.Second)
				goto RE
			}
		}
		wg.Done()
	}()

	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2000").Find(&code)
		for _, i := range code {
		RE:
			err := CrawlerGlobal.GetAllTicketTodayDetail(i.Code, i.Name, today, last_day_str, true)
			if err != nil {
				log.Println("爬虫错误， 休眠10秒继续...", i.Name)
				time.Sleep(10 * time.Second)
				goto RE
			}
		}
		wg.Done()
	}()
	wg.Wait()

	wg2.Add(2)
	// 抽出周 月线
	go func() {
		store.MysqlClient.GetDB().Exec("delete from magic_stock_history_week where date = ?", delete_day_week)
		if utils.TellEnv() == "loc" {
			store.MysqlClient.GetOnlineDB().Exec("delete from magic_stock_history_week where date = ?", delete_day_week)
		}

		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)

		for _, i := range code {
			//CrawlerGlobal.GetWeekDay(i, last_week, today_str, last_today_str) // 会删除 last_today_str 的所有数据
			CrawlerGlobal.GetWeekDay(i, week_begin, today_str)
		}
		wg2.Done()
	}()

	go func() {
		store.MysqlClient.GetDB().Exec("delete from magic_stock_history_month where date = ?", delete_day_month)
		if utils.TellEnv() == "loc" {
			store.MysqlClient.GetOnlineDB().Exec("delete from magic_stock_history_month where date = ?", delete_day_month)
		}
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetMonthDay(i, month_begin, today_str)
		}
		wg2.Done()
	}()
	wg2.Wait()

	// 计算分析数据
	GetData()
	end := time.Now()
	log.Println("一共耗时（s）:", end.Sub(start).Seconds())
}

// 获取具体日期的分析结果
func GetData() {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)
	for _, i := range code {
		x := &model.Params{i.Code, today_str, 0, 5, 10, 30, 60, 10, 40}
		y := CrawlerGlobal.CalcResultWithDefined(x)
		if y == nil {
			continue
		}
		CrawlerGlobal.Analyze(y, i.Code, i.Name)
	}
}

func TestGetData(t *testing.T) {
	GetData()
}

// 从日线中获取到周线的数据
func GetWeekDay(wg sync.WaitGroup) {
	// 需要删除的昨日数据 注意 如果为周线周五或者节假日收盘 请注释
	store.MysqlClient.GetDB().Exec("delete from magic_stock_history_week where date = ?", last_day_str)
	if utils.TellEnv() == "loc" {
		store.MysqlClient.GetOnlineDB().Exec("delete from magic_stock_history_week where date = ?", last_day_str)
	}

	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)

	for _, i := range code {
		//CrawlerGlobal.GetWeekDay(i, last_week, today_str, last_today_str) // 会删除 last_today_str 的所有数据
		CrawlerGlobal.GetWeekDay(i, week_begin, today_str)
	}
	defer wg.Done()
}

// 从日线中获取到月线的数据
func GetMouthDay(wg sync.WaitGroup) {
	store.MysqlClient.GetDB().Exec("delete from magic_stock_history_month where date = ?", last_day_str)
	if utils.TellEnv() == "loc" {
		store.MysqlClient.GetOnlineDB().Exec("delete from magic_stock_history_month where date = ?", last_day_str)
	}
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)
	for _, i := range code {
		CrawlerGlobal.GetMonthDay(i, month_begin, today_str)
	}
	defer wg.Done()
}


// 获取融资融券数据
func TestGetSecuritiesMarginTrading(t *testing.T) {
	dates := []string{"2020-03-05","2020-03-06","2020-03-09","2020-03-10"}
	for _, i := range dates {
		CrawlerGlobal.GetRongStock(i, false)
	}
}
