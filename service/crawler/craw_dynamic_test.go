// @Time:       2019/12/3 下午1:51

package crawler

import (
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"testing"
	"time"
)

// 得出基金排行并根据这些基金获取持仓股
func TestCrawler_GetFundHighHold(t *testing.T) {
	CrawlerGlobal.GetFundRanks()
	CrawlerGlobal.GetFundHighHold("2019-12-02")
}

// 把今日收盘信息加入周线
func TestAddTodayShouToWeek(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 0").Find(&code)
	for _, i := range code {
		CrawlerGlobal.AddTodayShouToWeek(i.Code, "2019-09-30", "", "2019-10-11") // 再次用的时候把 2019-06-28 全部删掉 用来计算均价用
	}

}

func TestGetAllTicketTodayDetail(t *testing.T) {
	today := "2019-12-09"
	// 超时重试
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id < 2000").Find(&code)
		for _, i := range code {
		RE:
			err := CrawlerGlobal.GetAllTicketTodayDetail(i.Code, i.Name, today, true)
			if err != nil {
				log.Println("爬虫错误， 休眠10秒继续...", i.Name)
				time.Sleep(10 * time.Second)
				goto RE
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2000").Find(&code)
		for _, i := range code {
		RE:
			err := CrawlerGlobal.GetAllTicketTodayDetail(i.Code, i.Name, today, true)
			if err != nil {
				log.Println("爬虫错误， 休眠10秒继续...", i.Name)
				time.Sleep(10 * time.Second)
				goto RE
			}
			time.Sleep(1 * time.Second)
		}
	}()
	select {}
}
