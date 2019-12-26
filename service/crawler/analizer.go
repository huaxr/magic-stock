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
	store.MysqlClient.GetDB().Model(&dal.FundHoldRank{}).Where("code = ?", code).Count(&count)
	return count
}

func STStock(code string) bool {
	var c int
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", code).Where("`concept` regexp ?", "ST").Count(&c)
	return c > 0
}

func GetStockPercent(result *model.CalcResult, score int) (zhangdie, huanshoulv, zhenfures string, scores int) {
	percent := result.RecentPercent[0]
	huanshou := result.RecentTurnoverRate[0]
	zhenfu := result.RecentAmplitude[0]
	if percent >= 8 {
		zhangdie = "涨幅高于8%; "
	}

	if percent >= 5 && percent <= 8 {
		zhangdie = "涨幅介于5% ~ 8%; "
	}

	if percent >= 2 && percent <= 5 {
		zhangdie = "涨幅介于2% ~ 5%; "
	}

	if percent >= -2 && percent <= 2 {
		zhangdie = "涨跌幅介于-2% ~ 2%; "
	}

	if percent <= -8 {
		zhangdie = "跌幅高于8%; "
	}

	if percent >= -8 && percent <= -5 {
		zhangdie = "跌幅介于-8% ~ -5%; "
	}

	if percent >= -5 && percent <= -2 {
		zhangdie = "跌幅介于-5% ~ -2%; "
	}

	if huanshou <= 5 {
		huanshoulv = "换手率小于5%; "
	}

	if huanshou >= 5 && huanshou <= 10 {
		huanshoulv = "换手率介于5% ~ 10%; "
	}

	if huanshou >= 10 {
		score += 1
		huanshoulv = "换手率大于10%; "
	}

	if zhenfu <= 5 {
		zhenfures = "振幅小于5%; "
	}

	if zhenfu >= 5 && zhenfu <= 10 {
		zhenfures = "振幅介于5% ~ 10%; "
	}

	if zhenfu >= 10 && zhenfu <= 15 {
		score += 1
		zhenfures = "振幅介于10% ~ 15%; "
	}

	if zhenfu >= 15 {
		score += 1
		zhenfures = "振幅大于15%; "
	}
	return zhangdie, huanshoulv, zhenfures, score
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

// recent_num type 涨幅滞涨 limit_percent
func (craw *Crawler) WeeklyPercentLimited(result *model.CalcResult, recent_num int, limit_percent float64, typ string) bool {
	switch typ {
	case "week":
		return ArraySumLessThan(result.RecentWeeklyPercent[0:recent_num], limit_percent)
	case "day":
		return ArraySumLessThan(result.RecentPercent[0:recent_num], limit_percent)
	}
	return false
}

func (craw *Crawler) HasLimitUpInTheseDays(result *model.CalcResult, recent_days int) bool {
	return ArrayHasLimitUp(result.RecentPercent, recent_days)
}

func xx(x interface{}) {

}

func (craw *Crawler) Analyze(result *model.CalcResult, code, name string) {
	// 今日线金叉 6 和 15线
	jincha1 := result.AveDailyPrice1[0] > result.AveDailyPrice2[0] && result.AveDailyPrice1[1] < result.AveDailyPrice2[1]
	sicha1 := result.AveDailyPrice1[0] < result.AveDailyPrice2[0] && result.AveDailyPrice1[1] > result.AveDailyPrice2[1]
	// 昨日线金叉 6 和 15线
	// jincha2 := result.AveDailyPrice1[1] > result.AveDailyPrice2[1] && result.AveDailyPrice1[2] < result.AveDailyPrice2[2]
	// 今日线金叉 15 和 30线
	jincha3 := result.AveDailyPrice2[0] > result.AveDailyPrice3[0] && result.AveDailyPrice2[1] < result.AveDailyPrice3[1]
	sicha3 := result.AveDailyPrice2[0] < result.AveDailyPrice3[0] && result.AveDailyPrice2[1] > result.AveDailyPrice3[1]
	// 昨日线金叉 15 和 30线
	//jincha4 := result.AveDailyPrice2[1] > result.AveDailyPrice3[1] && result.AveDailyPrice2[2] < result.AveDailyPrice3[2]
	// 本周线金叉
	jincha5 := result.AveWeeklyPrice1[0] > result.AveWeeklyPrice2[0] && result.AveWeeklyPrice1[1] < result.AveWeeklyPrice2[1]
	sicha5 := result.AveWeeklyPrice1[0] < result.AveWeeklyPrice2[0] && result.AveWeeklyPrice1[1] > result.AveWeeklyPrice2[1]
	// 上周金叉
	//jincha6 := result.AveWeeklyPrice1[1] > result.AveWeeklyPrice2[1] && result.AveWeeklyPrice1[2] < result.AveWeeklyPrice2[2]
	// 量能今日金叉
	jincha7 := result.AveCount1[0] > result.AveCount2[0] && result.AveCount1[1] < result.AveCount2[1]
	// 量能昨日金叉
	jincha8 := result.AveCount1[1] > result.AveCount2[1] && result.AveCount1[2] < result.AveCount2[2]

	// 6日均线 价格均线上扬
	priceshangyang1 := result.AveDailyPrice1[0] > result.AveDailyPrice1[1] && result.AveDailyPrice1[1] > result.AveDailyPrice1[2] && result.AveDailyPrice1[2] > result.AveDailyPrice1[3]
	//15日均线 价格均线上扬
	priceshangyang2 := result.AveDailyPrice2[0] > result.AveDailyPrice2[1] && result.AveDailyPrice2[1] > result.AveDailyPrice2[2] && result.AveDailyPrice2[2] > result.AveDailyPrice2[3]
	// 30日均线 价格均线上扬
	priceshangyang3 := result.AveDailyPrice3[0] > result.AveDailyPrice3[1] && result.AveDailyPrice3[1] > result.AveDailyPrice3[2] && result.AveDailyPrice3[2] > result.AveDailyPrice3[3]

	// 6周均线线 价格上扬
	priceshangyang4 := result.AveWeeklyPrice1[0] > result.AveWeeklyPrice1[1] && result.AveWeeklyPrice1[1] > result.AveWeeklyPrice1[2] && result.AveWeeklyPrice1[2] > result.AveWeeklyPrice1[3]
	// 15周均线线 价格上扬
	priceshangyang5 := result.AveWeeklyPrice2[0] > result.AveWeeklyPrice2[1] && result.AveWeeklyPrice2[1] > result.AveWeeklyPrice2[2] && result.AveWeeklyPrice2[2] > result.AveWeeklyPrice2[3]

	pricexiajiang1 := result.AveDailyPrice1[0] < result.AveDailyPrice1[1] && result.AveDailyPrice1[1] < result.AveDailyPrice1[2] && result.AveDailyPrice1[2] < result.AveDailyPrice1[3]
	//15日均线 价格均线下挫
	pricexiajiang2 := result.AveDailyPrice2[0] < result.AveDailyPrice2[1] && result.AveDailyPrice2[1] < result.AveDailyPrice2[2] && result.AveDailyPrice2[2] < result.AveDailyPrice2[3]
	// 30日均线 价格均线下挫
	pricexiajiang3 := result.AveDailyPrice3[0] < result.AveDailyPrice3[1] && result.AveDailyPrice3[1] < result.AveDailyPrice3[2] && result.AveDailyPrice3[2] < result.AveDailyPrice3[3]

	// 6周均线线 价格下挫
	pricexiajiang4 := result.AveWeeklyPrice1[0] < result.AveWeeklyPrice1[1] && result.AveWeeklyPrice1[1] < result.AveWeeklyPrice1[2] && result.AveWeeklyPrice1[2] < result.AveWeeklyPrice1[3]
	// 15周均线线 价格下挫
	pricexiajiang5 := result.AveWeeklyPrice2[0] < result.AveWeeklyPrice2[1] && result.AveWeeklyPrice2[1] < result.AveWeeklyPrice2[2] && result.AveWeeklyPrice2[2] < result.AveWeeklyPrice2[3]

	// 高位回调
	gaoweihuitiao1 := result.RecentPercent[0] < 3 && result.RecentPercent[1] < 9 && result.RecentPercent[2] > 9 && (result.RecentPercent[0]+result.RecentPercent[1]+result.RecentPercent[2] < 3)
	gaoweihuitiao2 := result.RecentPercent[0] < 3 && result.RecentPercent[1] < 3 && result.RecentPercent[2] < 9 && result.RecentPercent[3] > 9 && (result.RecentPercent[0]+result.RecentPercent[1]+result.RecentPercent[2]+result.RecentPercent[3] < 3)
	gaoweihuitiao3 := result.RecentPercent[0] < 3 && result.RecentPercent[1] < 3 && result.RecentPercent[2] < 3 && result.RecentPercent[3] < 9 && result.RecentPercent[4] > 9 && (result.RecentPercent[0]+result.RecentPercent[1]+result.RecentPercent[2]+result.RecentPercent[3]+result.RecentPercent[4] < 3)

	// 连续5日量能10均线上扬
	liangshangyang1 := result.AveCount1[0] > result.AveCount1[1] && result.AveCount1[1] > result.AveCount1[2] && result.AveCount1[2] > result.AveCount1[3] && result.AveCount1[3] > result.AveCount1[4] && result.AveCount1[4] > result.AveCount1[5]
	// 连续5日量能40均线上扬
	liangshangyang2 := result.AveCount2[0] > result.AveCount2[1] && result.AveCount2[1] > result.AveCount2[2] && result.AveCount2[2] > result.AveCount2[3] && result.AveCount2[3] > result.AveCount2[4] && result.AveCount2[4] > result.AveCount2[5]

	// 连续5日量能站上10均线
	liangnengbigger1 := result.AveCount1[0] < result.RecentCount[0] && result.AveCount1[1] < result.RecentCount[1] && result.AveCount1[2] < result.RecentCount[2] && result.AveCount1[3] < result.RecentCount[3] && result.AveCount1[4] < result.RecentCount[4]
	// 连续5日量能站上40均线
	liangnengbigger2 := result.AveCount2[0] < result.RecentCount[0] && result.AveCount2[1] < result.RecentCount[1] && result.AveCount2[2] < result.RecentCount[2] && result.AveCount2[3] < result.RecentCount[3] && result.AveCount2[4] < result.RecentCount[4]
	// 量能不断放大
	liangnengbuduanbigger := result.RecentCount[0] > result.RecentCount[1] && result.RecentCount[1] > result.RecentCount[2] && result.RecentCount[2] > result.RecentCount[3]
	// 突放巨量
	tufangjuliang := (result.RecentCount[0]-result.RecentCount[1])/result.RecentCount[1] > 5 || (result.RecentCount[1]-result.RecentCount[2])/result.RecentCount[1] > 5

	// 量能小于10日均线
	liangnengsmall1 := result.AveCount1[0] > result.RecentCount[0]
	// 量能小于40日均线
	liangnengsmall2 := result.AveCount2[0] > result.RecentCount[0]

	// 3 连阳
	sanlianyang := result.RecentPercent[0] > 0 && result.RecentPercent[1] > 0 && result.RecentPercent[2] > 0 && result.RecentPercent[3] > 0
	// 5连阳
	wulianyang := sanlianyang && result.RecentPercent[4] > 0 && result.RecentPercent[5] > 0

	// 近期长上影
	changshangying := craw.HasTopLine(result, 1)
	// 近期长下影
	changxiaying := craw.HasLowLine(result, 1)

	// 优良概念
	//goodconcept := GetConceptByCode(code, "预盈预增|业绩预升|高派息|独角兽|高送转|基金重仓|QFII|RQFII")
	// 私募持仓
	simuchicangcount := GetHolderByCode(code, "私募")
	// 基金持仓
	jigouchicangcount := GetFundByCode(code)

	// 均价粘合
	junjialianhe1 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice1, 5, 4)
	junjialianhe2 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice2, 5, 4)
	junjialianhe3 := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice3, 5, 4)

	// 涨停股
	zhangting := result.RecentPercent[0] > 9.94
	// 一字板
	yiziban := result.RecentPercent[0] > 9.94 && result.RecentClose[0] == result.RecentOpen[0]

	// 当前价格在短期均线上方
	priceaboveave6 := result.RecentClose[0] >= result.AveDailyPrice1[0]
	priceaboveave15 := result.RecentClose[0] >= result.AveDailyPrice2[0]
	priceaboveave30 := result.RecentClose[0] >= result.AveDailyPrice3[0]
	pricelowave6 := result.RecentClose[0] < result.AveDailyPrice1[0]
	pricelowave15 := result.RecentClose[0] < result.AveDailyPrice2[0]
	pricelowave30 := result.RecentClose[0] < result.AveDailyPrice3[0]

	// 成交过亿
	guoyi := result.CurrTotalMoney > 10000

	// 现金流量表 非负
	up1, up2, up3, up4 := GetUpManageCashFlow(code)
	// 利润表 不能为负
	pup1, pup2 := GetProfitNotMiner(code)
	// 资产负债表
	lup1, done1 := GetUpLiabilities(code)

	// st
	st := STStock(code)

	score := 0 // max 37  // low 17
	cond_str, bad_cond_str, finance := "", "", ""
	if priceaboveave6 {
		score += 1
		cond_str += fmt.Sprintf("收盘价在6日均线上方; ")
	}
	if priceaboveave15 {
		score += 1
		cond_str += fmt.Sprintf("收盘价在15日均线上方; ")
	}
	if priceaboveave30 {
		score += 1
		cond_str += fmt.Sprintf("收盘价在30日均线上方; ")
	}
	if pricelowave6 {
		score -= 1
		bad_cond_str += fmt.Sprintf("收盘价在6日均线下方; ")
	}
	if pricelowave15 {
		score -= 1
		bad_cond_str += fmt.Sprintf("收盘价在15日均线下方; ")
	}
	if pricelowave30 {
		score -= 1
		bad_cond_str += fmt.Sprintf("收盘价在30日均线下方; ")
	}
	if jincha1 {
		score += 2
		cond_str += "6日均线与15日均线交金叉; "
	}
	if sicha1 {
		score -= 2
		bad_cond_str += "6日均线与15日均线交死叉; "
	}
	if jincha3 {
		score += 2
		cond_str += "15日均线与30日均线交金叉; "
	}
	if sicha3 {
		score -= 2
		bad_cond_str += "15日均线与30日均线交死叉; "
	}
	if jincha5 {
		score += 2
		cond_str += "6周均线与15周均线交金叉; "
	}
	if sicha5 {
		score -= 2
		bad_cond_str += "6周均线与15周均线交死叉; "
	}
	if priceshangyang1 {
		score += 1
		cond_str += "6日均线上扬; "
	}
	if priceshangyang2 {
		score += 1
		cond_str += "15日均线上扬; "
	}
	if priceshangyang3 {
		score += 1
		cond_str += "30日均线上扬; "
	}

	if priceshangyang4 {
		score += 1
		cond_str += "6周均线上扬; "
	}
	if priceshangyang5 {
		score += 1
		cond_str += "15周均线上扬; "
	}

	if pricexiajiang1 {
		score -= 2
		bad_cond_str += "6日均线下挫; "
	}
	if pricexiajiang2 {
		score -= 1
		bad_cond_str += "15日均线下挫; "
	}
	if pricexiajiang3 {
		score -= 1
		bad_cond_str += "30日均线下挫; "
	}

	if pricexiajiang4 {
		score -= 1
		bad_cond_str += "6周均线下挫; "
	}

	if pricexiajiang5 {
		score -= 1
		bad_cond_str += "15周均线下挫; "
	}

	if changshangying {
		cond_str += "长上影; "
	}
	if changxiaying {
		cond_str += "长下影; "
	}
	if zhangting {
		score += 1
		cond_str += "涨停股; "
	}
	if yiziban {
		score += 1
		cond_str += "一字板; "
	}
	if gaoweihuitiao1 || gaoweihuitiao2 || gaoweihuitiao3 {
		cond_str += "高位回调; "
	}
	if junjialianhe1 {
		cond_str += "近5天6日均线与收盘价粘合; "
	}
	if junjialianhe2 {
		cond_str += "近5天15日均线与收盘价粘合; "
	}
	if junjialianhe3 {
		cond_str += "近5天30日均线与收盘价粘合; "
	}
	if sanlianyang {
		score += 1
		cond_str += "三连阳; "
	}
	if wulianyang {
		score += 1
		cond_str += "五连阳; "
	}
	// 量价
	if jincha7 || jincha8 {
		score += 1
		cond_str += "10日与40日量能均线交金叉; "
	}
	if liangshangyang1 {
		score += 1
		cond_str += "连续5日量能10日均线上扬; "
	}
	if liangshangyang2 {
		score += 1
		cond_str += "连续5日量能40日均线上扬; "
	}
	if liangnengbigger1 {
		score += 1
		cond_str += "连续5日量能站上10日均线; "
	}
	if liangnengbigger2 {
		score += 1
		cond_str += "连续5日量能站上40日均线; "
	}
	if liangnengsmall1 {
		score -= 2
		bad_cond_str += "成交量(不活跃)低于10日均线; "
	}
	if liangnengsmall2 {
		score -= 2
		bad_cond_str += "成交量(不活跃)低于40日均线; "
	}
	if liangnengbuduanbigger {
		score += 1
		cond_str += "量能不断放大; "
	}
	if tufangjuliang {
		score += 1
		cond_str += "突放巨量; "
	}
	if guoyi {
		cond_str += "成交额过亿; "
	}
	// 基本面 现金流量表
	if up1 {
		score += 1
		cond_str += "经营现金流量净额非负; "
	}
	if !up1 {
		score -= 1
		bad_cond_str += "经营现金流量净额出现负值(亏损可能); "
	}
	if up2 {
		score += 1
		cond_str += "投资现金流量净额非负; "
	}
	if !up2 {
		score -= 1
		bad_cond_str += "投资现金流量净额出现负值(亏损可能); "
	}
	if up3 {
		score += 1
		cond_str += "筹资现金流量净额非负; "
	}
	if !up3 {
		score -= 1
		bad_cond_str += "筹资现金流量净额出现负值(亏损可能); "
	}
	if up4 {
		score += 1
		cond_str += "期末现金及现金等价物余额非负; "
	}
	if !up4 {
		score -= 1
		bad_cond_str += "期末现金及现金等价物余额出现负值(亏损可能); "
	}
	// 基本面 利润表
	if pup1 {
		score += 1
		cond_str += "营业总收入非负; "
	}
	if !pup1 {
		score -= 2
		bad_cond_str += "营业总收入亏损; "
	}
	if pup2 {
		score += 1
		cond_str += "净利润非负; "
	}
	if !pup2 {
		score -= 2
		bad_cond_str += "净利润亏损; "
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
	// 机构持仓情况
	// 十大流通股东信息
	if simuchicangcount > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d个私募持仓; ", simuchicangcount)
	}
	if jigouchicangcount > 0 {
		score += 1
		cond_str += fmt.Sprintf("%d个基金持仓; ", jigouchicangcount)
	}

	if st {
		score -= 3
		bad_cond_str += "ST垃圾股; "
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
	score += high*2 + middle*1 + low*-1 + bad*-2
	fmt.Println(code, name, cond_str, bad_cond_str, finance)
	seed := []int{1, 2, 3, 4, 5, 6}
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(seed)

	p := dal.Predict{Code: code, Name: name, Condition: cond_str, BadCondition: bad_cond_str, Finance: finance, Date: result.CurrDate, FundCount: jigouchicangcount, SMCount: simuchicangcount, Score: score*2 + seed[n], Price: result.RecentClose[0], Percent: result.RecentPercent[0]}
	if utils.TellEnv() == "loc" {
		err := store.MysqlClient.GetOnlineDB().Save(&p).Error
		if err != nil {
			fmt.Println("写入线上错误")
		}
	}
	store.MysqlClient.GetDB().Save(&p)

}
