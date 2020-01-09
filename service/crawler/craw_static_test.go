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

// 获取所有股票概念信息 所属行业 GetAllTicketCodeBelong
// 首先更新数据表 update magic_stock_code set concept = null   新浪api
func TestGetAllTicketCodeConcept(T *testing.T) {
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id < 2000").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetAllTicketCodeConcept(i, false)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2000").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetAllTicketCodeConcept(i, true)
			time.Sleep(2 * time.Second)
		}
	}()

	select {}
}

// 获取所有股票简介信息
func TestCrawler_GetAllTicketCodeInfo(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("institutional_type = ?", "").Find(&code)
	go func() {
		for _, i := range code[0 : (len(code)-1)/2] {
			CrawlerGlobal.GetAllTicketCodeInfo2(i, false)
			time.Sleep(2 * time.Second)
		}
	}()
	go func() {
		for _, i := range code[(len(code)-1)/2 : (len(code) - 1)] {
			CrawlerGlobal.GetAllTicketCodeInfo2(i, true)
			time.Sleep(2 * time.Second)
		}
	}()
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

	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id > 562 and id <= 2000").Find(&code)
	for _, i := range code {
		log.Println("正在爬取", i.ID, i.Code, i.Name)
		CrawlerGlobal.GetSignalTicket(i, false)
	}

	//go func() {
	//	var code []dal.Code
	//	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id > 1000 and id <= 2000").Find(&code)
	//	for _, i := range code {
	//		log.Println("正在爬取", i.ID, i.Code, i.Name)
	//		CrawlerGlobal.GetSignalTicket(i, false)
	//	}
	//}()
	select {}
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

// 获取股票每股财务指标  // 运用能力 成长能力
func TestCrawler_GetStockPerTicket(t *testing.T) {
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id < 2000").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetStockPerTicket(i.Code, true)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2000").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetStockPerTicket(i.Code, false)
			time.Sleep(1 * time.Second)
		}
	}()
	select {}
}

// 获取分红配股
func TestGetProfitSharingAndStockOwnership(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id > 0").Find(&code)
	for _, i := range code {
		CrawlerGlobal.GetProfitSharingAndStockOwnership(i.Code, false)
		time.Sleep(5 * time.Second)
	}
}

// 获取增发
func TestGetZengFa(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2479").Find(&code)
	for _, i := range code {
		CrawlerGlobal.GetZengFa(i.Code, false)
		time.Sleep(1 * time.Second)
	}
}

// 获取控股公司记录
func TestCrawler_GetSubCompany(t *testing.T) {
	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id < 2000").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetSubCompany(i.Code, true)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		var code []dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id >= 2000").Find(&code)
		for _, i := range code {
			CrawlerGlobal.GetSubCompany(i.Code, false)
			time.Sleep(2 * time.Second)
		}
	}()
	select {}
}
