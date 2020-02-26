// @Time:       2019/12/2 下午5:16

package model

type Detail struct {
	FundCode   string  `json:"fundCode"`
	FundSName  string  `json:"fundSName"`
	LastMonth  float64 `json:"pct1Mon"`
	LastWeek   float64 `json:"pct1Week"`
	Last3Month float64 `json:"pct3Mon"`
	Last6Month float64 `json:"pct6Mon"`
	SinceBase  float64 `json:"pctBase"` // 成立以来涨跌幅
	ThisYear   float64 `json:"pctTYear"`
	LastYear   float64 `json:"pct1Year"`
	Last2Year  float64 `json:"pct2Year"`
	Last3Year  float64 `json:"pct3Year"`
	Time       int64   `json:"tradeDate"`
}

type Result struct {
	Limit     int       `json:"limit"`
	Curpage   int       `json:"curpage"`
	Total     int       `json:"total"`
	TradeDate int64     `json:"tradeDate"`
	List      []*Detail `json:"list"`
	Type      int       `json:"type"`
}

type Hold struct {
	Code    string  `json:"stockcode"`  // 证券代码
	Percent float64 `json:"totValProp"` // 持仓比例
}

// 预测结果
type RecentDailyData struct {
	RecentCount        []float64 // 最近成交量
	RecentClose        []float64 // 最近收盘价
	RecentOpen         []float64 // 最近的开盘价
	RecentHigh         []float64 // 最近的最高价
	RecentLow          []float64 // 最近的最低价
	RecentPercent      []float64 // 最近涨跌幅
	RecentAmplitude    []float64 // 最近振幅
	RecentTurnoverRate []float64 // 最近换手率
	//RecentNetFlow      []float64 // 最近资金净流入
	//RecentMainNetFlow  []float64 // 最近主力资金净流入
	CurrDate       string // 今天日期
	CurrTotalMoney float64
}

type RecentWeeklyData struct {
	RecentCountWeek        []float64 // 最近成交量
	RecentCloseWeek        []float64 // 最近收盘价
	RecentOpenWeek         []float64 // 最近的开盘价
	RecentHighWeek         []float64 // 最近的最高价
	RecentLowWeek          []float64 // 最近的最低价
	RecentPercentWeek      []float64 // 最近涨跌幅
	RecentAmplitudeWeek    []float64 // 最近振幅
	RecentTurnoverRateWeek []float64 // 最近换手率
}

type RecentMonthData struct {
	RecentCountMonth        []float64 // 最近成交量
	RecentCloseMonth        []float64 // 最近收盘价
	RecentOpenMonth         []float64 // 最近的开盘价
	RecentHighMonth         []float64 // 最近的最高价
	RecentLowMonth          []float64 // 最近的最低价
	RecentPercentMonth      []float64 // 最近涨跌幅
	RecentAmplitudeMonth    []float64 // 最近振幅
	RecentTurnoverRateMonth []float64 // 最近换手率
}

type RecentAverage struct {
	AveDailyPrice1 []float64
	AveDailyPrice2 []float64
	AveDailyPrice3 []float64
	AveDailyPrice4 []float64
	AveCount1      []float64
	AveCount2      []float64
}

type RecentAverageWeekly struct {
	AveWeeklyPrice1 []float64
	AveWeeklyPrice2 []float64
	AveWeeklyPrice3 []float64
	AveWeeklyPrice4 []float64

	AveCountWeekly1 []float64
	AveCountWeekly2 []float64
}

type RecentAverageMonth struct {
	AveMonthPrice1 []float64
	AveMonthPrice2 []float64
	AveMonthPrice3 []float64
	AveMonthPrice4 []float64

	AveCountMonth1 []float64
	AveCountMonth2 []float64
}

type Params struct {
	Code          string
	Date          string
	Offset        int
	AveragePrice1 int // 均线1 5
	AveragePrice2 int // 均线2 10
	AveragePrice3 int // 均线2 30
	AveragePrice4 int // 均线3 60

	AverageCount1 int // 量均线1 10
	AverageCount2 int // 量均线2 40
}

type CalcResult struct {
	*RecentDailyData
	*RecentWeeklyData
	*RecentMonthData
	*RecentAverage
	*RecentAverageWeekly
	*RecentAverageMonth
}

type PredictListResponse struct {
	Name            string  `json:"name"`
	Code            string  `json:"code"`
	Price           float64 `json:"price"`            // 当前价格
	Percent         float64 `json:"percent"`          // 涨幅
	Location        string  `json:"location"`         // 地区
	Form            string  `json:"form"`             // 组织形式
	Belong          string  `json:"belong"`           // 所属行业
	FundCount       int     `json:"fund_count"`       // 基金数量
	SimuCount       int     `json:"simu_count"`       // 私募数量
	FenghongCount   int     `json:"fenghong_count"`   // 分红次数
	SongguCount     int     `json:"songgu_count"`     // 送股次数
	ZhuangzengCount int     `json:"zhuangzeng_count"` // 转增次数
	PeiguCount      int     `json:"peigu_count"`      // 配股次数
	ZengfaCount     int     `json:"zengfa_count"`     // 增发次数
	SubcompCount    int     `json:"subcomp_count"`    // 控股数量
	Conditions      string  `json:"conditions"`       // 优点
	BadConditions   string  `json:"bad_conditions"`   // 缺点
	//Finance         string  `json:"finance"`          // 财务指标
	Kai          float64 `json:"kai"`
	High         float64 `json:"high"`
	Low          float64 `json:"low"`
	Amplitude    float64 `json:"amplitude"`     // 振幅
	TurnoverRate float64 `json:"turnover_rate"` // 换手率
	NumberRate   float64 `json:"number_rate"`   // 对昨量比

	Business string `json:"business"` // 业务
	Score    int    `json:"score"`    // 得分
	Tape     string `json:"tape"`     // 盘口
	Rong     bool   `json:"rong"`     // 是否融资融券股
	Date     string `json:"date"`
}
