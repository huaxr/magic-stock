// @Time:       2019/12/10 下午3:37

package model

import "magic/stock/dal"

type Query struct {
	Predicts           []string `json:"predicts"`            // 预测的打标
	Belongs            []string `json:"belongs"`             // 所属行业
	Locations          []string `json:"locations"`           // 所属地区
	OrganizationalForm []string `json:"organizational_form"` // 组织类型
	Concepts           []string `json:"concepts"`            // 所属概念 ex 腾讯概念
	Labels             []string `json:"labels"`              // 标签
}

type GetPredicts struct {
	Query Query  `json:"query"`
	Date  string `json:"date"`
	Save  bool   `json:"save"` // 查询并保存条件
}

type EditPredicts struct {
	Query Query  `json:"query"`
	Id    int    `json:"id"`
	Name  string `json:"name"`
}

type DeletePredicts struct {
	Id int `json:"id"`
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
