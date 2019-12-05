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
	Num      = 4
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
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", code).Where("concept regexp ?", concept).Count(&c)
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
	liangnengbigger1 := result.AveCount1[0] > result.RecentCount[0] && result.AveCount1[1] > result.RecentCount[1] && result.AveCount1[2] > result.RecentCount[2] && result.AveCount1[3] > result.RecentCount[3] && result.AveCount1[4] > result.RecentCount[4]
	// 连续5日量能站上40均线
	liangnengbigger2 := result.AveCount2[0] > result.RecentCount[0] && result.AveCount2[1] > result.RecentCount[1] && result.AveCount2[2] > result.RecentCount[2] && result.AveCount2[3] > result.RecentCount[3] && result.AveCount2[4] > result.RecentCount[4]

	// 5连阳
	wulianyang := result.RecentPercent[0] > 0 && result.RecentPercent[1] > 0 && result.RecentPercent[2] > 0 && result.RecentPercent[3] > 0 && result.RecentPercent[4] > 0 && result.RecentPercent[5] > 0

	// 近期长上影
	changshangying := craw.HasTopLine(result, 5)
	// 近期长下影
	changxiaying := craw.HasLowLine(result, 5)

	// 优良概念
	goodconcept := GetConceptByCode(code, "预盈预增|业绩预升|高派息|独角兽|高送转|基金重仓|QFII|RQFII")
	// 私募持仓
	simuchicangcount := GetHolderByCode(code, "私募")
	// 公募持仓
	gongmuchicangcount := GetHolderByCode(code, "公募")
	// 基金持仓
	jigouchicangcount := GetFundByCode(code)

	// 均价粘合上扬
	junjialianhe := RecentInRangeAveWithCond(result.RecentClose, result.AveDailyPrice1, 7, 4)

	// 涨停股
	zhangting := result.RecentPercent[0] > 9.94

	// 当前价格在短期均线上方
	priceaboveave := result.RecentClose[0] > result.AveDailyPrice1[0]

	// 成交过亿
	guoyi := result.CurrTotalMoney > 10000

	x := []interface{}{priceaboveave, guoyi, jigouchicangcount, jincha1, jincha2, jincha3, jincha4, jincha5, jincha6, jincha7, jincha8, priceshangyang1, priceshangyang2, priceshangyang3, priceshangyang4, priceshangyang5, gaoweihuitiao1, gaoweihuitiao2, gaoweihuitiao3, liangshangyang1, liangshangyang2, liangnengbigger1, liangnengbigger2, wulianyang, changshangying, changxiaying, goodconcept, simuchicangcount, junjialianhe, zhangting}
	xx(x)

	cond_str := ""
	if priceaboveave {
		cond_str += "当前价格在短期均线上方; "
	}
	if guoyi {
		cond_str += "成交额过亿; "
	}

	if jincha1 || jincha2 {
		cond_str += "短中均价金叉; "
	}
	if jincha3 || jincha4 {
		cond_str += "中长均价金叉; "
	}
	if jincha5 || jincha6 {
		cond_str += "周线均价金叉; "
	}
	if jincha7 || jincha8 {
		cond_str += "量能金叉; "
	}
	if priceshangyang1 && priceshangyang2 && priceshangyang3 {
		cond_str += "均线全部上扬; "
	}
	if priceshangyang4 && priceshangyang5 {
		cond_str += "短中周均线上扬; "
	}
	if gaoweihuitiao1 || gaoweihuitiao2 || gaoweihuitiao3 {
		cond_str += "高位回调; "
	}
	if liangshangyang1 || liangshangyang2 {
		cond_str += "连续5日量能10均线上扬; "
	}
	if liangnengbigger1 {
		cond_str += "连续5日量能站上十日均线; "
	}
	if liangnengbigger2 {
		cond_str += "连续5日量能站上四十日均线; "
	}
	if wulianyang {
		cond_str += "五连阳; "
	}
	if changshangying {
		cond_str += "长上影; "
	}
	if changxiaying {
		cond_str += "长下影; "
	}
	if goodconcept {
		cond_str += "优良概念; "
	}
	if simuchicangcount > 0 {
		cond_str += fmt.Sprintf("%d个私募持仓; ", simuchicangcount)
	}
	if gongmuchicangcount > 0 {
		cond_str += fmt.Sprintf("%d个公募持仓; ", gongmuchicangcount)
	}
	if jigouchicangcount > 0 {
		cond_str += fmt.Sprintf("%d个机构持仓; ", jigouchicangcount)
	}
	if junjialianhe {
		cond_str += "均价粘合上扬; "
	}
	if zhangting {
		cond_str += "涨停股; "
	}

	if len(cond_str) > 0 {
		fmt.Println(code, name, cond_str)
		p := dal.Predict{Code: code, Name: name, Condition: cond_str, Date: result.CurrDate, FundCount: jigouchicangcount, SMCount: simuchicangcount, GMCount: gongmuchicangcount}
		store.MysqlClient.GetDB().Save(&p)
	}

}
