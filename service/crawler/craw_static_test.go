// @Time:       2019/12/2 下午5:40

package crawler

import (
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"strings"
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

// 从股票概念中抽出详细的概念保存在表中
func TestGetConcept(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)
	xx := map[string]bool{}
	for _, i := range code {
		x := strings.Split(strings.TrimRight(i.Concept, ","), ",")
		for _, j := range x {
			xx[j] = true
		}
	}

	for i, _ := range xx {
		if len(i) == 0 {
			continue
		}
		if strings.Contains(i, "概念") {
			c := dal.StockConcept{Name: i}
			store.MysqlClient.GetDB().Save(&c)
		} else {
			c := dal.StockLabels{Name: i}
			store.MysqlClient.GetDB().Save(&c)
		}
	}
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

// 获取股票每股财务指标
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

// 上面获取到了个股的财务指标 以及 各项能力后 对按照大小的区间 对所有股打标
func TestCalcCaiWuForPreTicket(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)
	for _, i := range code {
		CrawlerGlobal.CalcCaiWuForPreTicket(i.Code)
	}
}
