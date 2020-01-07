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
	RecentNetFlow      []float64 // 最近资金净流入
	RecentMainNetFlow  []float64 // 最近主力资金净流入
	CurrDate           string    // 今天日期
	CurrTotalMoney     float64
}

type RecentWeeklyData struct {
	RecentWeeklyClose   []float64
	RecentWeeklyPercent []float64
}

type RecentAverage struct {
	AveDailyPrice1 []float64
	AveDailyPrice2 []float64
	AveDailyPrice3 []float64
	AveDailyPrice4 []float64
	AveDailyPrice5 []float64
	//AveWeeklyPrice1 []float64
	//AveWeeklyPrice2 []float64
	AveCount1 []float64
	AveCount2 []float64
}

type Params struct {
	Code          string
	Date          string
	Offset        int
	AveragePrice1 int // 均线1 5
	AveragePrice2 int // 均线2 10
	AveragePrice3 int // 均线2 15
	AveragePrice4 int // 均线3 30
	AveragePrice5 int // 均线4 60

	AverageCount1 int // 量均线1 10
	AverageCount2 int // 量均线2 40
}

type CalcResult struct {
	*RecentDailyData
	//*RecentWeeklyData
	*RecentAverage
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
	Finance         string  `json:"finance"`          // 财务指标
	Date            string  `json:"date"`
	Score           int     `json:"score"` // 得分
}
