// @Time:       2019/12/3 下午2:33

package crawler

import (
	"fmt"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/utils"
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	LIMIT_UP = 9.9
	Main     = 0.1
	Virtual  = 0.1
	Num      = 3
)

func ArraySumLessThan(array []float64, less float64) bool {
	if len(array) < 1 {
		return false
	}
	var x float64
	for _, i := range array {
		x += i
	}
	return x < less
}

// 有最大值
func ArrayHasLimitUp(array []float64, flag int) bool {
	if len(array) < 1 {
		return false
	}
	for _, i := range array[0:flag] {
		if i > LIMIT_UP {
			return true
		}
	}
	return false
}

// 指定大小
func Array1BiggerThanArray2ByNumber(arr1, arr2 []float64, recent int, total int) bool {
	var tmp []int
	for i := 0; i < recent-1; i++ {
		if arr1[i] > arr2[i] {
			tmp = append(tmp, 1)
		}
	}
	return len(tmp) >= total
}

// 指定百分比
func Array1BiggerThanArray2ByPercent(recent_count, recent_ave []float64, percent float64, recent, total int) bool {
	// percent = 1 就是两倍的意思
	var tmp []int
	for i := 0; i < recent-1; i++ {
		if (recent_count[i]-recent_ave[i])/recent_ave[i] > percent {
			tmp = append(tmp, 1)
		}
	}
	return len(tmp) >= total
}

func ConditionTopLine(array_shou, array_high, array_kai, percent []float64, recent int) bool {
	if len(array_shou) < 1 {
		return false
	}
	for i := 0; i <= recent-1; i++ {
		if ((array_high[i]-array_shou[i])/array_shou[i])*100/percent[i] > 5 {
			return true
		}
	}
	return false
}

// 如何判断实体柱有意义， 实体柱在昨日收盘价*0.1的数据中 占比%10即可算是有意义  也就是 main > 昨日收盘价*0.1*0.1, 同理影线长度也需要做限制
// 虚线柱子有4个点以上即可 简单粗暴

func (craw *Crawler) HasTopLine(result *model.CalcResult, recent int) bool {
	closes := result.RecentClose[0:recent]
	for index, i := range closes {
		// 开盘价大于收盘价
		if result.RecentOpen[index]-i >= 0 {
			virtual := result.RecentHigh[index] - result.RecentOpen[index]
			if virtual > result.RecentClose[index+1]*0.04 {
				return true
			}
		} else {
			virtual := result.RecentHigh[index] - result.RecentClose[index]
			if virtual > result.RecentClose[index+1]*0.04 {
				return true
			}
		}
	}
	return false
}

// 下影线是实体柱子4倍以上
func (craw *Crawler) HasLowLine(result *model.CalcResult, recent int) bool {
	closes := result.RecentClose[0:recent]
	for index, i := range closes {
		// 开盘价大于收盘价
		if result.RecentOpen[index]-i >= 0 {
			virtual := result.RecentClose[index] - result.RecentLow[index]
			if virtual > result.RecentClose[index+1]*0.04 {
				return true
			}
		} else {
			virtual := result.RecentOpen[index] - result.RecentLow[index]
			if virtual > result.RecentClose[index+1]*0.04 {
				return true
			}
		}
	}
	return false
}

func GetConceptByCode(code, concept string) bool {
	var c int
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", code).Where("`concept` regexp ?", concept).Count(&c)
	if c > 0 {
		return true
	}
	return false
}

func GetHolderByCode(code, concept string) int {
	var c int
	tmp := store.MysqlClient.GetDB().Model(&dal.Stockholder{}).Where("code = ?", code).Where("holder_name regexp ?", concept)
	tmp.Count(&c)
	return c
}

func GetFundByCode(code string) int {
	var count int
	store.MysqlClient.GetDB().Model(&dal.StockFund{}).Where("code = ?", code).Count(&count)
	return count
}

func GetFenghongByCode(code string) (int, int, int) {
	var fenhong, songgu, zhuangzeng int
	store.MysqlClient.GetDB().Model(&dal.StockFengHong{}).Where("pai_xi > ? and code = ?", 0, code).Count(&fenhong)
	store.MysqlClient.GetDB().Model(&dal.StockFengHong{}).Where("song_gu > ? and code = ?", 0, code).Count(&songgu)
	store.MysqlClient.GetDB().Model(&dal.StockFengHong{}).Where("zhuang_zeng > ? and code = ?", 0, code).Count(&zhuangzeng)
	return fenhong, songgu, zhuangzeng
}

func GetPeiGuByCode(code string) (int, int) {
	var pergu, zengfa int
	store.MysqlClient.GetDB().Model(&dal.StockPeiGu{}).Where("code = ?", code).Count(&pergu)
	store.MysqlClient.GetDB().Model(&dal.StockZengFa{}).Where("code = ?", code).Count(&zengfa)
	return pergu, zengfa
}

func STStock(code string) bool {
	var c int
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", code).Where("`concept` regexp ?", "ST").Count(&c)
	return c > 0
}

func GetHistoryNameByCode(code string) (int, bool) {
	var c dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", code).Find(&c)
	if c.HistoryNames != "暂无更名记录" {
		names := strings.Split(c.HistoryNames, " ")
		for _, i := range names {
			if strings.Contains(i, "ST") {
				return len(names), true
			}
		}
		return len(names), false
	} else {
		return 0, false
	}
}

// 不断增加
func calc(res []float64, flag string) bool {
	if len(res) == 0 {
		return false
	}
	var first = res[0]
	for _, i := range res {
		if flag == "desc" {
			if i < first {
				return false
			}
		}
		if flag == "asc" {
			if i > first {
				return false
			}
		}
		first = i
	}
	return true
}

// 不能为负
func calc2(res []float64) bool {
	if len(res) == 0 {
		return false
	}
	for _, i := range res {
		if i < 0 {
			return false
		}
	}
	return true
}

// 总和不能为负
func calc3(res []float64) bool {
	var a float64
	for _, i := range res {
		a += i
	}
	return a > 0
}

func GetUpManageCashFlow(code string) (up1, up2, up3, up4 bool) {
	var cash_flow []dal.StockCashFlow
	var ManageCashFlow, InvestCashFlow, FundraisingCashFlow, CashRemain []float64
	store.MysqlClient.GetDB().Model(&dal.StockCashFlow{}).Where("code = ?", code).Order("date desc").Find(&cash_flow)
	for _, i := range cash_flow {
		if i.ManageCashFlow != 0 {
			ManageCashFlow = append(ManageCashFlow, i.ManageCashFlow)
		}
		if i.InvestCashFlow != 0 {
			InvestCashFlow = append(InvestCashFlow, i.InvestCashFlow)
		}
		if i.FundraisingCashFlow != 0 {
			FundraisingCashFlow = append(FundraisingCashFlow, i.FundraisingCashFlow)
		}
		if i.CashRemain != 0 {
			CashRemain = append(CashRemain, i.CashRemain)
		}
	}
	up1 = calc2(ManageCashFlow)
	up2 = calc2(InvestCashFlow)
	up3 = calc2(FundraisingCashFlow)
	up4 = calc2(CashRemain)
	return
}

func GetProfitNotMiner(code string) (up1, up2 bool) {
	var profit []dal.StockProfit
	var NetProfit, GrossTradingIncome []float64
	store.MysqlClient.GetDB().Model(&dal.StockProfit{}).Where("code = ?", code).Order("date desc").Find(&profit)
	for _, i := range profit {
		if i.GrossTradingIncome != 0 {
			GrossTradingIncome = append(GrossTradingIncome, i.GrossTradingIncome)
		}
		if i.NetProfit != 0 {
			NetProfit = append(NetProfit, i.NetProfit)
		}
	}
	up1 = calc2(GrossTradingIncome)
	up2 = calc2(NetProfit)
	return
}

func GetUpLiabilities(code string) (up1, done1 bool) {
	var StockLiabilities []dal.StockLiabilities
	var TotalAssets, TotalLiabilities []float64 // 资产总计 负债总计
	store.MysqlClient.GetDB().Model(&dal.StockLiabilities{}).Where("code = ?", code).Order("date desc").Find(&StockLiabilities)
	for _, i := range StockLiabilities {
		if i.TotalAssets != 0 {
			TotalAssets = append(TotalAssets, i.TotalAssets)
		}
		if i.TotalLiabilities != 0 {
			TotalLiabilities = append(TotalLiabilities, i.TotalLiabilities)
		}
	}
	up1 = calc(TotalAssets, "desc")
	done1 = calc(TotalLiabilities, "asc")
	return
}

// 最近几日量价和均线有几日满足关系条件（在误差范围内）
func RecentInRangeAveWithCond(recent_money_or_count, recent_ave []float64, recent int, total int) bool {
	var tmp []int
	for i := 0; i < recent-1; i++ {
		if math.Abs((recent_money_or_count[i]-recent_ave[i])/recent_ave[i]) < 0.01 {
			tmp = append(tmp, 1)
		}
	}
	return len(tmp) >= total
}

func (craw *Crawler) HasLimitUpInTheseDays(result *model.CalcResult, recent_days int) bool {
	return ArrayHasLimitUp(result.RecentPercent, recent_days)
}

// 5 10 15 30 60
func (craw *Crawler) Analyze(result *model.CalcResult, code, name string) {
	// 5 10 金叉
	jincha1 := result.AveDailyPrice1[0] > result.AveDailyPrice2[0] && result.AveDailyPrice1[1] < result.AveDailyPrice2[1]
	// 5 15 金叉
	jincha2 := result.AveDailyPrice1[0] > result.AveDailyPrice3[0] && result.AveDailyPrice1[1] < result.AveDailyPrice3[1]
	// 5 30 金叉
	jincha3 := result.AveDailyPrice1[0] > result.AveDailyPrice4[0] && result.AveDailyPrice1[1] < result.AveDailyPrice4[1]
	// 5 60 金叉
	jincha4 := result.AveDailyPrice1[0] > result.AveDailyPrice5[0] && result.AveDailyPrice1[1] < result.AveDailyPrice5[1]
	// 10 15 金叉
	jincha5 := result.AveDailyPrice2[0] > result.AveDailyPrice3[0] && result.AveDailyPrice2[1] < result.AveDailyPrice3[1]
	// 10 30 金叉
	jincha6 := result.AveDailyPrice2[0] > result.AveDailyPrice4[0] && result.AveDailyPrice2[1] < result.AveDailyPrice4[1]
	// 10 60 金叉
	jincha7 := result.AveDailyPrice2[0] > result.AveDailyPrice5[0] && result.AveDailyPrice2[1] < result.AveDailyPrice5[1]
	// 15 30 金叉
	jincha8 := result.AveDailyPrice3[0] > result.AveDailyPrice4[0] && result.AveDailyPrice3[1] < result.AveDailyPrice4[1]
	// 15 60 金叉
	jincha9 := result.AveDailyPrice3[0] > result.AveDailyPrice5[0] && result.AveDailyPrice3[1] < result.AveDailyPrice5[1]
	// 30 60 金叉
	jincha10 := result.AveDailyPrice4[0] > result.AveDailyPrice5[0] && result.AveDailyPrice4[1] < result.AveDailyPrice5[1]
	// 10 40 量能金叉
	jincha11 := result.AveCount1[0] > result.AveCount2[0] && result.AveCount1[1] < result.AveCount2[1]

	// 死叉股
	sicha1 := result.AveDailyPrice1[0] < result.AveDailyPrice2[0] && result.AveDailyPrice1[1] > result.AveDailyPrice2[1]
	sicha2 := result.AveDailyPrice1[0] < result.AveDailyPrice3[0] && result.AveDailyPrice1[1] > result.AveDailyPrice3[1]
	sicha3 := result.AveDailyPrice1[0] < result.AveDailyPrice4[0] && result.AveDailyPrice1[1] > result.AveDailyPrice4[1]
	sicha4 := result.AveDailyPrice1[0] < result.AveDailyPrice5[0] && result.AveDailyPrice1[1] > result.AveDailyPrice5[1]
	sicha5 := result.AveDailyPrice2[0] < result.AveDailyPrice3[0] && result.AveDailyPrice2[1] > result.AveDailyPrice3[1]
	sicha6 := result.AveDailyPrice2[0] < result.AveDailyPrice4[0] && result.AveDailyPrice2[1] > result.AveDailyPrice4[1]
	sicha7 := result.AveDailyPrice2[0] < result.AveDailyPrice5[0] && result.AveDailyPrice2[1] > result.AveDailyPrice5[1]
	sicha8 := result.AveDailyPrice3[0] < result.AveDailyPrice4[0] && result.AveDailyPrice3[1] > result.AveDailyPrice4[1]
	sicha9 := result.AveDailyPrice3[0] < result.AveDailyPrice5[0] && result.AveDailyPrice3[1] > result.AveDailyPrice5[1]
	sicha10 := result.AveDailyPrice4[0] < result.AveDailyPrice5[0] && result.AveDailyPrice4[1] > result.AveDailyPrice5[1]
	sicha11 := result.AveCount1[0] < result.AveCount2[0] && result.AveCount1[1] < result.AveCount2[1]

	// 涨停股
	zhangting := result.RecentPercent[0] > 9.94
	// 一字板
	yiziban := result.RecentPercent[0] > 9.94 && result.RecentOpen[0] == result.RecentLow[0]
	// T 字板
	tziban := result.RecentPercent[0] > 9.94 && result.RecentOpen[0] == result.RecentClose[0] && result.RecentClose[0] > result.RecentLow[0]

	// 5条均线 价格均线上扬
	priceshangyang1 := result.AveDailyPrice1[0] > result.AveDailyPrice1[1] && result.AveDailyPrice1[1] > result.AveDailyPrice1[2] && result.AveDailyPrice1[2] > result.AveDailyPrice1[3] && result.AveDailyPrice1[3] > result.AveDailyPrice1[4] && result.AveDailyPrice1[4] > result.AveDailyPrice1[5]
	priceshangyang2 := result.AveDailyPrice2[0] > result.AveDailyPrice2[1] && result.AveDailyPrice2[1] > result.AveDailyPrice2[2] && result.AveDailyPrice2[2] > result.AveDailyPrice2[3] && result.AveDailyPrice2[3] > result.AveDailyPrice2[4] && result.AveDailyPrice2[4] > result.AveDailyPrice2[5]
	priceshangyang3 := result.AveDailyPrice3[0] > result.AveDailyPrice3[1] && result.AveDailyPrice3[1] > result.AveDailyPrice3[2] && result.AveDailyPrice3[2] > result.AveDailyPrice3[3] && result.AveDailyPrice3[3] > result.AveDailyPrice3[4] && result.AveDailyPrice3[4] > result.AveDailyPrice3[5]
	priceshangyang4 := result.AveDailyPrice4[0] > result.AveDailyPrice4[1] && result.AveDailyPrice4[1] > result.AveDailyPrice4[2] && result.AveDailyPrice4[2] > result.AveDailyPrice4[3] && result.AveDailyPrice4[3] > result.AveDailyPrice4[4] && result.AveDailyPrice4[4] > result.AveDailyPrice4[5]
	priceshangyang5 := result.AveDailyPrice5[0] > result.AveDailyPrice5[1] && result.AveDailyPrice5[1] > result.AveDailyPrice5[2] && result.AveDailyPrice5[2] > result.AveDailyPrice5[3] && result.AveDailyPrice5[3] > result.AveDailyPrice5[4] && result.AveDailyPrice5[4] > result.AveDailyPrice5[5]
	// 5条均线 价格均线下降
	pricexiajiang1 := result.AveDailyPrice1[0] < result.AveDailyPrice1[1] && result.AveDailyPrice1[1] < result.AveDailyPrice1[2] && result.AveDailyPrice1[2] < result.AveDailyPrice1[3] && result.AveDailyPrice1[3] < result.AveDailyPrice1[4] && result.AveDailyPrice1[4] < result.AveDailyPrice1[5]
	pricexiajiang2 := result.AveDailyPrice2[0] < result.AveDailyPrice2[1] && result.AveDailyPrice2[1] < result.AveDailyPrice2[2] && result.AveDailyPrice2[2] < result.AveDailyPrice2[3] && result.AveDailyPrice2[3] < result.AveDailyPrice2[4] && result.AveDailyPrice2[4] < result.AveDailyPrice2[5]
	pricexiajiang3 := result.AveDailyPrice3[0] < result.AveDailyPrice3[1] && result.AveDailyPrice3[1] < result.AveDailyPrice3[2] && result.AveDailyPrice3[2] < result.AveDailyPrice3[3] && result.AveDailyPrice3[3] < result.AveDailyPrice3[4] && result.AveDailyPrice3[4] < result.AveDailyPrice3[5]
	pricexiajiang4 := result.AveDailyPrice4[0] < result.AveDailyPrice4[1] && result.AveDailyPrice4[1] < result.AveDailyPrice4[2] && result.AveDailyPrice4[2] < result.AveDailyPrice4[3] && result.AveDailyPrice4[3] < result.AveDailyPrice4[4] && result.AveDailyPrice4[4] < result.AveDailyPrice4[5]
	pricexiajiang5 := result.AveDailyPrice5[0] < result.AveDailyPrice5[1] && result.AveDailyPrice5[1] < result.AveDailyPrice5[2] && result.AveDailyPrice5[2] < result.AveDailyPrice5[3] && result.AveDailyPrice5[3] < result.AveDailyPrice5[4] && result.AveDailyPrice5[4] < result.AveDailyPrice5[5]

	// 当前价格在短期均线上方 （取非为小于）
	priceaboveave1 := result.RecentClose[0] >= result.AveDailyPrice1[0]
	priceaboveave2 := result.RecentClose[0] >= result.AveDailyPrice2[0]
	priceaboveave3 := result.RecentClose[0] >= result.AveDailyPrice3[0]
	priceaboveave4 := result.RecentClose[0] >= result.AveDailyPrice4[0]
	priceaboveave5 := result.RecentClose[0] >= result.AveDailyPrice5[0]

	// 均价粘合
	junjialianhe1 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice1, 5, 4)
	junjialianhe2 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice2, 5, 4)
	junjialianhe3 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice3, 5, 4)
	junjialianhe4 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice4, 5, 4)
	junjialianhe5 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice5, 5, 4)

	// 低开高走
	dikaigaozou := (result.RecentClose[1]-result.RecentOpen[0])/result.RecentClose[1] > 0.02 && result.RecentPercent[0] > 3
	// 高开低走
	gaokaidizou := (result.RecentOpen[0]-result.RecentClose[1])/result.RecentOpen[0] > 0.02 && result.RecentPercent[0] < -3
	// 低开低走
	dikaidizou := (result.RecentClose[1]-result.RecentOpen[0])/result.RecentClose[1] > 0.02 && result.RecentPercent[0] < -5
	// 高开高走
	gaokaigaozou := (result.RecentOpen[0]-result.RecentClose[1])/result.RecentOpen[0] > 0.02 && result.RecentPercent[0] > 5

	// 3 连阳
	sanlianyang := result.RecentPercent[0] > 0 && result.RecentPercent[1] > 0 && result.RecentPercent[2] > 0 && result.RecentPercent[3] > 0
	// 4 连阳
	silianyang := sanlianyang && result.RecentPercent[4] > 0
	// 5连阳
	wulianyang := silianyang && result.RecentPercent[5] > 0

	// 长上影
	changshangying := craw.HasTopLine(result, 1)
	// 长下影
	changxiaying := craw.HasLowLine(result, 1)

	// 连续5日量能10均线上扬
	liangshangyang1 := result.AveCount1[0] > result.AveCount1[1] && result.AveCount1[1] > result.AveCount1[2] && result.AveCount1[2] > result.AveCount1[3] && result.AveCount1[3] > result.AveCount1[4] && result.AveCount1[4] > result.AveCount1[5]
	// 连续5日量能40均线上扬
	liangshangyang2 := result.AveCount2[0] > result.AveCount2[1] && result.AveCount2[1] > result.AveCount2[2] && result.AveCount2[2] > result.AveCount2[3] && result.AveCount2[3] > result.AveCount2[4] && result.AveCount2[4] > result.AveCount2[5]
	// 连续5日量能站上10均线
	liangnengbigger1 := result.AveCount1[0] < result.RecentCount[0] && result.AveCount1[1] < result.RecentCount[1] && result.AveCount1[2] < result.RecentCount[2] && result.AveCount1[3] < result.RecentCount[3] && result.AveCount1[4] < result.RecentCount[4]
	// 连续5日量能站上40均线
	liangnengbigger2 := result.AveCount2[0] < result.RecentCount[0] && result.AveCount2[1] < result.RecentCount[1] && result.AveCount2[2] < result.RecentCount[2] && result.AveCount2[3] < result.RecentCount[3] && result.AveCount2[4] < result.RecentCount[4]
	// 连续5日量能低于10均线
	liangnengsmaller1 := result.AveCount1[0] > result.RecentCount[0] && result.AveCount1[1] > result.RecentCount[1] && result.AveCount1[2] > result.RecentCount[2] && result.AveCount1[3] > result.RecentCount[3] && result.AveCount1[4] > result.RecentCount[4]
	// 连续5日量能低于40均线
	liangnengsmaller2 := result.AveCount2[0] > result.RecentCount[0] && result.AveCount2[1] > result.RecentCount[1] && result.AveCount2[2] > result.RecentCount[2] && result.AveCount2[3] > result.RecentCount[3] && result.AveCount2[4] > result.RecentCount[4]
	// 量能不断放大
	liangnengbuduanbigger := result.RecentCount[0] > result.RecentCount[1] && result.RecentCount[1] > result.RecentCount[2] && result.RecentCount[2] > result.RecentCount[3]
	// 突放巨量
	tufangjuliang := (result.RecentCount[0]-result.RecentCount[1])/result.RecentCount[1] > 5 || (result.RecentCount[1]-result.RecentCount[2])/result.RecentCount[1] > 5

	// 私募持仓
	simuchicangcount := GetHolderByCode(code, "私募")
	// 基金持仓
	jigouchicangcount := GetFundByCode(code)

	// 成交过亿
	guoyi := result.CurrTotalMoney > 10000

	// 现金流量表 非负
	up1, up2, up3, up4 := GetUpManageCashFlow(code)
	// 利润表 不能为负
	pup1, pup2 := GetProfitNotMiner(code)
	// 资产负债表
	lup1, done1 := GetUpLiabilities(code)
	// 分红配股的次数
	fenhong, songgu, zhuangzeng := GetFenghongByCode(code)
	pergu, zengfa := GetPeiGuByCode(code)
	// 历史更名次数
	changename, has_st := GetHistoryNameByCode(code)

	score := 0 // max 37  // low 17
	cond_str, bad_cond_str, finance := "", "", ""

	if jincha11 {
		score += 2
		cond_str += "(量能金叉)10x40; "
	}
	if jincha1 {
		score += 2
		cond_str += "(金叉)5x10; "
	}
	if jincha2 {
		score += 2
		cond_str += "(金叉)5x15; "
	}
	if jincha3 {
		score += 2
		cond_str += "(金叉)5x30; "
	}
	if jincha4 {
		score += 2
		cond_str += "(金叉)5x60; "
	}
	if jincha5 {
		score += 2
		cond_str += "(金叉)10x15; "
	}
	if jincha6 {
		score += 2
		cond_str += "(金叉)10x30; "
	}
	if jincha7 {
		score += 2
		cond_str += "(金叉)10x60; "
	}
	if jincha8 {
		score += 2
		cond_str += "(金叉)15x30; "
	}
	if jincha9 {
		score += 2
		cond_str += "(金叉)15x60; "
	}
	if jincha10 {
		score += 2
		cond_str += "(金叉)30x60; "
	}

	if sicha11 {
		score -= 2
		bad_cond_str += "(量能死叉)10x40; "
	}
	if sicha1 {
		score -= 2
		bad_cond_str += "(死叉)5x10; "
	}
	if sicha2 {
		score -= 2
		bad_cond_str += "(死叉)5x15; "
	}
	if sicha3 {
		score -= 2
		bad_cond_str += "(死叉)5x30; "
	}
	if sicha4 {
		score -= 2
		bad_cond_str += "(死叉)5x60; "
	}
	if sicha5 {
		score -= 2
		bad_cond_str += "(死叉)10x15; "
	}
	if sicha6 {
		score -= 2
		bad_cond_str += "(死叉)10x30; "
	}
	if sicha7 {
		score -= 2
		bad_cond_str += "(死叉)10x60; "
	}
	if sicha8 {
		score -= 2
		bad_cond_str += "(死叉)15x30; "
	}
	if sicha9 {
		score -= 2
		bad_cond_str += "(死叉)15x60; "
	}
	if sicha10 {
		score -= 2
		bad_cond_str += "(死叉)30x60; "
	}

	if yiziban {
		cond_str += "一字板; "
	}
	if tziban {
		cond_str += "T字板; "
	}
	if zhangting {
		score += 2
		cond_str += "涨停股; "
	}

	if priceshangyang1 {
		score += 2
		cond_str += "5均上升通道; "
	}
	if priceshangyang2 {
		score += 2
		cond_str += "10均上升通道; "
	}
	if priceshangyang3 {
		score += 2
		cond_str += "15均上升通道; "
	}
	if priceshangyang4 {
		score += 2
		cond_str += "30均上升通道; "
	}
	if priceshangyang5 {
		score += 2
		cond_str += "60均上升通道; "
	}

	if pricexiajiang1 {
		score -= 2
		bad_cond_str += "5均下降通道; "
	}
	if pricexiajiang2 {
		score -= 2
		bad_cond_str += "10均下降通道; "
	}
	if pricexiajiang3 {
		score -= 2
		bad_cond_str += "15均下降通道; "
	}
	if pricexiajiang4 {
		score -= 2
		bad_cond_str += "30均下降通道; "
	}
	if pricexiajiang5 {
		score -= 2
		bad_cond_str += "60均下降通道; "
	}

	if priceaboveave1 {
		score += 1
		cond_str += "价位在5均上方; "
	}
	if priceaboveave2 {
		score += 1
		cond_str += "价位在10均上方; "
	}
	if priceaboveave3 {
		score += 1
		cond_str += "价位在15均上方; "
	}
	if priceaboveave4 {
		score += 1
		cond_str += "价位在30均上方; "
	}
	if priceaboveave5 {
		score += 1
		cond_str += "价位在60均上方; "
	}

	if !priceaboveave1 {
		score -= 1
		bad_cond_str += "价位在5均下方; "
	}
	if !priceaboveave2 {
		score -= 1
		bad_cond_str += "价位在10均下方; "
	}
	if !priceaboveave3 {
		score -= 1
		bad_cond_str += "价位在15均下方; "
	}
	if !priceaboveave4 {
		score -= 1
		bad_cond_str += "价位在30均下方; "
	}
	if !priceaboveave5 {
		score -= 1
		bad_cond_str += "价位在60均下方; "
	}

	if junjialianhe1 {
		cond_str += "5均价黏合; "
	}
	if junjialianhe2 {
		cond_str += "10均价黏合; "
	}
	if junjialianhe3 {
		cond_str += "15均价黏合; "
	}
	if junjialianhe4 {
		cond_str += "30均价黏合; "
	}
	if junjialianhe5 {
		cond_str += "60均价黏合; "
	}

	if dikaigaozou {
		cond_str += "低开高走; "
	}
	if gaokaigaozou {
		cond_str += "高开高走; "
	}
	if gaokaidizou {
		cond_str += "高开低走; "
	}
	if dikaidizou {
		cond_str += "低开低走; "
	}

	if wulianyang {
		score += 3
		cond_str += "五连阳; "
	}
	if !wulianyang && silianyang {
		score += 2
		cond_str += "四连阳; "
	}
	if !wulianyang && !silianyang && sanlianyang {
		score += 3
		cond_str += "三连阳; "
	}
	if changshangying {
		cond_str += "长上影; "
	}
	if changxiaying {
		cond_str += "长下影; "
	}

	// 量价
	if liangshangyang1 {
		score += 3
		cond_str += "量能10均上升通道; "
	}
	if liangshangyang2 {
		score += 3
		cond_str += "量能40均上升通道; "
	}
	if liangnengbigger1 || liangnengbigger2 {
		score += 3
		cond_str += "量能活跃; "
	}
	if liangnengsmaller1 || liangnengsmaller2 {
		score -= 2
		bad_cond_str += "量能萎靡; "
	}
	if liangnengbuduanbigger {
		score += 1
		cond_str += "量能不断放大; "
	}
	if tufangjuliang {
		score += 2
		cond_str += "突放巨量; "
	}
	if guoyi {
		cond_str += "成交额过亿; "
	}

	if simuchicangcount > 0 {
		score += 3
		cond_str += fmt.Sprintf("%d个私募持仓; ", simuchicangcount)
	}
	if jigouchicangcount > 0 {
		score += 3
		cond_str += fmt.Sprintf("%d个基金持仓; ", jigouchicangcount)
	}
	if fenhong > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d次分红; ", fenhong)
	}
	if songgu > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d次送股; ", songgu)
	}
	if zhuangzeng > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d次转增; ", zhuangzeng)
	}
	if pergu > 0 {
		cond_str += fmt.Sprintf("%d次配股; ", pergu)
	}
	if zengfa > 0 {
		cond_str += fmt.Sprintf("%d次增发; ", zengfa)
	}
	if changename > 0 {
		bad_cond_str += fmt.Sprintf("历史更名%d次; ", changename)
		if has_st {
			bad_cond_str += "曾ST带帽; "
		}
	}

	// 基本面 现金流量表
	if up1 {
		score += 1
		cond_str += "经营现金流量净额非负; "
	}
	if up2 {
		score += 1
		cond_str += "投资现金流量净额非负; "
	}
	if up3 {
		score += 1
		cond_str += "筹资现金流量净额非负; "
	}
	if up4 {
		score += 1
		cond_str += "期末现金及现金等价物余额非负; "
	}
	// 基本面 利润表
	if pup1 {
		score += 1
		cond_str += "营业总收入非负; "
	}
	if pup2 {
		score += 1
		cond_str += "净利润非负; "
	}
	// 基本面 资产负债表
	if lup1 {
		score += 1
		cond_str += "总资产不断增加; "
	}
	if done1 {
		score += 1
		cond_str += "总负债不断减小; "
	}

	if !up1 {
		score -= 1
		bad_cond_str += "(亏损可能)经营现金流量净额出现负值; "
	}
	if !up2 {
		score -= 1
		bad_cond_str += "(亏损可能)投资现金流量净额出现负值; "
	}
	if !up3 {
		score -= 1
		bad_cond_str += "(亏损可能)筹资现金流量净额出现负值; "
	}
	if !up4 {
		score -= 1
		bad_cond_str += "(亏损可能)期末现金及现金等价物余额出现负值; "
	}
	if !pup1 {
		score -= 2
		bad_cond_str += "营业总收入亏损; "
	}
	if !pup2 {
		score -= 2
		bad_cond_str += "净利润亏损; "
	}

	if score < 0 {
		score = 0
	}

	finance += ""
	var per dal.StockPerTicket
	err := store.MysqlClient.GetDB().Model(&dal.StockPerTicket{}).Where("code = ?", code).Find(&per).Error
	if err != nil {
		finance = ""
	} else {
		finance = per.RankCaiwu
	}
	high := strings.Count(finance, "高")
	middle := strings.Count(finance, "一般")
	low := strings.Count(finance, "偏低")
	bad := strings.Count(finance, "较差")
	score += high*4 + middle*3 + low*-1 + bad*-2
	fmt.Println(code, name, cond_str, bad_cond_str, finance)
	seed := []int{1, 2, 3, 4, 5, 6}
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(seed)
	if score < 0 {
		score = 0
	}

	p := dal.Predict{Code: code, Name: name, Condition: cond_str, BadCondition: bad_cond_str, Finance: finance,
		Date: result.CurrDate, Score: score + seed[n], Price: result.RecentClose[0], Percent: result.RecentPercent[0],
		FundCount: jigouchicangcount, SMCount: simuchicangcount, FenghongCount: fenhong, PeiguCount: pergu, ZhuangzengCount: zhuangzeng,
		SongguCount: songgu, ZengfaCount: zengfa}
	if utils.TellEnv() == "loc" {
		err := store.MysqlClient.GetOnlineDB().Save(&p).Error
		if err != nil {
			fmt.Println("写入线上错误")
		}
	}
	store.MysqlClient.GetDB().Save(&p)
}
