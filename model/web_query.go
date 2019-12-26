// @Time:       2019/12/10 下午3:37

package model

import "magic/stock/dal"

type QueryPerTicket struct {
	// 每股收益调整后  指扣除与主营业务无关的一次性损益后的净利润除以总股本得出的每股收益。
	Shouyiafter map[string]float64 `json:"shouyiafter"`
	// 使用加权平均法来计算每股收益，这样才可以更准确、更合理地反映公司客观的盈利能力。
	Jiaquanshouyi map[string]float64 `json:"jiaquanshouyi"`
	// 每股净资产_调整后(元) 这一指标反映每股股票所拥有的资产现值。每股净资产越高，股东拥有的每股资产价值越多；
	Jinzichanafter map[string]float64 `json:"jinzichanafter"`
	// 每股经营性现金流(元)范围  即每股经营活动产生的现金流量净额
	Jingyingxianjinliu map[string]float64 `json:"jingyingxianjinliu"`
	// 每股资本公积金(元)范围  资本公积金是指从公司的利润以外的收入中提取的一种公积金。其主要来源有股票溢价收入，财产重估增值，以及接受捐赠资产等。每股资本公积金=资本公积金/总股本
	Gubengongjijin map[string]float64 `json:"gubengongjijin"`
	// 每股未分配利润(元)范围 (1)每股未分配利润越多，不仅表明该公司盈利能力强，也意味着该公司未来分红、送股的能力强、概率比较大。
	//(2)一般而言，如果某公司的每股未分配利润超过一元，该公司就具有每10股送10股或每股派现1元红利的能力。
	//(3)每股未分配利润较多的上市公司，往往被各类投资者青睐，因为该类公司盈利和分红能力强，投资回报高。
	Weifenpeilirun map[string]float64 `json:"weifenpeilirun"`
}

type LastDayRange struct {
	LastPrice        map[string]float64 `json:"last_price"`        // 昨日收盘股价区间
	LastPercent      map[string]float64 `json:"last_percent"`      // 昨日涨幅区间 %
	LastAmplitude    map[string]float64 `json:"last_amplitude"`    // 昨日振幅区间
	LastTurnoverrate map[string]float64 `json:"last_turnoverrate"` // 昨日换手率区间
}

// 盈利能力区间
type YlAbility struct {
	// 总资产利润率(%)  总资产利润率=利润总额/资产平均总额×100% 可用来说明企业运用其全部资产获取利润的能力。
	YlZongzichanlirunlv map[string]float64 `json:"yl_zongzichanlirunlv"`
	// 主营业务利润率(%) 主营利润率，是公司主业所产生的利润率。比如公司主业是房地产，那么经营房地产所产生的利润，与主营业务收入的比率，就是主营利润率。
	YlZhuyingyewulirunlv map[string]float64 `json:"yl_zhuyingyewulirunlv"`
	// 总资产净利润率(%) // 又称总资产收益率，是企业净利润总额与企业资产平均总额的比率，即过去所说的资金利润率
	YlZongzichanjinglirunlv map[string]float64 `json:"yl_zongzichanjinglirunlv"`
	// 营业利润率(%)  // 营业利润率是企业付清一切帐项后剩下的金额称为利润。在会计学上，利润可分为毛利
	YlYingyelirunlv map[string]float64 `json:"yl_yingyelirunlv"`
	// 销售净利率(%) 是净利润占销售收入的百分比。 该指标反映每一元销售收入带来的净利润的多少，表示销售收入的收益水平。
	YlXiaoshoujinglilv map[string]float64 `json:"yl_xiaoshoujinglilv"`
	// 股本报酬率(%)  股本报酬率是指公司税后利润与其股本的比率，表明公司股本总额中平均每百元股本所获得的纯利润。
	YlGubenbaochoulv map[string]float64 `json:"yl_gubenbaochoulv"`
	// 净资产报酬率(%) // 该指标反映股东权益的收益水平，用以衡量公司运用自有资本的效率。指标值越高，说明投资带来的收益越高。该指标体现了自有资本获得净收益的能力。
	YlJingzichanbaochoulv map[string]float64 `json:"yl_jingzichanbaochoulv"`
	// 资产报酬率(%) 用以评价企业运用全部资产的总体获利能力，是评价企业资产运营效益的重要指标。
	YlZichanbaochoulv map[string]float64 `json:"yl_zichanbaochoulv"`
}

// 成长能力
type CzAbility struct {
	// 主营业务收入增长率(%)
	CzZhuyingyewushouruzengzhanglv map[string]float64 `json:"cz_zhuyingyewushouruzengzhanglv"`
	// 净利润增长率(%)  净利润增长率是指企业当期净利润比上期净利润的增长幅度，指标值越大代表企业盈利能力越强。
	CzJinglirunzengzhanglv map[string]float64 `json:"cz_jinglirunzengzhanglv"`
	// 净资产增长率(%)
	CzJingzichanzengzhanglv map[string]float64 `json:"cz_jingzichanzengzhanglv"`
	// 总资产增长率(%)
	CzZongzichanzengzhanglv map[string]float64 `json:"cz_zongzichanzengzhanglv"`
}

// 运营能力
type YyAbility struct {
	// 应收账款周转率(次) 应收账款周转率越高越好，应收账示周转率高，表明收账迅速，账龄较短；
	YyYingshouzhangkuanzhouzhuanlv map[string]float64 `json:"yy_yingshouzhangkuanzhouzhuanlv"`
	// 存货周转率(次) 存货周转率越高，表明企业存货资产变现能力越强，存货及占用在存货上的资金周转速度越快。
	YyCunhuozhouzhuanglv map[string]float64 `json:"yy_cunhuozhouzhuanglv"`
	// 流动资产周转率(次) 该指标越高，说明企业流动资产的利用效率越好。
	YyLiudongzichanzhouzhuanglv map[string]float64 `json:"yy_liudongzichanzhouzhuanglv"`
	// 总资产周转率(次) 总资产周转率越高，说明企业销售能力越强,资产投资的效益越好
	YyZongzichanzhouzhuanglv map[string]float64 `json:"yy_zongzichanzhouzhuanglv"`
	// 股东权益周转率(次)  指标说明公司运用所有制的资产的效率。 该比率越高，表明所有者资产的运用效率高，营运能力强
	YyGudongquanyizhouzhuanglv map[string]float64 `json:"yy_gudongquanyizhouzhuanglv"`
}

type Query struct {
	Predicts  []string `json:"predicts"`  // 预测的打标
	Belongs   []string `json:"belongs"`   // 所属行业
	Locations []string `json:"locations"` // 所属地区
	Concepts  []string `json:"concepts"`  // 所属概念 ex 腾讯概念
	Labels    []string `json:"labels"`    // 标签

	PerTickets   QueryPerTicket `json:"per_tickets"`    // 高级条件区间
	LastDayRange LastDayRange   `json:"last_day_range"` // 昨日收盘情况区间
	YlAbility    YlAbility      `json:"yl_ability"`     // 盈利能力区间
	CzAbility    CzAbility      `json:"cz_ability"`     // 成长能力区间
	YyAbility    YyAbility      `json:"yy_ability"`     // 运营能力区间
}

type GetPredicts struct {
	Query Query  `json:"query"`
	Date  string `json:"date"`
	Order string `json:"order"`
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

	PerTicket dal.StockPerTicket `json:"per_ticket"` // 每股指标
}

type StockFund struct {
	FundInfo dal.FundRank `json:"fund"`    // 机构信息
	Percent  float64      `json:"percent"` // 机构购买此股占比
}
