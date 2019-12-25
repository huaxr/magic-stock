// @Time:       2019/12/2 下午2:21

package crawler

import (
	"magic/stock/dal"
	"magic/stock/model"

	"github.com/PuerkitoBio/goquery"
)

type CrawlerIF interface {
	NewDocumentWithProxy(uri string) (*goquery.Document, error)
	NewDocument(url string) (*goquery.Document, error)
	// 把今日的收盘价 加入到周线均价的表中
	AddTodayShouToWeek(code, last_week, last_day_to_delete, today string)
	// 获取短期回报基金排行前300名
	GetFundRanks()
	// 通过上面获取重仓股
	GetFundHighHold(date string)

	// 获取上证所有数据代码和名称
	GetAllTicketCode()
	// 股票的所属概念信息 记录到 code 表中
	GetAllTicketCodeConcept(code dal.Code, proxy bool)
	GetAllTicketCodeBelong(code dal.Code, proxy bool)
	// 获取股票的公司简介信息
	GetAllTicketCodeInfo(code dal.Code, proxy bool)
	// 网易api获得十大流通股东（带新进变化趋势
	GetTopStockholder(code, namer string, proxy bool)
	// 获取股票历史记录
	GetSignalTicket(code, name string, proxy bool) error
	// 获取股票的 收益表 数据
	GetStockProfit(code string, proxy bool)
	// 获取股票的 资产负债表 数据
	GetStockLiabilities(code string, proxy bool)
	// 获取股票的 现金流量表 数据
	GetStockCashFlow(code string, proxy bool)
	// 获取股票每股的收益情况 运营能力 成长能力
	GetStockPerTicket(code string, proxy bool)
	// 通过前几名来判断股票的价值 能力
	CalcCaiWuForPreTicket(code string)
	// 生成周线表
	GenerateWeekHistory(code string)
	// 计算周线表百分比
	CalcPercentTicketWeekly(code string)
	// 获取今日股价
	GetAllTicketTodayDetail(code, name, today string, proxy bool) error

	// 计算返回
	CalcResultWithDefined(params *model.Params) *model.CalcResult
	// 分析
	// 近几周百分比小于less
	WeeklyPercentLimited(result *model.CalcResult, recent_num int, limit_percent float64, typ string) bool
	Analyze(result *model.CalcResult, code, name string)
}

var CrawlerGlobal CrawlerIF

func init() {
	tmp := new(Crawler)
	CrawlerGlobal = tmp
}

type Crawler struct {
}
