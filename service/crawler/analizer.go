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

func GetSubCompByCode(code string) int {
	var count int
	store.MysqlClient.GetDB().Model(&dal.StockSubCompany{}).Where("code = ?", code).Count(&count)
	return count
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
		if math.Abs((recent_money_or_count[i]-recent_ave[i])/recent_ave[i]) < 0.003 {
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
	score := 0
	cond_str, bad_cond_str, finance := "", "", ""
	cond_str_, bad_cond_str_ := "", ""
	// 5 10 金叉  并且 今日10 均大于昨日 10均
	jincha1 := result.AveDailyPrice1[0] > result.AveDailyPrice2[0] && result.AveDailyPrice1[1] < result.AveDailyPrice2[1] && result.AveDailyPrice2[0] >= result.AveDailyPrice2[1]
	// 5 30 金叉
	jincha2 := result.AveDailyPrice1[0] > result.AveDailyPrice3[0] && result.AveDailyPrice1[1] < result.AveDailyPrice3[1] && result.AveDailyPrice3[0] >= result.AveDailyPrice3[1]
	// 5 60 金叉
	jincha3 := result.AveDailyPrice1[0] > result.AveDailyPrice4[0] && result.AveDailyPrice1[1] < result.AveDailyPrice4[1] && result.AveDailyPrice4[0] >= result.AveDailyPrice4[1]
	// 10 30 金叉
	jincha5 := result.AveDailyPrice2[0] > result.AveDailyPrice3[0] && result.AveDailyPrice2[1] < result.AveDailyPrice3[1] && result.AveDailyPrice3[0] >= result.AveDailyPrice3[1]
	// 10 60 金叉
	jincha6 := result.AveDailyPrice2[0] > result.AveDailyPrice4[0] && result.AveDailyPrice2[1] < result.AveDailyPrice4[1] && result.AveDailyPrice4[0] >= result.AveDailyPrice4[1]
	// 30 60 金叉
	jincha8 := result.AveDailyPrice3[0] > result.AveDailyPrice4[0] && result.AveDailyPrice3[1] < result.AveDailyPrice4[1] && result.AveDailyPrice4[0] >= result.AveDailyPrice4[1]
	// 10 40 量能金叉
	jincha11 := result.AveCount1[0] > result.AveCount2[0] && result.AveCount1[1] < result.AveCount2[1] && result.AveCount2[0] > result.AveCount2[1]
	// 死叉股 并且今日慢均线 比 昨日减小
	sicha1 := result.AveDailyPrice1[0] < result.AveDailyPrice2[0] && result.AveDailyPrice1[1] > result.AveDailyPrice2[1] && result.AveDailyPrice2[0] <= result.AveDailyPrice2[1]
	sicha2 := result.AveDailyPrice1[0] < result.AveDailyPrice3[0] && result.AveDailyPrice1[1] > result.AveDailyPrice3[1] && result.AveDailyPrice3[0] <= result.AveDailyPrice3[1]
	sicha3 := result.AveDailyPrice1[0] < result.AveDailyPrice4[0] && result.AveDailyPrice1[1] > result.AveDailyPrice4[1] && result.AveDailyPrice4[0] <= result.AveDailyPrice4[1]
	sicha5 := result.AveDailyPrice2[0] < result.AveDailyPrice3[0] && result.AveDailyPrice2[1] > result.AveDailyPrice3[1] && result.AveDailyPrice3[0] <= result.AveDailyPrice3[1]
	sicha6 := result.AveDailyPrice2[0] < result.AveDailyPrice4[0] && result.AveDailyPrice2[1] > result.AveDailyPrice4[1] && result.AveDailyPrice4[0] <= result.AveDailyPrice4[1]
	sicha8 := result.AveDailyPrice3[0] < result.AveDailyPrice4[0] && result.AveDailyPrice3[1] > result.AveDailyPrice4[1] && result.AveDailyPrice4[0] <= result.AveDailyPrice4[1]
	sicha11 := result.AveCount1[0] < result.AveCount2[0] && result.AveCount1[1] > result.AveCount2[1] && result.AveCount2[0] < result.AveCount2[1]

	// 涨停股
	zhangting := result.RecentPercent[0] > 9.94
	// 一字板
	yiziban := result.RecentPercent[0] > 9.94 && result.RecentHigh[0] == result.RecentLow[0]
	// T 字板
	tziban := result.RecentPercent[0] > 9.94 && result.RecentOpen[0] == result.RecentClose[0] && result.RecentClose[0] > result.RecentLow[0]
	// 一字跌停板
	dietingban := result.RecentPercent[0] < -9.94 && result.RecentHigh[0] == result.RecentLow[0]
	// 倒T板
	daotban := result.RecentPercent[0] < -9.94 && result.RecentOpen[0] == result.RecentClose[0] && result.RecentHigh[0] > result.RecentClose[0]

	// 4条均线 价格均线上扬
	priceshangyang1 := result.AveDailyPrice1[0] > result.AveDailyPrice1[1] && result.AveDailyPrice1[1] > result.AveDailyPrice1[2] && result.AveDailyPrice1[2] > result.AveDailyPrice1[3] && result.AveDailyPrice1[3] > result.AveDailyPrice1[4] && result.AveDailyPrice1[4] > result.AveDailyPrice1[5]
	priceshangyang2 := result.AveDailyPrice2[0] > result.AveDailyPrice2[1] && result.AveDailyPrice2[1] > result.AveDailyPrice2[2] && result.AveDailyPrice2[2] > result.AveDailyPrice2[3] && result.AveDailyPrice2[3] > result.AveDailyPrice2[4] && result.AveDailyPrice2[4] > result.AveDailyPrice2[5]
	priceshangyang3 := result.AveDailyPrice3[0] > result.AveDailyPrice3[1] && result.AveDailyPrice3[1] > result.AveDailyPrice3[2] && result.AveDailyPrice3[2] > result.AveDailyPrice3[3] && result.AveDailyPrice3[3] > result.AveDailyPrice3[4] && result.AveDailyPrice3[4] > result.AveDailyPrice3[5]
	priceshangyang4 := result.AveDailyPrice4[0] > result.AveDailyPrice4[1] && result.AveDailyPrice4[1] > result.AveDailyPrice4[2] && result.AveDailyPrice4[2] > result.AveDailyPrice4[3] && result.AveDailyPrice4[3] > result.AveDailyPrice4[4] && result.AveDailyPrice4[4] > result.AveDailyPrice4[5]
	// 4条均线 价格均线下降
	pricexiajiang1 := result.AveDailyPrice1[0] < result.AveDailyPrice1[1] && result.AveDailyPrice1[1] < result.AveDailyPrice1[2] && result.AveDailyPrice1[2] < result.AveDailyPrice1[3] && result.AveDailyPrice1[3] < result.AveDailyPrice1[4] && result.AveDailyPrice1[4] < result.AveDailyPrice1[5]
	pricexiajiang2 := result.AveDailyPrice2[0] < result.AveDailyPrice2[1] && result.AveDailyPrice2[1] < result.AveDailyPrice2[2] && result.AveDailyPrice2[2] < result.AveDailyPrice2[3] && result.AveDailyPrice2[3] < result.AveDailyPrice2[4] && result.AveDailyPrice2[4] < result.AveDailyPrice2[5]
	pricexiajiang3 := result.AveDailyPrice3[0] < result.AveDailyPrice3[1] && result.AveDailyPrice3[1] < result.AveDailyPrice3[2] && result.AveDailyPrice3[2] < result.AveDailyPrice3[3] && result.AveDailyPrice3[3] < result.AveDailyPrice3[4] && result.AveDailyPrice3[4] < result.AveDailyPrice3[5]
	pricexiajiang4 := result.AveDailyPrice4[0] < result.AveDailyPrice4[1] && result.AveDailyPrice4[1] < result.AveDailyPrice4[2] && result.AveDailyPrice4[2] < result.AveDailyPrice4[3] && result.AveDailyPrice4[3] < result.AveDailyPrice4[4] && result.AveDailyPrice4[4] < result.AveDailyPrice4[5]

	// 当前价格在短期均线上方 （取非为小于）
	priceaboveave1 := result.RecentClose[0] >= result.AveDailyPrice1[0]
	priceaboveave2 := result.RecentClose[0] >= result.AveDailyPrice2[0]
	priceaboveave3 := result.RecentClose[0] >= result.AveDailyPrice3[0]
	priceaboveave4 := result.RecentClose[0] >= result.AveDailyPrice4[0]

	//// 均价粘合
	//junjialianhe1 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice1, 5, 4)
	//junjialianhe2 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice2, 5, 4)
	//junjialianhe3 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice3, 5, 4)
	//junjialianhe4 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice4, 5, 4)

	// 低开高走
	dikaigaozou := (result.RecentClose[1]-result.RecentOpen[0])/result.RecentClose[1] > 0.02 && result.RecentPercent[0] > 3
	// 高开低走
	gaokaidizou := (result.RecentOpen[0]-result.RecentClose[1])/result.RecentOpen[0] > 0.02 && result.RecentPercent[0] < -3
	// 低开低走
	dikaidizou := (result.RecentClose[1]-result.RecentOpen[0])/result.RecentClose[1] > 0.02 && result.RecentPercent[0] < -5
	// 高开高走
	gaokaigaozou := (result.RecentOpen[0]-result.RecentClose[1])/result.RecentOpen[0] > 0.02 && result.RecentPercent[0] > 5

	// 3 连阳
	sanlianyang := result.RecentPercent[0] > 0 && result.RecentPercent[1] > 0 && result.RecentPercent[2] > 0
	// 4 连阳
	silianyang := sanlianyang && result.RecentPercent[3] > 0
	// 5连阳
	wulianyang := silianyang && result.RecentPercent[4] > 0

	// 3连阴
	sanlianyin := result.RecentPercent[0] < 0 && result.RecentPercent[1] < 0 && result.RecentPercent[2] < 0
	// 4连阴
	silianyin := sanlianyin && result.RecentPercent[3] < 0
	// 5连阴
	wulianyin := silianyin && result.RecentPercent[4] < 0

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
	tufangjuliang := (result.RecentCount[0]-result.RecentCount[1])/result.RecentCount[1] > 3 || (result.RecentCount[1]-result.RecentCount[2])/result.RecentCount[1] > 3
	// 量能突破均线
	liangnengtupo1 := result.AveCount1[0] < result.RecentCount[0] && result.AveCount1[1] > result.RecentCount[1] && result.AveCount1[2] > result.RecentCount[2] && result.AveCount1[3] > result.RecentCount[3] && result.AveCount1[4] > result.RecentCount[4] && result.AveCount1[5] > result.RecentCount[5]
	liangnengtupo2 := result.AveCount2[0] < result.RecentCount[0] && result.AveCount2[1] > result.RecentCount[1] && result.AveCount2[2] > result.RecentCount[2] && result.AveCount2[3] > result.RecentCount[3] && result.AveCount2[4] > result.RecentCount[4] && result.AveCount2[5] > result.RecentCount[5]

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
	subcomp := GetSubCompByCode(code)
	// 历史更名次数
	changename, has_st := GetHistoryNameByCode(code)

	// 周线数据如下
	// 周线金叉
	var wjincha1, wjincha2, wjincha3, wjincha5, wjincha6, wjincha8, wjincha11 bool
	var wsicha1, wsicha2, wsicha3, wsicha5, wsicha6, wsicha8, wsicha11 bool
	var wpriceshangyang1, wpriceshangyang2, wpriceshangyang3, wpriceshangyang4 bool
	var wpricexiajiang1, wpricexiajiang2, wpricexiajiang3, wpricexiajiang4 bool
	var wpriceaboveave1, wpriceaboveave2, wpriceaboveave3, wpriceaboveave4, wliangshangyang1, wliangshangyang2, wliangnengbuduanbigger, wtufangjuliang, wliangnengtupo1, wliangnengtupo2 bool
	var tupo5, tupo6, tupo7, tupo8, jichuang5, jichuang6, jichuang7, jichuang8 bool // 突破5,10,30,60yue

	if result.AveWeeklyPrice1 != nil && result.AveWeeklyPrice2 != nil { // 两条短均线一定不能为nil
		// 周线5 10 金叉
		wjincha1 = result.AveWeeklyPrice1[0] > result.AveWeeklyPrice2[0] && result.AveWeeklyPrice1[1] < result.AveWeeklyPrice2[1] && result.AveWeeklyPrice2[0] >= result.AveWeeklyPrice2[1]
		// 周线死叉股 5*10
		wsicha1 = result.AveWeeklyPrice1[0] < result.AveWeeklyPrice2[0] && result.AveWeeklyPrice1[1] > result.AveWeeklyPrice2[1] && result.AveWeeklyPrice2[0] <= result.AveWeeklyPrice2[1]
		// 5，10周线上扬
		wpriceshangyang1 = result.AveWeeklyPrice1[0] > result.AveWeeklyPrice1[1] && result.AveWeeklyPrice1[1] > result.AveWeeklyPrice1[2] && result.AveWeeklyPrice1[2] > result.AveWeeklyPrice1[3] && result.AveWeeklyPrice1[3] > result.AveWeeklyPrice1[4] && result.AveWeeklyPrice1[4] > result.AveWeeklyPrice1[5]
		wpriceshangyang2 = result.AveWeeklyPrice2[0] > result.AveWeeklyPrice2[1] && result.AveWeeklyPrice2[1] > result.AveWeeklyPrice2[2] && result.AveWeeklyPrice2[2] > result.AveWeeklyPrice2[3] && result.AveWeeklyPrice2[3] > result.AveWeeklyPrice2[4] && result.AveWeeklyPrice2[4] > result.AveWeeklyPrice2[5]
		// 5，10周线下降
		wpricexiajiang1 = result.AveWeeklyPrice1[0] < result.AveWeeklyPrice1[1] && result.AveWeeklyPrice1[1] < result.AveWeeklyPrice1[2] && result.AveWeeklyPrice1[2] < result.AveWeeklyPrice1[3] && result.AveWeeklyPrice1[3] < result.AveWeeklyPrice1[4] && result.AveWeeklyPrice1[4] < result.AveWeeklyPrice1[5]
		wpricexiajiang2 = result.AveWeeklyPrice2[0] < result.AveWeeklyPrice2[1] && result.AveWeeklyPrice2[1] < result.AveWeeklyPrice2[2] && result.AveWeeklyPrice2[2] < result.AveWeeklyPrice2[3] && result.AveWeeklyPrice2[3] < result.AveWeeklyPrice2[4] && result.AveWeeklyPrice2[4] < result.AveWeeklyPrice2[5]
		// 周线当前价格在短期均线上方 （取非为小于）
		wpriceaboveave1 = result.RecentCloseWeek[0] >= result.AveWeeklyPrice1[0]
		wpriceaboveave2 = result.RecentCloseWeek[0] >= result.AveWeeklyPrice2[0]
		tupo5 = result.AveWeeklyPrice1[0] < result.RecentClose[0] && result.AveWeeklyPrice1[1] >= result.RecentClose[1] && result.AveWeeklyPrice1[2] > result.RecentClose[2] && result.AveWeeklyPrice1[3] > result.RecentClose[3]
		tupo6 = result.AveWeeklyPrice2[0] < result.RecentClose[0] && result.AveWeeklyPrice2[1] >= result.RecentClose[1] && result.AveWeeklyPrice2[2] > result.RecentClose[2] && result.AveWeeklyPrice2[3] > result.RecentClose[3]
		jichuang5 = result.AveWeeklyPrice1[0] > result.RecentClose[0] && result.AveWeeklyPrice1[1] <= result.RecentClose[1] && result.AveWeeklyPrice1[2] < result.RecentClose[2] && result.AveWeeklyPrice1[3] < result.RecentClose[3]
		jichuang6 = result.AveWeeklyPrice2[0] > result.RecentClose[0] && result.AveWeeklyPrice2[1] <= result.RecentClose[1] && result.AveWeeklyPrice2[2] < result.RecentClose[2] && result.AveWeeklyPrice2[3] < result.RecentClose[3]

		if result.AveWeeklyPrice3 != nil {
			tupo7 = result.AveWeeklyPrice3[0] < result.RecentClose[0] && result.AveWeeklyPrice3[1] >= result.RecentClose[1] && result.AveWeeklyPrice3[2] > result.RecentClose[2] && result.AveWeeklyPrice3[3] > result.RecentClose[3]
			jichuang7 = result.AveWeeklyPrice3[0] > result.RecentClose[0] && result.AveWeeklyPrice3[1] <= result.RecentClose[1] && result.AveWeeklyPrice3[2] < result.RecentClose[2] && result.AveWeeklyPrice3[3] < result.RecentClose[3]

			wjincha2 = result.AveWeeklyPrice1[0] > result.AveWeeklyPrice3[0] && result.AveWeeklyPrice1[1] < result.AveWeeklyPrice3[1] && result.AveWeeklyPrice3[0] >= result.AveWeeklyPrice3[1]
			wjincha5 = result.AveWeeklyPrice2[0] > result.AveWeeklyPrice3[0] && result.AveWeeklyPrice2[1] < result.AveWeeklyPrice3[1] && result.AveWeeklyPrice3[0] >= result.AveWeeklyPrice3[1]
			wsicha2 = result.AveWeeklyPrice1[0] < result.AveWeeklyPrice3[0] && result.AveWeeklyPrice1[1] > result.AveWeeklyPrice3[1] && result.AveWeeklyPrice3[0] <= result.AveWeeklyPrice3[1]
			wsicha5 = result.AveWeeklyPrice2[0] < result.AveWeeklyPrice3[0] && result.AveWeeklyPrice2[1] > result.AveWeeklyPrice3[1] && result.AveWeeklyPrice3[0] <= result.AveWeeklyPrice3[1]
			wpriceshangyang3 = result.AveWeeklyPrice3[0] > result.AveWeeklyPrice3[1] && result.AveWeeklyPrice3[1] > result.AveWeeklyPrice3[2] && result.AveWeeklyPrice3[2] > result.AveWeeklyPrice3[3] && result.AveWeeklyPrice3[3] > result.AveWeeklyPrice3[4] && result.AveWeeklyPrice3[4] > result.AveWeeklyPrice3[5]
			wpricexiajiang3 = result.AveWeeklyPrice3[0] < result.AveWeeklyPrice3[1] && result.AveWeeklyPrice3[1] < result.AveWeeklyPrice3[2] && result.AveWeeklyPrice3[2] < result.AveWeeklyPrice3[3] && result.AveWeeklyPrice3[3] < result.AveWeeklyPrice3[4] && result.AveWeeklyPrice3[4] < result.AveWeeklyPrice3[5]
			wpriceaboveave3 = result.RecentCloseWeek[0] >= result.AveWeeklyPrice3[0]
			if result.AveWeeklyPrice4 != nil {
				tupo8 = result.AveWeeklyPrice4[0] < result.RecentClose[0] && result.AveWeeklyPrice4[1] >= result.RecentClose[1] && result.AveWeeklyPrice4[2] > result.RecentClose[2] && result.AveWeeklyPrice4[3] > result.RecentClose[3]
				jichuang8 = result.AveWeeklyPrice4[0] > result.RecentClose[0] && result.AveWeeklyPrice4[1] <= result.RecentClose[1] && result.AveWeeklyPrice4[2] < result.RecentClose[2] && result.AveWeeklyPrice4[3] < result.RecentClose[3]
				wjincha8 = result.AveWeeklyPrice3[0] > result.AveWeeklyPrice4[0] && result.AveWeeklyPrice3[1] < result.AveWeeklyPrice4[1] && result.AveWeeklyPrice4[0] >= result.AveWeeklyPrice4[1]
				wsicha8 = result.AveWeeklyPrice3[0] < result.AveWeeklyPrice4[0] && result.AveWeeklyPrice3[1] > result.AveWeeklyPrice4[1] && result.AveWeeklyPrice4[0] <= result.AveWeeklyPrice4[1]
				wsicha3 = result.AveWeeklyPrice1[0] < result.AveWeeklyPrice4[0] && result.AveWeeklyPrice1[1] > result.AveWeeklyPrice4[1] && result.AveWeeklyPrice4[0] <= result.AveWeeklyPrice3[1]
				wjincha3 = result.AveWeeklyPrice1[0] > result.AveWeeklyPrice4[0] && result.AveWeeklyPrice1[1] < result.AveWeeklyPrice4[1] && result.AveWeeklyPrice4[0] >= result.AveWeeklyPrice4[1]
				wjincha6 = result.AveWeeklyPrice2[0] > result.AveWeeklyPrice4[0] && result.AveWeeklyPrice2[1] < result.AveWeeklyPrice4[1] && result.AveWeeklyPrice4[0] >= result.AveWeeklyPrice4[1]
				wsicha6 = result.AveWeeklyPrice2[0] < result.AveWeeklyPrice4[0] && result.AveWeeklyPrice2[1] > result.AveWeeklyPrice4[1] && result.AveWeeklyPrice4[0] <= result.AveWeeklyPrice4[1]
				wpriceshangyang4 = result.AveWeeklyPrice4[0] > result.AveWeeklyPrice4[1] && result.AveWeeklyPrice4[1] > result.AveWeeklyPrice4[2] && result.AveWeeklyPrice4[2] > result.AveWeeklyPrice4[3] && result.AveWeeklyPrice4[3] > result.AveWeeklyPrice4[4] && result.AveWeeklyPrice4[4] > result.AveWeeklyPrice4[5]
				wpricexiajiang4 = result.AveWeeklyPrice4[0] < result.AveWeeklyPrice4[1] && result.AveWeeklyPrice4[1] < result.AveWeeklyPrice4[2] && result.AveWeeklyPrice4[2] < result.AveWeeklyPrice4[3] && result.AveWeeklyPrice4[3] < result.AveWeeklyPrice4[4] && result.AveWeeklyPrice4[4] < result.AveWeeklyPrice4[5]
				wpriceaboveave4 = result.RecentCloseWeek[0] >= result.AveWeeklyPrice4[0]
			}
		}
		if result.AveCountWeekly1 != nil && result.AveCountWeekly2 != nil {
			// 周线量能金叉 10x40
			wjincha11 = result.AveCountWeekly1[0] > result.AveCountWeekly2[0] && result.AveCountWeekly1[1] < result.AveCountWeekly2[1] && result.AveCountWeekly2[0] > result.AveCountWeekly2[1]
			// 周线量能死叉 10x40
			wsicha11 = result.AveCountWeekly1[0] < result.AveCountWeekly2[0] && result.AveCountWeekly1[1] > result.AveCountWeekly2[1] && result.AveCountWeekly2[0] < result.AveCountWeekly2[1]
			// 连续5周量能10均线上扬
			wliangshangyang1 = result.AveCountWeekly1[0] > result.AveCountWeekly1[1] && result.AveCountWeekly1[1] > result.AveCountWeekly1[2] && result.AveCountWeekly1[2] > result.AveCountWeekly1[3] && result.AveCountWeekly1[3] > result.AveCountWeekly1[4] && result.AveCountWeekly1[4] > result.AveCountWeekly1[5]
			// 连续5周量能40均线上扬
			wliangshangyang2 = result.AveCountWeekly2[0] > result.AveCountWeekly2[1] && result.AveCountWeekly2[1] > result.AveCountWeekly2[2] && result.AveCountWeekly2[2] > result.AveCountWeekly2[3] && result.AveCountWeekly2[3] > result.AveCountWeekly2[4] && result.AveCountWeekly2[4] > result.AveCountWeekly2[5]
			// 量能不断放大
			wliangnengbuduanbigger = result.RecentCountWeek[0] > result.RecentCountWeek[1] && result.RecentCountWeek[1] > result.RecentCountWeek[2] && result.RecentCountWeek[2] > result.RecentCountWeek[3]
			// 突放巨量
			wtufangjuliang = (result.RecentCountWeek[0]-result.RecentCountWeek[1])/result.RecentCountWeek[1] > 3 || (result.RecentCountWeek[1]-result.RecentCountWeek[2])/result.RecentCountWeek[1] > 3
			// 量能突破均线
			wliangnengtupo1 = result.AveCountWeekly1[0] < result.RecentCountWeek[0] && result.AveCountWeekly1[1] > result.RecentCountWeek[1] && result.AveCountWeekly1[2] > result.RecentCountWeek[2] && result.AveCountWeekly1[3] > result.RecentCountWeek[3] && result.AveCountWeekly1[4] > result.RecentCountWeek[4] && result.AveCountWeekly1[5] > result.RecentCountWeek[5]
			wliangnengtupo2 = result.AveCountWeekly2[0] < result.RecentCountWeek[0] && result.AveCountWeekly2[1] > result.RecentCountWeek[1] && result.AveCountWeekly2[2] > result.RecentCountWeek[2] && result.AveCountWeekly2[3] > result.RecentCountWeek[3] && result.AveCountWeekly2[4] > result.RecentCountWeek[4] && result.AveCountWeekly2[5] > result.RecentCountWeek[5]
		}
	}

	// 月线数据如下
	// 月线金叉
	var yjincha1, yjincha2, yjincha3, yjincha5, yjincha6, yjincha8, yjincha11 bool
	var ysicha1, ysicha2, ysicha3, ysicha5, ysicha6, ysicha8, ysicha11 bool
	var ypriceshangyang1, ypriceshangyang2, ypriceshangyang3, ypriceshangyang4 bool
	var ypricexiajiang1, ypricexiajiang2, ypricexiajiang3, ypricexiajiang4 bool
	var ypriceaboveave1, ypriceaboveave2, ypriceaboveave3, ypriceaboveave4, yliangshangyang1, yliangshangyang2, yliangnengbuduanbigger, ytufangjuliang, yliangnengtupo1, yliangnengtupo2 bool
	var tupo9, tupo10, tupo11, tupo12, jichuang9, jichuang10, jichuang11, jichuang12 bool // 突破5,10,30,60yue
	// 5 10 金叉  并且 今日10 均大于昨日 10均
	if result.AveMonthPrice1 != nil && result.AveMonthPrice2 != nil { // 两条短均线一定不能为nil
		yjincha1 = result.AveMonthPrice1[0] > result.AveMonthPrice2[0] && result.AveMonthPrice1[1] < result.AveMonthPrice2[1] && result.AveMonthPrice2[0] >= result.AveMonthPrice2[1]
		ysicha1 = result.AveMonthPrice1[0] < result.AveMonthPrice2[0] && result.AveMonthPrice1[1] > result.AveMonthPrice2[1] && result.AveMonthPrice2[0] <= result.AveMonthPrice2[1]
		ypriceshangyang1 = result.AveMonthPrice1[0] > result.AveMonthPrice1[1] && result.AveMonthPrice1[1] > result.AveMonthPrice1[2] && result.AveMonthPrice1[2] > result.AveMonthPrice1[3] && result.AveMonthPrice1[3] > result.AveMonthPrice1[4] && result.AveMonthPrice1[4] > result.AveMonthPrice1[5]
		ypriceshangyang2 = result.AveMonthPrice2[0] > result.AveMonthPrice2[1] && result.AveMonthPrice2[1] > result.AveMonthPrice2[2] && result.AveMonthPrice2[2] > result.AveMonthPrice2[3] && result.AveMonthPrice2[3] > result.AveMonthPrice2[4] && result.AveMonthPrice2[4] > result.AveMonthPrice2[5]
		ypricexiajiang1 = result.AveMonthPrice1[0] < result.AveMonthPrice1[1] && result.AveMonthPrice1[1] < result.AveMonthPrice1[2] && result.AveMonthPrice1[2] < result.AveMonthPrice1[3] && result.AveMonthPrice1[3] < result.AveMonthPrice1[4] && result.AveMonthPrice1[4] < result.AveMonthPrice1[5]
		ypricexiajiang2 = result.AveMonthPrice2[0] < result.AveMonthPrice2[1] && result.AveMonthPrice2[1] < result.AveMonthPrice2[2] && result.AveMonthPrice2[2] < result.AveMonthPrice2[3] && result.AveMonthPrice2[3] < result.AveMonthPrice2[4] && result.AveMonthPrice2[4] < result.AveMonthPrice2[5]
		ypriceaboveave1 = result.RecentCloseMonth[0] >= result.AveMonthPrice1[0]
		ypriceaboveave2 = result.RecentCloseMonth[0] >= result.AveMonthPrice2[0]
		tupo9 = result.AveMonthPrice1[0] < result.RecentClose[0] && result.AveMonthPrice1[0] >= result.RecentClose[1] && result.AveMonthPrice1[2] > result.RecentClose[2] && result.AveMonthPrice1[3] > result.RecentClose[3]
		tupo10 = result.AveMonthPrice2[0] < result.RecentClose[0] && result.AveMonthPrice2[0] >= result.RecentClose[1] && result.AveMonthPrice2[2] > result.RecentClose[2] && result.AveMonthPrice2[3] > result.RecentClose[3]
		jichuang9 = result.AveMonthPrice1[0] > result.RecentClose[0] && result.AveMonthPrice1[0] <= result.RecentClose[1] && result.AveMonthPrice1[2] < result.RecentClose[2] && result.AveMonthPrice1[3] < result.RecentClose[3]
		jichuang10 = result.AveMonthPrice2[0] > result.RecentClose[0] && result.AveMonthPrice2[0] <= result.RecentClose[1] && result.AveMonthPrice2[2] < result.RecentClose[2] && result.AveMonthPrice2[3] < result.RecentClose[3]
		if result.AveMonthPrice3 != nil {
			yjincha2 = result.AveMonthPrice1[0] > result.AveMonthPrice3[0] && result.AveMonthPrice1[1] < result.AveMonthPrice3[1] && result.AveMonthPrice3[0] >= result.AveMonthPrice3[1]
			yjincha5 = result.AveMonthPrice2[0] > result.AveMonthPrice3[0] && result.AveMonthPrice2[1] < result.AveMonthPrice3[1] && result.AveMonthPrice3[0] >= result.AveMonthPrice3[1]
			ysicha2 = result.AveMonthPrice1[0] < result.AveMonthPrice3[0] && result.AveMonthPrice1[1] > result.AveMonthPrice3[1] && result.AveMonthPrice1[0] <= result.AveMonthPrice1[1]
			ysicha5 = result.AveMonthPrice2[0] < result.AveMonthPrice3[0] && result.AveMonthPrice2[1] > result.AveMonthPrice3[1] && result.AveMonthPrice3[0] <= result.AveMonthPrice3[1]
			ypriceshangyang3 = result.AveMonthPrice3[0] > result.AveMonthPrice3[1] && result.AveMonthPrice3[1] > result.AveMonthPrice3[2] && result.AveMonthPrice3[2] > result.AveMonthPrice3[3] && result.AveMonthPrice3[3] > result.AveMonthPrice3[4] && result.AveMonthPrice3[4] > result.AveMonthPrice3[5]
			ypricexiajiang3 = result.AveMonthPrice3[0] < result.AveMonthPrice3[1] && result.AveMonthPrice3[1] < result.AveMonthPrice3[2] && result.AveMonthPrice3[2] < result.AveMonthPrice3[3] && result.AveMonthPrice3[3] < result.AveMonthPrice3[4] && result.AveMonthPrice3[4] < result.AveMonthPrice3[5]
			ypriceaboveave3 = result.RecentCloseMonth[0] >= result.AveMonthPrice3[0]
			tupo11 = result.AveMonthPrice3[0] < result.RecentClose[0] && result.AveMonthPrice3[0] >= result.RecentClose[1] && result.AveMonthPrice3[2] > result.RecentClose[2] && result.AveMonthPrice3[3] > result.RecentClose[3]
			jichuang11 = result.AveMonthPrice3[0] > result.RecentClose[0] && result.AveMonthPrice3[0] <= result.RecentClose[1] && result.AveMonthPrice3[2] < result.RecentClose[2] && result.AveMonthPrice3[3] < result.RecentClose[3]

			if result.AveMonthPrice4 != nil {
				yjincha8 = result.AveMonthPrice3[0] > result.AveMonthPrice4[0] && result.AveMonthPrice3[1] < result.AveMonthPrice4[1] && result.AveMonthPrice4[0] >= result.AveMonthPrice4[1]
				ysicha8 = result.AveMonthPrice4[0] < result.AveMonthPrice4[0] && result.AveMonthPrice3[1] > result.AveMonthPrice4[1] && result.AveMonthPrice4[0] <= result.AveMonthPrice4[1]
				ysicha3 = result.AveMonthPrice1[0] < result.AveMonthPrice4[0] && result.AveMonthPrice1[1] > result.AveMonthPrice4[1] && result.AveMonthPrice4[0] <= result.AveMonthPrice3[1]
				yjincha3 = result.AveMonthPrice1[0] > result.AveMonthPrice4[0] && result.AveMonthPrice1[1] < result.AveMonthPrice4[1] && result.AveMonthPrice4[0] >= result.AveMonthPrice4[1]
				yjincha6 = result.AveMonthPrice2[0] > result.AveMonthPrice4[0] && result.AveMonthPrice2[1] < result.AveMonthPrice4[1] && result.AveMonthPrice4[0] >= result.AveMonthPrice4[1]
				ysicha6 = result.AveMonthPrice2[0] < result.AveMonthPrice4[0] && result.AveMonthPrice2[1] > result.AveMonthPrice4[1] && result.AveMonthPrice4[0] <= result.AveMonthPrice4[1]
				ypriceshangyang4 = result.AveMonthPrice4[0] > result.AveMonthPrice4[1] && result.AveMonthPrice4[1] > result.AveMonthPrice4[2] && result.AveMonthPrice4[2] > result.AveMonthPrice4[3] && result.AveMonthPrice4[3] > result.AveMonthPrice4[4] && result.AveMonthPrice4[4] > result.AveMonthPrice4[5]
				ypricexiajiang4 = result.AveMonthPrice4[0] < result.AveMonthPrice4[1] && result.AveMonthPrice4[1] < result.AveMonthPrice4[2] && result.AveMonthPrice4[2] < result.AveMonthPrice4[3] && result.AveMonthPrice4[3] < result.AveMonthPrice4[4] && result.AveMonthPrice4[4] < result.AveMonthPrice4[5]
				ypriceaboveave4 = result.RecentCloseMonth[0] >= result.AveMonthPrice4[0]
				tupo12 = result.AveMonthPrice4[0] < result.RecentClose[0] && result.AveMonthPrice4[0] >= result.RecentClose[1] && result.AveMonthPrice4[2] > result.RecentClose[2] && result.AveMonthPrice4[3] > result.RecentClose[3]
				jichuang12 = result.AveMonthPrice4[0] > result.RecentClose[0] && result.AveMonthPrice4[0] <= result.RecentClose[1] && result.AveMonthPrice4[2] < result.RecentClose[2] && result.AveMonthPrice4[3] < result.RecentClose[3]

			}
		}
		if result.AveCountMonth1 != nil && result.AveCountMonth2 != nil {
			// 量能金叉 前提是长均线要是上升状态
			yjincha11 = result.AveMonthPrice1[0] > result.AveMonthPrice2[0] && result.AveMonthPrice1[1] < result.AveMonthPrice2[1] && result.AveMonthPrice2[0] > result.AveMonthPrice2[1]
			// 量能死叉 前提是长均线要是下降状态
			ysicha11 = result.AveCountMonth1[0] < result.AveCountMonth2[0] && result.AveCountMonth1[1] > result.AveCountMonth2[1] && result.AveCountMonth2[0] < result.AveCountMonth2[1]
			yliangshangyang1 = result.AveCountMonth1[0] > result.AveCountMonth1[1] && result.AveCountMonth1[1] > result.AveCountMonth1[2] && result.AveCountMonth1[2] > result.AveCountMonth1[3] && result.AveCountMonth1[3] > result.AveCountMonth1[4] && result.AveCountMonth1[4] > result.AveCountMonth1[5]
			yliangshangyang2 = result.AveCountMonth2[0] > result.AveCountMonth2[1] && result.AveCountMonth2[1] > result.AveCountMonth2[2] && result.AveCountMonth2[2] > result.AveCountMonth2[3] && result.AveCountMonth2[3] > result.AveCountMonth2[4] && result.AveCountMonth2[4] > result.AveCountMonth2[5]
			yliangnengbuduanbigger = result.RecentCountMonth[0] > result.RecentCountMonth[1] && result.RecentCountMonth[1] > result.RecentCountMonth[2] && result.RecentCountMonth[2] > result.RecentCountMonth[3]
			ytufangjuliang = (result.RecentCountMonth[0]-result.RecentCountMonth[1])/result.RecentCountMonth[1] > 3 || (result.RecentCountMonth[1]-result.RecentCountMonth[2])/result.RecentCountMonth[1] > 3
			yliangnengtupo1 = result.AveCountMonth1[0] < result.RecentCountMonth[0] && result.AveCountMonth1[1] > result.RecentCountMonth[1] && result.AveCountMonth1[2] > result.RecentCountMonth[2] && result.AveCountMonth1[3] > result.RecentCountMonth[3] && result.AveCountMonth1[4] > result.RecentCountMonth[4] && result.AveCountMonth1[5] > result.RecentCountMonth[5]
			yliangnengtupo2 = result.AveCountMonth2[0] < result.RecentCountMonth[0] && result.AveCountMonth2[1] > result.RecentCountMonth[1] && result.AveCountMonth2[2] > result.RecentCountMonth[2] && result.AveCountMonth2[3] > result.RecentCountMonth[3] && result.AveCountMonth2[4] > result.RecentCountMonth[4] && result.AveCountMonth2[5] > result.RecentCountMonth[5]
		}
	}
	// 价格突破5日压力位
	tupo1 := result.AveDailyPrice1[0] < result.RecentClose[0] && result.AveDailyPrice1[1] >= result.RecentClose[1] && result.AveDailyPrice1[2] > result.RecentClose[2]
	tupo2 := result.AveDailyPrice2[0] < result.RecentClose[0] && result.AveDailyPrice2[1] >= result.RecentClose[1] && result.AveDailyPrice2[2] > result.RecentClose[2]
	tupo3 := result.AveDailyPrice3[0] < result.RecentClose[0] && result.AveDailyPrice3[1] >= result.RecentClose[1] && result.AveDailyPrice3[2] > result.RecentClose[2]
	tupo4 := result.AveDailyPrice4[0] < result.RecentClose[0] && result.AveDailyPrice4[1] >= result.RecentClose[1] && result.AveDailyPrice4[2] > result.RecentClose[2]

	jichuang1 := result.AveDailyPrice1[0] > result.RecentClose[0] && result.AveDailyPrice1[1] <= result.RecentClose[1] && result.AveDailyPrice1[2] < result.RecentClose[2]
	jichuang2 := result.AveDailyPrice2[0] > result.RecentClose[0] && result.AveDailyPrice2[1] <= result.RecentClose[1] && result.AveDailyPrice2[2] < result.RecentClose[2]
	jichuang3 := result.AveDailyPrice3[0] > result.RecentClose[0] && result.AveDailyPrice3[1] <= result.RecentClose[1] && result.AveDailyPrice3[2] < result.RecentClose[2]
	jichuang4 := result.AveDailyPrice4[0] > result.RecentClose[0] && result.AveDailyPrice4[1] <= result.RecentClose[1] && result.AveDailyPrice4[2] < result.RecentClose[2]

	if jincha11 {
		score += 5
		cond_str += "(量能金叉)10日上穿40日均线; "
		cond_str_ += "(量能金叉)10日上穿40日均线; "
	}
	if jincha1 {
		score += 5
		cond_str += "(金叉)5日上穿10日均线; "
		cond_str_ += "(金叉)5日上穿10日均线; "
	}
	if jincha2 {
		score += 5
		cond_str += "(金叉)5日上穿30日均线; "
		cond_str_ += "(金叉)5日上穿30日均线; "
	}
	if jincha3 {
		score += 5
		cond_str += "(金叉)5日上穿60日均线; "
		cond_str_ += "(金叉)5日上穿60日均线; "
	}
	if jincha5 {
		score += 5
		cond_str += "(金叉)10日上穿30日均线; "
		cond_str_ += "(金叉)10日上穿30日均线; "
	}
	if jincha6 {
		score += 5
		cond_str += "(金叉)10日上穿60日均线; "
		cond_str_ += "(金叉)10日上穿60日均线; "
	}
	if jincha8 {
		score += 5
		cond_str += "(金叉)30日上穿60日均线; "
		cond_str_ += "(金叉)30日上穿60日均线; "
	}

	if wjincha11 {
		score += 5
		cond_str += "(量能金叉)10周上穿40周均线; "
		cond_str_ += "(量能金叉)10周上穿40周均线; "
	}
	if wjincha1 {
		score += 5
		cond_str += "(金叉)5周上穿10周均线; "
		cond_str_ += "(金叉)5周上穿10周均线; "
	}
	if wjincha2 {
		score += 5
		cond_str += "(金叉)5周上穿30周均线; "
		cond_str_ += "(金叉)5周上穿30周均线; "
	}
	if wjincha3 {
		score += 5
		cond_str += "(金叉)5周上穿60周均线; "
		cond_str_ += "(金叉)5周上穿60周均线; "
	}
	if wjincha5 {
		score += 5
		cond_str += "(金叉)10周上穿30周均线; "
		cond_str_ += "(金叉)10周上穿30周均线; "
	}
	if wjincha6 {
		score += 5
		cond_str += "(金叉)10周上穿60周均线; "
		cond_str_ += "(金叉)10周上穿60周均线; "
	}
	if wjincha8 {
		score += 5
		cond_str += "(金叉)30周上穿60周均线; "
		cond_str_ += "(金叉)30周上穿60周均线; "
	}

	if yjincha11 {
		score += 5
		cond_str += "(量能金叉)10月上穿40月均线; "
		cond_str_ += "(量能金叉)10月上穿40月均线; "
	}
	if yjincha1 {
		score += 5
		cond_str += "(金叉)5月上穿10月周均线; "
		cond_str_ += "(金叉)5月上穿10月周均线; "
	}
	if yjincha2 {
		score += 5
		cond_str += "(金叉)5月上穿30月均线; "
		cond_str_ += "(金叉)5月上穿30月均线; "
	}
	if yjincha3 {
		score += 5
		cond_str += "(金叉)5月上穿60月均线; "
		cond_str_ += "(金叉)5月上穿60月均线; "
	}
	if yjincha5 {
		score += 5
		cond_str += "(金叉)10月上穿30月均线; "
		cond_str_ += "(金叉)10月上穿30月均线; "
	}
	if yjincha6 {
		score += 5
		cond_str += "(金叉)10月上穿60月均线; "
		cond_str_ += "(金叉)10月上穿60月均线; "
	}
	if yjincha8 {
		score += 5
		cond_str += "(金叉)30月上穿60月均线; "
		cond_str_ += "(金叉)30月上穿60月均线; "
	}

	if sicha11 {
		score -= 5
		bad_cond_str += "(量能死叉)10日下穿40日均线; "
		bad_cond_str_ += "(量能死叉)10日下穿40日均线; "
	}
	if sicha1 {
		score -= 5
		bad_cond_str += "(死叉)5日下穿10日均线; "
		bad_cond_str_ += "(死叉)5日下穿10日均线; "
	}
	if sicha2 {
		score -= 5
		bad_cond_str += "(死叉)5日下穿30日均线; "
		bad_cond_str_ += "(死叉)5日下穿30日均线; "
	}
	if sicha3 {
		score -= 5
		bad_cond_str += "(死叉)5日下穿60日均线; "
		bad_cond_str_ += "(死叉)5日下穿60日均线; "
	}
	if sicha5 {
		score -= 5
		bad_cond_str += "(死叉)10日下穿30日均线; "
		bad_cond_str_ += "(死叉)10日下穿30日均线; "
	}
	if sicha6 {
		score -= 5
		bad_cond_str += "(死叉)10日下穿60日均线; "
		bad_cond_str_ += "(死叉)10日下穿60日均线; "
	}
	if sicha8 {
		score -= 5
		bad_cond_str += "(死叉)30下穿60日均线; "
		bad_cond_str_ += "(死叉)30下穿60日均线; "
	}

	if wsicha11 {
		score -= 5
		bad_cond_str += "(量能死叉)10周下穿40周均线; "
		bad_cond_str_ += "(量能死叉)10周下穿40周均线; "
	}
	if wsicha1 {
		score -= 5
		bad_cond_str += "(死叉)5周下穿10周均线; "
		bad_cond_str_ += "(死叉)5周下穿10周均线; "
	}
	if wsicha2 {
		score -= 5
		bad_cond_str += "(死叉)5周下穿30周均线; "
		bad_cond_str_ += "(死叉)5周下穿30周均线; "
	}
	if wsicha3 {
		score -= 5
		bad_cond_str += "(死叉)5周下穿60周均线; "
		bad_cond_str_ += "(死叉)5周下穿60周均线; "
	}
	if wsicha5 {
		score -= 5
		bad_cond_str += "(死叉)10周下穿30周均线; "
		bad_cond_str_ += "(死叉)10周下穿30周均线; "
	}
	if wsicha6 {
		score -= 5
		bad_cond_str += "(死叉)10周下穿60周均线; "
		bad_cond_str_ += "(死叉)10周下穿60周均线; "
	}
	if wsicha8 {
		score -= 5
		bad_cond_str += "(死叉)30周下穿60周均线; "
		bad_cond_str_ += "(死叉)30周下穿60周均线; "
	}

	if ysicha11 {
		score -= 5
		bad_cond_str += "(量能死叉)10月下穿40月均线; "
		bad_cond_str_ += "(量能死叉)10月下穿40月均线; "
	}
	if ysicha1 {
		score -= 5
		bad_cond_str += "(死叉)5月下穿10月均线; "
		bad_cond_str_ += "(死叉)5月下穿10月均线; "
	}
	if ysicha2 {
		score -= 5
		bad_cond_str += "(死叉)5月下穿30月均线; "
		bad_cond_str_ += "(死叉)5月下穿30月均线; "
	}
	if ysicha3 {
		score -= 5
		bad_cond_str += "(死叉)5月下穿60月均线; "
		bad_cond_str_ += "(死叉)5月下穿60月均线; "
	}
	if ysicha5 {
		score -= 5
		bad_cond_str += "(死叉)10月下穿30月均线; "
		bad_cond_str_ += "(死叉)10月下穿30月均线; "
	}
	if ysicha6 {
		score -= 5
		bad_cond_str += "(死叉)10月下穿60月均线; "
		bad_cond_str_ += "(死叉)10月下穿60月均线; "
	}
	if ysicha8 {
		score -= 5
		bad_cond_str += "(死叉)30月下穿60月均线; "
		bad_cond_str_ += "(死叉)30月下穿60月均线; "
	}

	if yiziban {
		score += 3
		cond_str += "一字板; "
		cond_str_ += "一字板; "
	}
	if tziban {
		score += 2
		cond_str += "T字板; "
		cond_str_ += "T字板; "
	}
	if zhangting {
		score += 1
		cond_str += "涨停股; "
		cond_str_ += "涨停股; "
	}
	if dietingban {
		score -= 3
		bad_cond_str += "一字跌停； "
		bad_cond_str_ += "一字跌停； "
	}
	if daotban {
		score -= 2
		bad_cond_str += "倒T跌停； "
		bad_cond_str_ += "倒T跌停； "
	}

	if dikaigaozou {
		score += 2
		cond_str += "低开高走; "
		cond_str_ += "低开高走; "
	}
	if gaokaigaozou {
		score += 2
		cond_str += "高开高走; "
		cond_str_ += "高开高走; "
	}
	if gaokaidizou {
		score -= 2
		bad_cond_str += "高开低走; "
		bad_cond_str_ += "高开低走; "
	}
	if dikaidizou {
		score -= 2
		bad_cond_str += "低开低走; "
		bad_cond_str_ += "低开低走; "
	}
	if wulianyang {
		score += 3
		cond_str += "五连阳; "
		cond_str_ += "五连阳; "
	}
	if silianyang {
		score += 2
		cond_str += "四连阳; "
		cond_str_ += "四连阳; "
	}
	if sanlianyang {
		score += 1
		cond_str += "三连阳; "
		cond_str_ += "三连阳; "
	}

	if wulianyin {
		score -= 2
		bad_cond_str += "五连阴; "
		bad_cond_str_ += "五连阴; "
	}

	if silianyin {
		score -= 2
		bad_cond_str += "四连阴; "
		bad_cond_str_ += "四连阴; "
	}

	if sanlianyin {
		score -= 2
		bad_cond_str += "三连阴; "
		bad_cond_str_ += "三连阴; "
	}

	if changshangying {
		score += 1
		cond_str += "长上影; "
		cond_str_ += "长上影; "
	}
	if changxiaying {
		score += 1
		cond_str += "长下影; "
		cond_str_ += "长下影; "
	}

	var tupo string
	if tupo1 {
		score += 3
		cond_str += "收盘突破5日压力位; "
		tupo += "5日,"
	}
	if tupo2 {
		score += 3
		cond_str += "收盘突破10日压力位; "
		tupo += "10日,"
	}
	if tupo3 {
		score += 3
		cond_str += "收盘突破30日压力位; "
		tupo += "30日,"
	}

	if tupo4 {
		score += 3
		cond_str += "收盘突破60日压力位; "
		tupo += "60日,"
	}

	if tupo5 {
		score += 3
		cond_str += "收盘突破5周压力位; "
		tupo += "5周,"
	}
	if tupo6 {
		score += 3
		cond_str += "收盘突破10周压力位; "
		tupo += "10周,"
	}
	if tupo7 {
		score += 3
		cond_str += "收盘突破30周压力位; "
		tupo += "30周,"
	}
	if tupo8 {
		score += 3
		cond_str += "收盘突破60周压力位; "
		tupo += "60周,"
	}
	if tupo9 {
		score += 3
		cond_str += "收盘突破5月压力位; "
		tupo += "5月,"
	}
	if tupo10 {
		score += 3
		cond_str += "收盘突破10月压力位; "
		tupo += "10月,"
	}
	if tupo11 {
		score += 3
		cond_str += "收盘突破30月压力位; "
		tupo += "30月,"
	}
	if tupo12 {
		score += 3
		cond_str += "收盘突破60月压力位; "
		tupo += "60月,"
	}
	if tupo != "" {
		tupo = strings.TrimRight(tupo, ",")
		cond_str_ += fmt.Sprintf("收盘价突破%s均线压力位; ", tupo)
	}

	var jichuang string
	if jichuang1 {
		score -= 3
		bad_cond_str += "收盘击穿5日支撑位; "
		jichuang += "5日,"
	}
	if jichuang2 {
		score -= 3
		bad_cond_str += "收盘击穿10日支撑位; "
		jichuang += "10日,"
	}
	if jichuang3 {
		score -= 3
		bad_cond_str += "收盘击穿30日支撑位; "
		jichuang += "30日,"
	}
	if jichuang4 {
		score -= 3
		bad_cond_str += "收盘击穿60日支撑位; "
		jichuang += "60日,"
	}

	if jichuang5 {
		score -= 3
		bad_cond_str += "收盘击穿5周支撑位; "
		jichuang += "5周,"
	}
	if jichuang6 {
		score -= 3
		bad_cond_str += "收盘击穿10周支撑位; "
		jichuang += "10周,"
	}
	if jichuang7 {
		score -= 3
		bad_cond_str += "收盘击穿30周支撑位; "
		jichuang += "30周,"
	}
	if jichuang8 {
		score -= 3
		bad_cond_str += "收盘击穿60周支撑位; "
		jichuang += "60周,"
	}
	if jichuang9 {
		score -= 3
		bad_cond_str += "收盘击穿5月支撑位; "
		jichuang += "5月,"
	}
	if jichuang10 {
		score -= 3
		bad_cond_str += "收盘击穿10月支撑位; "
		jichuang += "10月,"
	}
	if jichuang11 {
		score -= 2
		bad_cond_str += "收盘击穿30月支撑位; "
		jichuang += "30月,"
	}
	if jichuang12 {
		score -= 3
		bad_cond_str += "收盘击穿60月支撑位; "
		jichuang += "60月,"
	}

	if jichuang != "" {
		jichuang = strings.TrimRight(jichuang, ",")
		bad_cond_str_ += fmt.Sprintf("收盘价击穿%s均线支撑位; ", jichuang)
	}

	var shangshengtongdao string
	if priceshangyang1 {
		score += 2
		cond_str += "5日均线处于上升通道; "
		shangshengtongdao += "5日,"
	}
	if priceshangyang2 {
		score += 2
		cond_str += "10日均线处于上升通道; "
		shangshengtongdao += "10日,"
	}
	if priceshangyang3 {
		score += 2
		cond_str += "30日均线处于上升通道; "
		shangshengtongdao += "30日,"
	}
	if priceshangyang4 {
		score += 2
		cond_str += "60日均线处于上升通道; "
		shangshengtongdao += "60日,"
	}

	if wpriceshangyang1 {
		score += 2
		cond_str += "5周均线处于上升通道; "
		shangshengtongdao += "5周,"
	}
	if wpriceshangyang2 {
		score += 2
		cond_str += "10周均线处于上升通道; "
		shangshengtongdao += "10周,"
	}
	if wpriceshangyang3 {
		score += 2
		cond_str += "30周均线处于上升通道; "
		shangshengtongdao += "30周,"
	}
	if wpriceshangyang4 {
		score += 2
		cond_str += "60周均线处于上升通道; "
		shangshengtongdao += "60周,"
	}

	if ypriceshangyang1 {
		score += 2
		cond_str += "5月均线处于上升通道; "
		shangshengtongdao += "5月,"
	}
	if ypriceshangyang2 {
		score += 2
		cond_str += "10月均线处于上升通道; "
		shangshengtongdao += "10月,"
	}
	if ypriceshangyang3 {
		score += 2
		cond_str += "30月均线处于上升通道; "
		shangshengtongdao += "30月,"
	}
	if ypriceshangyang4 {
		score += 2
		cond_str += "60月均线处于上升通道; "
		shangshengtongdao += "60月,"
	}
	if shangshengtongdao != "" {
		shangshengtongdao = strings.TrimRight(shangshengtongdao, ",") + "均线处于上升通道; "
		cond_str_ += shangshengtongdao
	}

	var xiajiangtongdao string
	if pricexiajiang1 {
		score -= 2
		bad_cond_str += "5日均线处于下降通道; "
		xiajiangtongdao += "5日,"
	}
	if pricexiajiang2 {
		score -= 2
		bad_cond_str += "10日均线处于下降通道; "
		xiajiangtongdao += "10日,"
	}
	if pricexiajiang3 {
		score -= 2
		bad_cond_str += "30日均线处于下降通道; "
		xiajiangtongdao += "30日,"
	}
	if pricexiajiang4 {
		score -= 2
		bad_cond_str += "60日均线处于下降通道; "
		xiajiangtongdao += "60日,"
	}

	if wpricexiajiang1 {
		score -= 2
		bad_cond_str += "5周均线处于下降通道; "
		xiajiangtongdao += "5周,"
	}
	if wpricexiajiang2 {
		score -= 2
		bad_cond_str += "10周均线处于下降通道; "
		xiajiangtongdao += "10周,"
	}
	if wpricexiajiang3 {
		score -= 2
		bad_cond_str += "30周均线处于下降通道; "
		xiajiangtongdao += "30周,"
	}
	if wpricexiajiang4 {
		score -= 2
		bad_cond_str += "60周均线处于下降通道; "
		xiajiangtongdao += "60周,"
	}

	if ypricexiajiang1 {
		score -= 2
		bad_cond_str += "5月均线处于下降通道; "
		xiajiangtongdao += "5月,"
	}
	if ypricexiajiang2 {
		score -= 2
		bad_cond_str += "10月均线处于下降通道; "
		xiajiangtongdao += "10月,"
	}
	if ypricexiajiang3 {
		score -= 2
		bad_cond_str += "30月均线处于下降通道; "
		xiajiangtongdao += "30月,"
	}
	if ypricexiajiang4 {
		score -= 2
		bad_cond_str += "60月均线处于下降通道; "
		xiajiangtongdao += "60月,"
	}
	if xiajiangtongdao != "" {
		xiajiangtongdao = strings.TrimRight(xiajiangtongdao, ",") + "均线处于下降通道; "
		bad_cond_str_ += xiajiangtongdao
	}

	var shangfang string
	if priceaboveave1 {
		score += 2
		cond_str += "当前价位在5日均线上方; "
		shangfang += "5日,"
	}
	if priceaboveave2 {
		score += 2
		cond_str += "当前价位在10日均线上方; "
		shangfang += "10日,"
	}
	if priceaboveave3 {
		score += 2
		cond_str += "当前价位在30日均线上方; "
		shangfang += "30日,"
	}
	if priceaboveave4 {
		score += 2
		cond_str += "当前价位在60日均线上方; "
		shangfang += "60日,"
	}

	if wpriceaboveave1 {
		score += 2
		cond_str += "当前价位在5周均线上方; "
		shangfang += "5周,"
	}
	if wpriceaboveave2 {
		score += 2
		cond_str += "当前价位在10周均线上方; "
		shangfang += "10周,"
	}
	if wpriceaboveave3 {
		score += 2
		cond_str += "当前价位在30周均线上方; "
		shangfang += "30周,"
	}
	if wpriceaboveave4 {
		score += 2
		cond_str += "当前价位在60周均线上方; "
		shangfang += "60周,"
	}

	if ypriceaboveave1 {
		score += 2
		cond_str += "当前价位在5月均线上方; "
		shangfang += "5月,"
	}
	if ypriceaboveave2 {
		score += 2
		cond_str += "当前价位在10月均线上方; "
		shangfang += "10月,"
	}
	if ypriceaboveave3 {
		score += 2
		cond_str += "当前价位在30月均线上方; "
		shangfang += "30月,"
	}
	if ypriceaboveave4 {
		score += 2
		cond_str += "当前价位在60月均线上方; "
		shangfang += "60月,"
	}
	if shangfang != "" {
		shangfang = strings.TrimRight(shangfang, ",")
		cond_str_ += fmt.Sprintf("收盘价位在%s均线上方; ", shangfang)
	}

	var xiafang string
	if !priceaboveave1 {
		score -= 2
		bad_cond_str += "当前价位在5日均线下方; "
		xiafang += "5日,"
	}
	if !priceaboveave2 {
		score -= 2
		bad_cond_str += "当前价位在10日均线下方; "
		xiafang += "10日,"
	}
	if !priceaboveave3 {
		score -= 2
		bad_cond_str += "当前价位在30日均线下方; "
		xiafang += "30日,"
	}
	if !priceaboveave4 {
		score -= 2
		bad_cond_str += "当前价位在60日均线下方; "
		xiafang += "60日,"
	}

	if !wpriceaboveave1 {
		score -= 2
		bad_cond_str += "当前价位在5周均线下方; "
		xiafang += "5周,"
	}
	if !wpriceaboveave2 {
		score -= 2
		bad_cond_str += "当前价位在10周均线下方; "
		xiafang += "10周,"
	}
	if !wpriceaboveave3 && result.AveWeeklyPrice3 != nil {
		score -= 2
		bad_cond_str += "当前价位在30周均线下方; "
		xiafang += "30周,"
	}
	if !wpriceaboveave4 && result.AveWeeklyPrice4 != nil {
		score -= 2
		bad_cond_str += "当前价位在60周均线下方; "
		xiafang += "60周,"
	}

	if !ypriceaboveave1 && result.AveMonthPrice1 != nil {
		score -= 2
		bad_cond_str += "当前价位在5月均线下方; "
		xiafang += "5月,"
	}
	if !ypriceaboveave2 && result.AveMonthPrice2 != nil {
		score -= 2
		bad_cond_str += "当前价位在10月均线下方; "
		xiafang += "10月,"
	}
	if !ypriceaboveave3 && result.AveMonthPrice3 != nil {
		score -= 2
		bad_cond_str += "当前价位在30月均线下方; "
		xiafang += "30月,"
	}
	if !ypriceaboveave4 && result.AveMonthPrice4 != nil {
		score -= 2
		bad_cond_str += "当前价位在60月均线下方; "
		xiafang += "60月,"
	}

	if xiafang != "" {
		xiafang = strings.TrimRight(xiafang, ",")
		bad_cond_str_ += fmt.Sprintf("收盘价位在%s均线下方; ", xiafang)
	}

	//if junjialianhe1 {
	//	cond_str += "近期5日均线与收盘价黏合; "
	//}
	//if junjialianhe2 {
	//	cond_str += "近期10日均线与收盘价黏合; "
	//}
	//if junjialianhe3 {
	//	cond_str += "近期30日均线与收盘价黏合; "
	//}
	//if junjialianhe4 {
	//	cond_str += "近期60日均线与收盘价黏合; "
	//}

	// 量价
	if liangshangyang1 {
		score += 3
		cond_str += "量能10日均线处于上升通道; "
		cond_str_ += "量能10日均线处于上升通道; "
	}
	if liangshangyang2 {
		score += 3
		cond_str += "量能40日均线处于上升通道; "
		cond_str_ += "量能40日均线处于上升通道; "
	}
	if liangnengbigger1 || liangnengbigger2 {
		score += 3
		cond_str += "近期量能相对活跃; "
		cond_str_ += "近期量能相对活跃; "
	}
	if liangnengsmaller1 || liangnengsmaller2 {
		score -= 3
		bad_cond_str += "近期量能相对萎靡; "
		bad_cond_str_ += "近期量能相对萎靡; "
	}
	if liangnengbuduanbigger {
		score += 2
		cond_str += "量能不断放大; "
		cond_str_ += "量能不断放大; "
	}
	if tufangjuliang {
		score += 2
		cond_str += "突放巨量; "
		cond_str_ += "突放巨量; "
	}
	if liangnengtupo1 {
		score += 2
		cond_str += "量能突破10日均线; "
		cond_str_ += "量能突破10日均线; "

	}
	if liangnengtupo2 {
		score += 2
		cond_str += "量能突破40日均线; "
		cond_str_ += "量能突破40日均线; "
	}

	// 周量价
	if wliangshangyang1 {
		score += 3
		cond_str += "量能10周均线处于上升通道; "
		cond_str_ += "量能10周均线处于上升通道; "
	}
	if wliangshangyang2 {
		score += 3
		cond_str += "量能40周均线处于上升通道; "
		cond_str_ += "量能40周均线处于上升通道; "
	}
	if wliangnengbuduanbigger {
		score += 2
		cond_str += "量能周线不断放大; "
		cond_str_ += "量能周线不断放大; "
	}
	if wtufangjuliang {
		score += 2
		cond_str += "量能周线突放巨量; "
		cond_str_ += "量能周线突放巨量; "
	}
	if wliangnengtupo1 {
		score += 2
		cond_str += "量能突破10周均线; "
		cond_str_ += "量能突破10周均线; "
	}
	if wliangnengtupo2 {
		score += 2
		cond_str += "量能突破40周均线; "
		cond_str_ += "量能突破40周均线; "
	}

	// 月量价
	if yliangshangyang1 {
		score += 3
		cond_str += "量能10月均线处于上升通道; "
		cond_str_ += "量能10月均线处于上升通道; "
	}
	if yliangshangyang2 {
		score += 3
		cond_str += "量能40月均线处于上升通道; "
		cond_str_ += "量能40月均线处于上升通道; "
	}
	if yliangnengbuduanbigger {
		score += 2
		cond_str += "量能月线不断放大; "
		cond_str_ += "量能月线不断放大; "
	}
	if ytufangjuliang {
		score += 2
		cond_str += "量能月线突放巨量; "
		cond_str_ += "量能月线突放巨量; "
	}
	if yliangnengtupo1 {
		score += 2
		cond_str += "量能突破10月均线; "
		cond_str_ += "量能突破10月均线; "
	}
	if yliangnengtupo2 {
		score += 2
		cond_str += "量能突破40月均线; "
		cond_str_ += "量能突破40月均线; "
	}

	if guoyi {
		cond_str += "当日成交总额上亿; "
		cond_str_ += "当日成交总额上亿; "
	}

	if simuchicangcount > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d个私募持仓; ", simuchicangcount)
		cond_str_ += fmt.Sprintf("%d个私募持仓; ", simuchicangcount)
	}
	if jigouchicangcount > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d个基金持仓; ", jigouchicangcount)
		cond_str_ += fmt.Sprintf("%d个基金持仓; ", jigouchicangcount)
	}
	if fenhong > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d次分红; ", fenhong)
		cond_str_ += fmt.Sprintf("%d次分红; ", fenhong)
	}
	if songgu > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d次送股; ", songgu)
		cond_str_ += fmt.Sprintf("%d次送股; ", songgu)
	}
	if zhuangzeng > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d次转增; ", zhuangzeng)
		cond_str_ += fmt.Sprintf("%d次转增; ", zhuangzeng)
	}
	if pergu > 0 {
		score -= 1
		bad_cond_str += fmt.Sprintf("%d次配股; ", pergu)
		bad_cond_str_ += fmt.Sprintf("%d次配股; ", pergu)
	}
	if zengfa > 0 {
		score -= 1
		bad_cond_str += fmt.Sprintf("%d次增发; ", zengfa)
		bad_cond_str_ += fmt.Sprintf("%d次增发; ", zengfa)
	}
	if subcomp > 0 {
		cond_str += fmt.Sprintf("%d个参股公司; ", subcomp)
		cond_str_ += fmt.Sprintf("%d个参股公司; ", subcomp)
	}
	if changename > 0 {
		bad_cond_str += fmt.Sprintf("历史更名%d次; ", changename)
		bad_cond_str_ += fmt.Sprintf("历史更名%d次; ", changename)
		if has_st {
			score -= 2
			bad_cond_str += "曾ST带帽; "
			bad_cond_str_ += "曾ST带帽; "
		}
	}
	// 基本面 现金流量表
	if up1 {
		score += 1
		cond_str += "经营现金流量净额非负; "
		cond_str_ += "经营现金流量净额非负; "
	}
	if up2 {
		score += 1
		cond_str += "投资现金流量净额非负; "
		cond_str_ += "投资现金流量净额非负; "
	}
	if up3 {
		score += 1
		cond_str += "筹资现金流量净额非负; "
		cond_str_ += "筹资现金流量净额非负; "
	}
	if up4 {
		score += 1
		cond_str += "期末现金及现金等价物余额非负; "
		cond_str_ += "期末现金及现金等价物余额非负; "
	}
	// 基本面 利润表
	if pup1 {
		score += 1
		cond_str += "营业总收入非负; "
		cond_str_ += "营业总收入非负; "
	}
	if pup2 {
		score += 1
		cond_str += "净利润非负; "
		cond_str_ += "净利润非负; "
	}
	// 基本面 资产负债表
	if lup1 {
		score += 1
		cond_str += "总资产不断增加; "
		cond_str_ += "总资产不断增加; "
	}
	if done1 {
		score += 1
		cond_str += "总负债不断减小; "
		cond_str_ += "总负债不断减小; "
	}

	if !up1 {
		score -= 1
		bad_cond_str += "(亏损可能)经营现金流量净额出现负值; "
		bad_cond_str_ += "(亏损可能)经营现金流量净额出现负值; "
	}
	if !up2 {
		score -= 1
		bad_cond_str += "(亏损可能)投资现金流量净额出现负值; "
		bad_cond_str_ += "(亏损可能)投资现金流量净额出现负值; "
	}
	if !up3 {
		score -= 1
		bad_cond_str += "(亏损可能)筹资现金流量净额出现负值; "
		bad_cond_str_ += "(亏损可能)筹资现金流量净额出现负值; "
	}
	if !up4 {
		score -= 1
		bad_cond_str += "(亏损可能)期末现金及现金等价物余额出现负值; "
		bad_cond_str_ += "(亏损可能)期末现金及现金等价物余额出现负值; "

	}
	if !pup1 {
		score -= 2
		bad_cond_str += "营业总收入亏损; "
		bad_cond_str_ += "营业总收入亏损; "

	}
	if !pup2 {
		score -= 2
		bad_cond_str += "净利润亏损; "
		bad_cond_str_ += "净利润亏损; "

	}
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
	score += high*4 + middle*2 + low*-2 + bad*-4

	if score < 0 {
		score = 1
	}

	if score >= 100 {
		score = 95
	}

	fmt.Println(code, name, cond_str, bad_cond_str, finance, cond_str_, bad_cond_str_)
	seed := []int{1, 2, 3, 4, 5}
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(seed)

	p := dal.Predict{Code: code, Name: name, Condition: cond_str_, BadCondition: bad_cond_str_, Condition_: cond_str, BadCondition_: bad_cond_str, Finance: finance,
		Date: result.CurrDate, Score: score + seed[n], Price: result.RecentClose[0], Percent: result.RecentPercent[0],
		FundCount: jigouchicangcount, SMCount: simuchicangcount, FenghongCount: fenhong, PeiguCount: pergu, ZhuangzengCount: zhuangzeng,
		SongguCount: songgu, ZengfaCount: zengfa, SubcompCount: subcomp}
	if utils.TellEnv() == "loc" {
		err := store.MysqlClient.GetOnlineDB().Save(&p).Error
		if err != nil {
			fmt.Println("写入线上错误")
		}
	}
	store.MysqlClient.GetDB().Save(&p)
}
