// @Time:       2019/12/2 下午5:40

package crawler

import (
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"testing"
	"time"
)

// 注意: 这里的static要注意使用， 不要重复多次使用

// 获取所有股票证券代码
func TestCrawler(t *testing.T) {
	CrawlerGlobal.GetAllTicketCode()
}

// 获取所有股票概念信息
func TestGetAllTicketCodeConcept(T *testing.T) {
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 1919 and id < 2000").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetAllTicketCodeBelong(i, false)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 3684").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetAllTicketCodeBelong(i, true)
			time.Sleep(2 * time.Second)
		}
	}()

	select {}
}

// 获取所有股票简介信息
func TestCrawler_GetAllTicketCodeInfo(t *testing.T) {
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2000 and id <= 2005").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetAllTicketCodeInfo(i, false)
		}
	}()

	//go func() {
	//	var code []dal.Code
	//	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2000").Find(&code)
	//	for _, i := range code {
	//		CrawlerGlobal.GetAllTicketCodeInfo(i, true)
	//	}
	//}()

	select {}

}

// 获取十大流通股东（done）
func TestCrawler_GetTopStockholder(t *testing.T) {
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 0 and id < 2000").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetTopStockholder(i.Code, i.Name, false)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2000").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetTopStockholder(i.Code, i.Name, true)
			time.Sleep(2 * time.Second)
		}
	}()
	select {}
}

// 获取股票历史记录
func TestGetSignalTicket(T *testing.T) {
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 0").Find(&code)
		for _, i := range code {
		RE:
			err := CrawlerGlobal.GetSignalTicket(i.Code, i.Name, true)
			if err != nil {
				log.Println("爬虫错误， 休眠10秒继续...", i.Name)
				time.Sleep(10 * time.Second)
				goto RE
			}
			time.Sleep(2 * time.Second)
		}
	}()
	select {}
}

// 获取股票历史记录--资金流入流出数据
func TestGetSignalTicketFlow(T *testing.T) {
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 686 and id < 2000").Find(&code)
		for _, i := range code {
		RE:
			err := CrawlerGlobal.GetSignalTicketFlow(i, true)
			if err != nil {
				log.Println("爬虫错误， 休眠10秒继续...", i.Name)
				time.Sleep(10 * time.Second)
				goto RE
			}
		}
	}()
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2776").Find(&code)
		for _, i := range code {
		RE:
			err := CrawlerGlobal.GetSignalTicketFlow(i, false)
			if err != nil {
				log.Println("爬虫错误， 休眠10秒继续...", i.Name)
				time.Sleep(10 * time.Second)
				goto RE
			}
		}
	}()
	select {}
}

// 从股票历史记录计算出周线数据
func TestCalcWeekPrice(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)
	for _, i := range code {
		CrawlerGlobal.GenerateWeekHistory(i.Code)
	}
	for _, i := range code {
		CrawlerGlobal.CalcPercentTicketWeekly(i.Code)
	}
}

// 获取基本面信息
func TestGetStockProfit(t *testing.T) {
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id < 2000 and id >= 735").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetStockProfit(i.Code, true)
			time.Sleep(1 * time.Second)
			CrawlerGlobal.GetStockLiabilities(i.Code, true)
			time.Sleep(1 * time.Second)
			CrawlerGlobal.GetStockCashFlow(i.Code, true)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 3163").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetStockProfit(i.Code, false)
			time.Sleep(1 * time.Second)
			CrawlerGlobal.GetStockLiabilities(i.Code, false)
			time.Sleep(1 * time.Second)
			CrawlerGlobal.GetStockCashFlow(i.Code, false)
			time.Sleep(1 * time.Second)
		}
	}()

	select {}
}

//func TestSync(t *testing.T) {
//	var code []dal.TicketHistory
//	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("date = ?", "2019-12-10").Find(&code)
//	for _, i := range code {
//		store.MysqlClient.GetTmpDb().Save(&i)
//	}
//}
