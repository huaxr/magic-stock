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

var wg sync.WaitGroup             //定义一个同步等待的组
var last_today_str = "2020-01-02" // 可以计算量比用
var today_str = "2020-01-03"

// 得出基金排行并根据这些基金获取持仓股
func TestCrawler_GetFundHighHold(t *testing.T) {
	CrawlerGlobal.GetFundRanks()
	CrawlerGlobal.GetFundHighHold(today_str)
}

// 获取每只股票的基金持仓情况
func TestGetStockAllFund(t *testing.T) {
	//go func() {
	//	var code []dal.Code
	//	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 1538 and id < 2000").Find(&code)
	//	for _, i := range code {
	//		CrawlerGlobal.GetStockAllFund(i.Code, false)
	//	}
	//}()

	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 3357").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetStockAllFund(i.Code, true)
		}
	}()
	select {}
}

// 把今日收盘信息加入周线 （自动加入到线上）
func TestAddTodayShouToWeek(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 0").Find(&code)
	for _, i := range code {
		CrawlerGlobal.AddTodayShouToWeek(i.Code, "2019-12-20", "", "2019-12-27") // 再次用的时候把 2019-06-28 全部删掉 用来计算均价用
	}
}

// 获取今日的所有股票 （自动加入到线上）
func TestGetAllTicketTodayDetail(t *testing.T) {
	today := today_str
	wg.Add(2)
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id < 2000").Find(&code)
		for _, i := range code {
		RE:
			err := CrawlerGlobal.GetAllTicketTodayDetail(i.Code, i.Name, today, last_today_str, true)
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
			err := CrawlerGlobal.GetAllTicketTodayDetail(i.Code, i.Name, today, last_today_str, true)
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
