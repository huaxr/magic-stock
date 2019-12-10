// @Time:       2019/12/10 下午3:37

package model

import "magic/stock/dal"

type GetPredicts struct {
	Predicts []string `json:"predicts"`
}

type StockDetail struct {
	TicketHistory       []dal.TicketHistory       `json:"ticket_history"`
	TicketHistoryWeekly []dal.TicketHistoryWeekly `json:"ticket_history_weekly"`
	//FundHold            []dal.FundRank            `json:"fund_hold"`   // 持仓机构详情
	Stockholder      []dal.Stockholder      `json:"stockholder"` // 十大流通股东
	Stock            dal.Code               `json:"stock"`       // 股票详情
	Predict          dal.Predict            `json:"predict"`     // 计算结果
	StockCashFlow    []dal.StockCashFlow    `json:"stock_cash_flow"`
	StockLiabilities []dal.StockLiabilities `json:"stock_liabilities"`
	StockProfit      []dal.StockProfit      `json:"stock_profit"`
}

type StockFund struct {
	FundInfo dal.FundRank `json:"fund"`    // 机构信息
	Percent  float64      `json:"percent"` // 机构购买此股占比
}
