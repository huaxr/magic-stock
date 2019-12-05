package dao

import (
	"github.com/jinzhu/gorm"
)


type Stock struct {
	gorm.Model
	Code string  // models.CharField(verbose_name="股票", default="", max_length=20)

	Bone float64 //= models.FloatField(verbose_name="买一档", default=0)
	Btwo float64 //= models.FloatField(verbose_name="买二档", default=0)
	Bthr float64 //= models.FloatField(verbose_name="买三档", default=0)
	Bfor float64 //= models.FloatField(verbose_name="买四档", default=0)
	Bfir float64 //= models.FloatField(verbose_name="买五档", default=0)

	Sone float64 //= models.FloatField(verbose_name="卖一档", default=0)
	Stwo float64 // = models.FloatField(verbose_name="卖二档", default=0)
	Sthr float64 //= models.FloatField(verbose_name="卖三档", default=0)
	Sfor float64 //= models.FloatField(verbose_name="卖四档", default=0)
	Sfir float64 //= models.FloatField(verbose_name="卖五档", default=0)

	Kai float64 //= models.FloatField(verbose_name="开盘价", default=0)
	Shou float64 //= models.FloatField(verbose_name="昨收盘", default=0)
	Cur float64 //= models.FloatField(verbose_name="当前价", default=0)

	B1c float64 //= models.FloatField(verbose_name="买一档笔数", default=0)
	B1m float64 //= models.FloatField(verbose_name="买一档金额", default=0)
	B2c float64 //= models.FloatField(verbose_name="买2档笔数", default=0)
	B2m float64 //= models.FloatField(verbose_name="买2档金额", default=0)
	B3c float64 //= models.FloatField(verbose_name="买3档笔数", default=0)
	B3m float64 //= models.FloatField(verbose_name="买3档金额", default=0)
	B4c float64 //= models.FloatField(verbose_name="买4档笔数", default=0)
	B4m float64 //= models.FloatField(verbose_name="买4档金额", default=0)
	B5c float64 //= models.FloatField(verbose_name="买5档笔数", default=0)
	B5m float64 //= models.FloatField(verbose_name="买5档金额", default=0)

	S1c float64 //= models.FloatField(verbose_name="卖一档笔数", default=0)
	S1m float64 //= models.FloatField(verbose_name="卖一档金额", default=0)
	S2c float64 //= models.FloatField(verbose_name="卖2档笔数", default=0)
	S2m float64 //= models.FloatField(verbose_name="卖2档金额", default=0)
	S3c float64 //= models.FloatField(verbose_name="卖3档笔数", default=0)
	S3m float64 //= models.FloatField(verbose_name="卖3档金额", default=0)
	S4c float64 //= models.FloatField(verbose_name="卖4档笔数", default=0)
	S4m float64 //= models.FloatField(verbose_name="卖4档金额", default=0)
	S5c float64 //= models.FloatField(verbose_name="卖5档笔数", default=0)
	S5m float64 //= models.FloatField(verbose_name="卖5档金额", default=0)

	Max float64 //= models.FloatField(verbose_name="今日最高", default=0)
	Min float64 //= models.FloatField(verbose_name="今日最低", default=0)

	Jin_buy_1 float64 //= models.FloatField(verbose_name="竞买1", default=0)
	Jin_sale_1 float64 //= models.FloatField(verbose_name="竞卖1", default=0)

	Total_count float64 //= models.FloatField(verbose_name="总成交数", default=0)
	Total_money float64 //= models.FloatField(verbose_name="总成交金额", default=0)
}

func (Stock) TableName() string {
	return "stock"
}

type Code struct {
	ID  uint `gorm:"primary_key"`
	Code string `gorm:"not null;unique"`
	Name string
	Monitor bool
	Belong string // 所属行业板块
	Concept string // 所属概念板块
}


type Detail struct {
	gorm.Model
	Code string
	Date string
	CurrTime string
	Curr string
	Inner float64
	Outer float64
}


type InOut struct {
	gorm.Model
	Code string
	Zin float64 // 主力
	Zout float64
	Zjin float64
	Zpercent float64
	Sin float64 // 散户
	Sout float64
	Sjin float64
	Spercent float64

	Total float64
	Date string
}

func (Code) TableName() string {
	return "code"
}

func (Detail) TableName() string {
	return "detail"
}

func (InOut) TableName() string {
	return "inout"
}

type GrailHistory struct {
	gorm.Model
	Date string
	Kai, High, Low, Shou, ChangeAmount, Percent, TotalCount, TotalMoney float64
}

func (GrailHistory) TableName() string {
	return "grail_history"
}


type TicketHistory struct {
	gorm.Model
	Code string `gorm:"index"`
	Name string
	Date string `gorm:"index"`
	Kai, High, Shou, Low, TotalCount, TotalMoney float64
	Percent float64
}

func (TicketHistory) TableName() string {
	return "ticket_history"
}

type TicketHistoryTmp struct {
	gorm.Model
	Code string
	Name string
	Date string
	Kai, High, Shou, Low, TotalCount, TotalMoney float64
	Percent float64
}

func (TicketHistoryTmp) TableName() string {
	return "ticket_history_tmp"
}


type TicketHistoryWeekly struct {
	gorm.Model
	Code string `gorm:"index"`
	Name string
	Date string `gorm:"index"`
	Shou float64
	Percent float64  // 周线涨跌幅
}

func (TicketHistoryWeekly) TableName() string {
	return "ticket_history_week"
}


type TicketHistoryHour struct {
	gorm.Model
	Code string
	Name string
	Date string
	Time string

	Open float64
	High float64
	Low float64
	Close float64
	Volume float64
}

func (TicketHistoryHour) TableName() string {
	return "ticket_history_hour"
}

type Analyze struct {
	gorm.Model
	Code string
	Name string
	Date string
	OnSix bool
	OnFifteen bool
	OnThirty bool
	OnSixWeek bool
	OnFifteenWeek bool
}

func (Analyze) TableName() string {
	return "analyze"
}

type Predict struct {
	ID  int `gorm:"primary_key"`
	Code string
	Name string
	Date string
	Condition string
	RealPercent float64 // 真实的涨跌情况， 用来判断准确性

	Attrs string
	Holder string

	FundCount int // 一共几只基金持有
}

func (Predict) TableName() string {
	return "predict"
}


type PredictDebug struct {
	ID  int `gorm:"primary_key"`
	Code string
	Name string
	Date string
	Condition string
	RealPercent float64 // 真实的涨跌情况， 用来判断准确性

	Attrs string
	Holder string

	FundCount int // 一共几只基金持有

}

func (PredictDebug) TableName() string {
	return "predict_debug"
}

type MACD struct {
	//EMA（12）= 前一日EMA（12）×11/13＋今日收盘价×2/13
	//EMA（26）= 前一日EMA（26）×25/27＋今日收盘价×2/27
	//DIFF=今日EMA（12）- 今日EMA（26）
	//DEA（MACD）= 前一日DEA×8/10＋今日DIF×2/10
	//MACD =2×(DIFF－DEA)
	gorm.Model
	Code string
	Date string
	Ema1 float64
	Ema2 float64
	Diff float64
	Dea float64
	Macd float64
}

func (MACD) TableName() string {
	return "macd"
}


type AllTicketHistory struct {
	gorm.Model
	Code string
	Name string
	Date string
	Kai, High, Shou, Low, TotalCount, TotalMoney float64
}

func (AllTicketHistory) TableName() string {
	return "all_ticket_history"
}

type Stockholder struct {
	gorm.Model
	Code string `gorm:"index"`
	Name string
	HolderName string
	Count string
	Percent string
	Change string
}

func (Stockholder) TableName() string {
	return "stock_holder"
}

type LongHuRankDepartmentTop30 struct {
	gorm.Model
	Department string
	Count int
	BuyMoney float64
	BuyCount int
	SaleMoney float64
	SaleCount int
	Tickets string
}

func (LongHuRankDepartmentTop30) TableName() string {
	return "longhu_department_top30"
}

type LongHuRankDepartmentTop10 struct {
	gorm.Model
	Department string
	Count int
	BuyMoney float64
	BuyCount int
	SaleMoney float64
	SaleCount int
	Tickets string
}

func (LongHuRankDepartmentTop10) TableName() string {
	return "longhu_department_top10"
}

type LongHuRankDepartmentTop5 struct {
	gorm.Model
	Department string
	Count int
	BuyMoney float64
	BuyCount int
	SaleMoney float64
	SaleCount int
	Tickets string
}

func (LongHuRankDepartmentTop5) TableName() string {
	return "longhu_department_top5"
}


type LongHuRankTicket struct {
	gorm.Model
	Code string
	Name string
	Count int  // 上板次数
	BuyCount int
	SaleCount int
	BuyTotal float64
	SaleTotal float64
	Total float64 // 买入 - 卖出 净额
}


func (LongHuRankTicket) TableName() string {
	return "longhu_ticket"
}


type LongHuRankTicketDetail struct {
	gorm.Model
	Code string
	Name string
	YYBCode   string //yyb
	YYBName   string
	BuyTotal  float64
	SaleTotal float64
	Date string
}


func (LongHuRankTicketDetail) TableName() string {
	return "longhu_ticket_detail"
}


// 基金排行
type FundRank struct {
	gorm.Model
	Type string
	FundCode string
	FundSName string
	LastMonth float64
	LastWeek float64
	Last3Month float64
	Last6Month float64
	SinceBase float64  // 成立以来涨跌幅
	ThisYear float64
	LastYear float64
	Last2Year float64
	Last3Year float64
	Time string
}

func (FundRank) TableName() string {
	return "fund_rank"
}


// 基金持仓
type FundHoldRank struct {
	gorm.Model
	Type string  // 基金类型
	FundCode string
	FundSName string
	Code string
	Name string
	Percent float64
	Concept string
	Time string
}

func (FundHoldRank) TableName() string {
	return "fund_hold"
}

