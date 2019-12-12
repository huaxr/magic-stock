// @Time:       2019/12/3 下午2:33

package crawler

import (
	"fmt"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"math"
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
// 1 main 为 1 个点的时候， 虚线为3个点以上
// 2 main 为 2 个点的时候, 虚线为6个点以上
// 3 main 为 3个点  虚线为 7 点以上
// 4 main 为4个点 虚线5个点以上

// 上影线是实体柱子4倍以上
func (craw *Crawler) HasTopLine(result *model.CalcResult, recent int) bool {
	closes := result.RecentClose[0:recent]
	for index, i := range closes {
		// 开盘价大于收盘价
		if result.RecentOpen[index]-i >= 0 {
			main := result.RecentOpen[index] - result.RecentClose[index]
			virtual := result.RecentHigh[index] - result.RecentOpen[index]
			if main <= result.RecentClose[index+1]*0.1*Main || virtual <= result.RecentClose[index+1]*0.1*Virtual {
				continue
			}
			if virtual/main > Num {
				return true
			}
		} else {
			main := result.RecentClose[index] - result.RecentOpen[index]
			virtual := result.RecentHigh[index] - result.RecentClose[index]
			if main <= result.RecentClose[index+1]*0.1*Main || virtual <= result.RecentClose[index+1]*0.1*Virtual {
				continue
			}
			if virtual/main > Num {
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
			main := result.RecentOpen[index] - i
			virtual := result.RecentClose[index] - result.RecentLow[index]
			if main <= result.RecentClose[index+1]*0.1*Main || virtual <= result.RecentClose[index+1]*0.1*Virtual {
				continue
			}
			if virtual/main > Num {
				return true
			}
		} else {
			main := result.RecentClose[index] - result.RecentOpen[index]
			virtual := result.RecentOpen[index] - result.RecentLow[index]
			if main <= result.RecentClose[index+1]*0.1*Main || virtual <= result.RecentClose[index+1]*0.1*Virtual {
				continue
			}
			if virtual/main > Num {
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
	// 昨日线金叉 6 和 15线
	jincha2 := result.AveDailyPrice1[1] > result.AveDailyPrice2[1] && result.AveDailyPrice1[2] < result.AveDailyPrice2[2]
	// 今日线金叉 15 和 30线
	jincha3 := result.AveDailyPrice2[0] > result.AveDailyPrice3[0] && result.AveDailyPrice2[1] < result.AveDailyPrice3[1]
	// 昨日线金叉 15 和 30线
	jincha4 := result.AveDailyPrice2[1] > result.AveDailyPrice3[1] && result.AveDailyPrice2[2] < result.AveDailyPrice3[2]
	// 本周线金叉
	jincha5 := result.AveWeeklyPrice1[0] > result.AveWeeklyPrice2[0] && result.AveWeeklyPrice1[1] < result.AveWeeklyPrice2[1]
	// 上周金叉
	jincha6 := result.AveWeeklyPrice1[1] > result.AveWeeklyPrice2[1] && result.AveWeeklyPrice1[2] < result.AveWeeklyPrice2[2]
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
	priceshangyang5 := result.AveWeeklyPrice2[0] > result.AveWeeklyPrice2[1] && result.AveWeeklyPrice2[1] > result.AveWeeklyPrice2[2]

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
	// 5连阳
	wulianyang := result.RecentPercent[0] > 0 && result.RecentPercent[1] > 0 && result.RecentPercent[2] > 0 && result.RecentPercent[3] > 0 && result.RecentPercent[4] > 0 && result.RecentPercent[5] > 0

	// 近期长上影
	changshangying := craw.HasTopLine(result, 1)
	// 近期长下影
	changxiaying := craw.HasLowLine(result, 1)

	// 优良概念
	goodconcept := GetConceptByCode(code, "预盈预增|业绩预升|高派息|独角兽|高送转|基金重仓|QFII|RQFII")
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
	priceaboveave6 := result.RecentClose[0] > result.AveDailyPrice1[0]
	priceaboveave15 := result.RecentClose[0] > result.AveDailyPrice2[0]
	priceaboveave30 := result.RecentClose[0] > result.AveDailyPrice3[0]
	pricelowave6 := result.RecentClose[0] < result.AveDailyPrice1[0]
	pricelowave15 := result.RecentClose[0] < result.AveDailyPrice2[0]
	pricelowave30 := result.RecentClose[0] < result.AveDailyPrice3[0]

	// 成交过亿
	guoyi := result.CurrTotalMoney > 10000

	// 近期高振幅
	gaozhenfu := result.RecentAmplitude[0] > 15 || result.RecentAmplitude[1] > 15 || result.RecentAmplitude[2] > 15

	// 高换手
	huanshou := result.RecentTurnoverRate[0] > 30 || result.RecentTurnoverRate[1] > 30 || result.RecentTurnoverRate[2] > 30

	// 现金流量表 非负
	up1, up2, up3, up4 := GetUpManageCashFlow(code)
	// 利润表 不能为负
	pup1, pup2 := GetProfitNotMiner(code)
	// 资产负债表
	lup1, done1 := GetUpLiabilities(code)

	// 近5日资金净流入综合非负
	netflow := calc3(result.RecentNetFlow[0:4])
	// 近5日主力资金净流入综合非负
	mainnetflow := calc3(result.RecentMainNetFlow[0:4])

	x := []interface{}{liangnengbuduanbigger, yiziban, netflow, mainnetflow, pup1, pup2, huanshou, pricelowave6, pricelowave15, pricelowave30, priceaboveave15, priceaboveave30, priceaboveave6, guoyi, jigouchicangcount, jincha1, jincha2, jincha3, jincha4, jincha5, jincha6, jincha7, jincha8, priceshangyang1, priceshangyang2, priceshangyang3, priceshangyang4, priceshangyang5, gaoweihuitiao1, gaoweihuitiao2, gaoweihuitiao3, liangshangyang1, liangshangyang2, liangnengbigger1, liangnengbigger2, wulianyang, changshangying, changxiaying, goodconcept, simuchicangcount, junjialianhe1, zhangting}
	xx(x)

	// 1 K线形态
	// 2 量价形态
	// 3 财务数据
	// 4 机构持仓情况
	// 5 其它
	// 自定义数值

	cond_str := ""
	if priceaboveave6 {
		cond_str += fmt.Sprintf("收盘价在6日均线上方; ")
	}
	if priceaboveave15 {
		cond_str += fmt.Sprintf("收盘价在15日均线上方; ")
	}
	if priceaboveave30 {
		cond_str += fmt.Sprintf("收盘价在30日均线上方; ")
	}
	if pricelowave6 {
		cond_str += fmt.Sprintf("收盘价在6日均线下方; ")
	}
	if pricelowave15 {
		cond_str += fmt.Sprintf("收盘价在15日均线下方; ")
	}
	if pricelowave30 {
		cond_str += fmt.Sprintf("收盘价在30日均线下方; ")
	}
	if jincha1 || jincha2 {
		cond_str += "6日均线与15日均线交金叉; "
	}
	if jincha3 || jincha4 {
		cond_str += "15日均线与30日均线交金叉; "
	}
	if jincha5 || jincha6 {
		cond_str += "6周均线与15周均线交金叉; "
	}
	if priceshangyang1 {
		cond_str += "6日均线上扬; "
	}
	if priceshangyang2 {
		cond_str += "15日均线上扬; "
	}
	if priceshangyang3 {
		cond_str += "30日均线上扬; "
	}
	if priceshangyang4 {
		cond_str += "6周均线上扬; "
	}
	if priceshangyang5 {
		cond_str += "15周均线上扬; "
	}
	if changshangying {
		cond_str += "长上影; "
	}
	if changxiaying {
		cond_str += "长下影; "
	}
	if zhangting {
		cond_str += "涨停股; "
	}
	if yiziban {
		cond_str += "一字板; "
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
	if gaozhenfu {
		cond_str += "近期股价振幅高达15%; "
	}
	if wulianyang {
		cond_str += "五连阳; "
	}
	// 量价
	if jincha7 || jincha8 {
		cond_str += "10日与40日量能均线交金叉; "
	}
	if liangshangyang1 {
		cond_str += "连续5日量能10日均线上扬; "
	}
	if liangshangyang2 {
		cond_str += "连续5日量能40日均线上扬; "
	}
	if liangnengbigger1 {
		cond_str += "连续5日量能站上10日均线; "
	}
	if liangnengbigger2 {
		cond_str += "连续5日量能站上40日均线; "
	}
	if huanshou {
		cond_str += "近期成交量换手率高达30%; "
	}
	if liangnengbuduanbigger {
		cond_str += "量能不断放大; "
	}
	if tufangjuliang {
		cond_str += "突放巨量; "
	}
	if guoyi {
		cond_str += "成交额过亿; "
	}
	// 基本面 现金流量表
	if up1 {
		cond_str += "经营现金流量净额非负; "
	}
	if up2 {
		cond_str += "投资现金流量净额非负; "
	}
	if up3 {
		cond_str += "筹资现金流量净额非负; "
	}
	if up4 {
		cond_str += "期末现金及现金等价物余额非负; "
	}
	// 基本面 利润表
	if pup1 {
		cond_str += "营业总收入非负; "
	}
	if pup2 {
		cond_str += "净利润非负; "
	}
	// 基本面 资产负债表
	if lup1 {
		cond_str += "总资产不断增加; "
	}
	if done1 {
		cond_str += "总负债不断减小; "
	}
	// 机构持仓情况
	// 十大流通股东信息
	if simuchicangcount > 0 {
		cond_str += fmt.Sprintf("%d个私募持仓; ", simuchicangcount)
	}
	if jigouchicangcount > 0 {
		cond_str += fmt.Sprintf("%d个基金持仓; ", jigouchicangcount)
	}

	// 其它
	if gaoweihuitiao1 || gaoweihuitiao2 || gaoweihuitiao3 {
		cond_str += "高位回调; "
	}
	if goodconcept {
		cond_str += "优良概念; "
	}
	if netflow {
		cond_str += "近5日资金净流入总和非负; "
	}
	if mainnetflow {
		cond_str += "近5日主力资金净流入总和非负; "
	}

	if len(cond_str) > 0 {
		fmt.Println(code, name, cond_str)
		p := dal.Predict{Code: code, Name: name, Condition: cond_str, Date: result.CurrDate, FundCount: jigouchicangcount, SMCount: simuchicangcount}
		store.MysqlClient.GetDB().Save(&p)
	}
}
