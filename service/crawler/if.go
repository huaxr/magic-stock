// @Time:       2019/12/2 下午2:21

package crawler

import (
	"magic/stock/dal"
	"magic/stock/model"

	"github.com/PuerkitoBio/goquery"
)

type CrawlerIF interface {
	//// Deprecate
	//GetFundRanks()
	//// Deprecate 通过上面获取重仓股
	//GetFundHighHold(date string)
	//// 把今日的收盘价 加入到周线均价的表中
	//AddTodayShouToWeek(code, last_week, last_day_to_delete, today string)
	//// 生成周线表
	//GenerateWeekHistory(code string)
	//// 计算周线表百分比
	//CalcPercentTicketWeekly(code string)
	NewDocumentWithProxy(uri string) (*goquery.Document, error)
	NewDocument(url string) (*goquery.Document, error)

	// 获取短期回报基金排行前300名
	GetStockAllFund(code string, proxy bool)

	// 获取上证所有数据代码和名称
	GetAllTicketCode()
	// 股票的所属概念信息 记录到 code 表中
	GetAllTicketCodeConcept(code dal.Code, proxy bool)
	GetAllTicketCodeBelong(code dal.Code, proxy bool)
	// 获取股票的公司简介信息
	GetAllTicketCodeInfo(code dal.Code, proxy bool)
	GetAllTicketCodeInfo2(code dal.Code, proxy bool)
	// 获取公司的公告新闻
	GetPublicNews(code dal.Code, proxy bool)
	GetPublicReports(code dal.Code, proxy bool)

	// 网易api获得十大流通股东（带新进变化趋势
	GetTopStockholder(code, namer string, proxy bool)
	// 获取股票历史记录
	GetSignalTicket(code dal.Code, proxy bool) error
	// 获取股票的 收益表 数据
	GetStockProfit(code string, proxy bool)
	// 获取股票的 资产负债表 数据
	GetStockLiabilities(code string, proxy bool)
	// 获取股票的 现金流量表 数据
	GetStockCashFlow(code string, proxy bool)
	// 获取股票每股的收益情况 运营能力 成长能力
	GetStockPerTicket(code string, proxy bool)
	// 获取分红 配股数据
	GetProfitSharingAndStockOwnership(code string, proxy bool)
	// 获取增发数据
	GetZengFa(code string, proxy bool)
	// 获取今日股价
	GetAllTicketTodayDetail(code, name, today, last_today_str string, proxy bool) error
	// 获取所有子公司以及控股公司记录
	GetSubCompany(code string, proxy bool)

	// 计算返回
	CalcResultWithDefined(params *model.Params) *model.CalcResult
	// 分析
	Analyze(result *model.CalcResult, code, name string)

	// 通过算法获取周线数据
	GetWeekDays(code dal.Code, d1, d2 string)
	// 通过算法获取月线数据
	GetMonthDays(code dal.Code)

	// 只把今日的周线月线加入
	GetWeekDay(code dal.Code, week_begin, today string)
	GetMonthDay(code dal.Code, month_begin string, today string)
}

var CrawlerGlobal CrawlerIF

func init() {
	tmp := new(Crawler)
	CrawlerGlobal = tmp
}

type Crawler struct {
}
