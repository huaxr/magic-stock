// @Time:       2019/12/3 下午1:51

package crawler

import (
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup //定义一个同步等待的组

const (
	today_str    = "2020-01-14"
	last_day_str = "2020-01-13" // 上一个交易日数据 可以计算量比用, 删除昨日周线等
	week_begin   = "2020-01-09"
	month_begin  = "2020-01-01"
)

// 获取今日的所有股票 （自动加入到线上）
func TestGetAllTicketTodayDetail(t *testing.T) {
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
}

// 从日线中获取到周线的数据
func TestGetWeekDay(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)

	for _, i := range code {
		//CrawlerGlobal.GetWeekDay(i, last_week, today_str, last_today_str) // 会删除 last_today_str 的所有数据
		CrawlerGlobal.GetWeekDay(i, week_begin, today_str, "")
	}
}

// 从日线中获取到月线的数据
func TestGetMouthDay(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id < 2").Find(&code)
	for _, i := range code {
		CrawlerGlobal.GetMonthDay(i, month_begin, today_str, "")
	}
}
