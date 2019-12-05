// @Contact:    huaxinrui
// @Time:       2019/7/5 下午4:54

package dynamic

import (
	"huaxinrui/stock/crawler/utils"
	d "huaxinrui/stock/dao"
	"math"
	"strconv"
	"strings"
)

// 周线抑价 count 表示几周
func WeekLimitedPrice(recent_percent_week []float64, percent float64, count int) bool {
	return ConditionLimitedPrice(recent_percent_week[0:count], percent)
}

// 日线抑价
func DayLimitedPrice(recent_percent_day []float64, percent float64, count int) bool {
	return ConditionLimitedPrice(recent_percent_day[0:count], percent)
}

// 近几天有涨停
// recent 表示天数
func TopRecentDays(recent_percent []float64, recent int) bool {
	return ConditionTop(recent_percent, recent)
}

// 长上影
func TopLineRecentDays(array_shou, array_high, array_kai, percent []float64, recent int) bool {
	return ConditionTopLine(array_shou, array_high, array_kai, percent, recent)
}

// 当前价格必须在6日均线上方
func CurrentPriceBiggerAve6(cur_money, ave_6 float64) bool {
	return cur_money > ave_6
}

// 连续几日量价和均线的关系
func RecentBiggerAve(recent_money_or_count, recent_ave []float64, recent int) bool {
	for i := 0; i < recent-1; i++ {
		if recent_money_or_count[i] < recent_ave[i] {
			return false
		}
	}
	return true
}

func GetConceptByCode(code, concept string) bool {
	var c int
	d.Backend.DB.Model(&d.Code{}).Where("code = ?", code).Where("concept regexp ?", concept).Count(&c)
	if c > 0 {
		return true
	}
	return false
}

func GetSimu(code, concept string) (string, bool) {
	var holder d.Stockholder
	var c int
	tmp := d.Backend.DB.Model(&d.Stockholder{}).Where("code = ?", code).Where("holder_name regexp ?", concept)
	tmp.Count(&c)
	if c == 1 {
		tmp.Find(&holder)
		return holder.Percent, true

	}
	return "", false
}

// 最近几日量价和均线有几日满足关系条件
func RecentBiggerAveWithCond(recent_money_or_count, recent_ave []float64, recent int, total int) bool {
	var tmp []int
	for i := 0; i < recent-1; i++ {
		if recent_money_or_count[i] > recent_ave[i] {
			tmp = append(tmp, 1)
		}
	}
	return len(tmp) >= total
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

// 近几日放巨量, 且指定百分比
func RecentBigCount(recent_count, recent_ave []float64, percent float64, recent, total int) bool {
	// percent = 1 就是两倍的意思
	var tmp []int
	for i := 0; i < recent-1; i++ {
		if (recent_count[i]-recent_ave[i])/recent_ave[i] > percent {
			tmp = append(tmp, 1)
		}
	}
	return len(tmp) >= total
}

// 没有所谓纯粹的技术可言
// 量能， 均线分析都存在缺陷
// 但是可以直接关注高位回调, 大涨股票, 可在未来一段时间内持续性关注
func GetTicketPredict(code, name string, offset int, debug bool) {
	if strings.HasPrefix(code, "3") {
		return
	}

	recent_count, recent_money, recent_money_kai, recent_money_high, recent_money_low, recent_percent, recent_money_weekly, recent_percent_weekly,
		recent_ave_price_1, recent_ave_price_2, recent_ave_price_3, recent_ave_price_1_weekly, recent_ave_price_2_weekly, recent_ave_count_1, recent_ave_count_2, date, total_money := utils.CalcResultWithDefined(code, offset, 6, 15, 30, 10, 40)
	if len(recent_count) != 50 || len(recent_money_weekly) != 20 {
		return
	}
	// 抑价10%(3周滞涨)
	cond := WeekLimitedPrice(recent_percent_weekly, 10, 3)
	// 抑价2%(两天滞涨)
	cond2 := DayLimitedPrice(recent_percent, 5, 3)
	// 近60天有
	cond3 := TopRecentDays(recent_percent, 60)
	// 当前价格必须在6日均线上方
	cond4 := CurrentPriceBiggerAve6(recent_money[0], recent_ave_price_1[0])
	// 量能-近5天有4天量能站上均线
	cond5 := RecentBiggerAveWithCond(recent_count, recent_ave_count_1, 8, 5)
	cond6 := RecentBiggerAveWithCond(recent_count, recent_ave_count_2, 8, 5)
	// 近7天有3天量能是40线的两倍以上
	cond7 := RecentBigCount(recent_count, recent_ave_count_2, 1.5, 8, 4)

	// 昨日线金叉
	cond8 := recent_ave_price_1[0] > recent_ave_price_2[0] && recent_ave_price_1[1] < recent_ave_price_2[1]
	// 前日线金叉
	cond9 := recent_ave_price_1[1] > recent_ave_price_2[1] && recent_ave_price_1[2] < recent_ave_price_2[2]
	// 上周线金叉
	cond10 := recent_ave_price_1_weekly[0] > recent_ave_price_2_weekly[0] && recent_ave_price_1_weekly[1] < recent_ave_price_2_weekly[1]
	// 前周金叉
	cond11 := recent_ave_price_1_weekly[1] > recent_ave_price_2_weekly[1] && recent_ave_price_1_weekly[2] < recent_ave_price_2_weekly[2]

	// ave_price_1 价格均线上扬
	cond12 := recent_ave_price_1[0] > recent_ave_price_1[1] && recent_ave_price_1[1] > recent_ave_price_1[2] && recent_ave_price_1[2] > recent_ave_price_1[3]
	// ave_price_2 价格均线上扬
	cond13 := recent_ave_price_2[0] > recent_ave_price_2[1] && recent_ave_price_2[1] > recent_ave_price_2[2] && recent_ave_price_2[2] > recent_ave_price_2[3]
	// ave_price_3 价格均线上扬
	cond14 := recent_ave_price_3[0] > recent_ave_price_3[1] && recent_ave_price_3[1] > recent_ave_price_3[2] && recent_ave_price_3[2] > recent_ave_price_3[3]
	// 周线价格上扬
	cond15 := recent_ave_price_1_weekly[0] > recent_ave_price_1_weekly[1] && recent_ave_price_1_weekly[1] > recent_ave_price_1_weekly[2]
	cond16 := recent_ave_price_2_weekly[0] > recent_ave_price_2_weekly[1] && recent_ave_price_2_weekly[1] > recent_ave_price_2_weekly[2]

	// 高位回调
	cond17 := recent_percent[0] < -6 && recent_percent[1] > 8 && (recent_percent[0]+recent_percent[1] < 3)
	cond18 := recent_percent[0] < 3 && recent_percent[1] < 3 && recent_percent[2] > 8 && (recent_percent[0]+recent_percent[1]+recent_percent[2] < 3)
	cond19 := recent_percent[0] < 3 && recent_percent[1] < 3 && recent_percent[2] < 3 && recent_percent[3] > 8 && (recent_percent[0]+recent_percent[1]+recent_percent[2]+recent_percent[3] < 3)
	cond20 := recent_percent[0] < 3 && recent_percent[1] < 3 && recent_percent[2] < 3 && recent_percent[3] < 3 && recent_percent[4] > 8 && (recent_percent[0]+recent_percent[1]+recent_percent[2]+recent_percent[3]+recent_percent[4] < 3)
	cond28 := recent_percent[0] < 3 && recent_percent[1] < 3 && recent_percent[2] < 3 && recent_percent[3] < 3 && recent_percent[4] < 3 && recent_percent[5] > 8 && (recent_percent[0]+recent_percent[1]+recent_percent[2]+recent_percent[3]+recent_percent[4]+recent_percent[5] < 3)

	// 量能均线上扬
	cond21 := recent_ave_count_1[0] > recent_ave_count_1[1] && recent_ave_count_1[1] > recent_ave_count_1[2] && recent_ave_count_1[2] > recent_ave_count_1[3] && recent_ave_count_1[3] > recent_ave_count_1[4] && recent_ave_count_1[4] > recent_ave_count_1[5] && recent_ave_count_1[5] > recent_ave_count_1[6] && recent_ave_count_1[6] > recent_ave_count_1[7]
	cond22 := recent_ave_count_2[0] > recent_ave_count_2[1] && recent_ave_count_2[1] > recent_ave_count_2[2] && recent_ave_count_2[2] > recent_ave_count_2[3] && recent_ave_count_2[3] > recent_ave_count_2[4] && recent_ave_count_2[4] > recent_ave_count_2[5] && recent_ave_count_2[5] > recent_ave_count_2[6] && recent_ave_count_2[6] > recent_ave_count_2[7]
	// 量能均线昨日金叉
	cond23 := recent_ave_count_1[0] > recent_ave_count_2[0] && recent_ave_count_1[1] < recent_ave_count_2[1]
	// 量能均线前日金叉
	cond24 := recent_ave_count_1[1] > recent_ave_count_2[1] && recent_ave_count_1[2] < recent_ave_count_2[2]
	// 量能大前日金叉
	cond32 := recent_ave_count_1[2] > recent_ave_count_2[2] && recent_ave_count_1[3] < recent_ave_count_2[3]
	cond33 := recent_ave_count_1[3] > recent_ave_count_2[3] && recent_ave_count_1[4] < recent_ave_count_2[4]
	// 量能连续上扬
	cond25 := recent_count[0] > recent_count[1] && recent_count[1] > recent_count[2] && recent_count[2] > recent_count[3] && recent_count[3] > recent_count[4]
	// 5连阳
	cond26 := recent_percent[0] > 0 && recent_percent[1] > 0 && recent_percent[2] > 0 && recent_percent[3] > 0 && recent_percent[4] > 0 && DayLimitedPrice(recent_percent, 5, 5)

	// 上影洗盘
	cond27 := TopLineRecentDays(recent_money, recent_money_high, recent_money_kai, recent_percent, 5) && DayLimitedPrice(recent_percent, 5, 2)
	// 成交额必须上5qw
	cond29 := total_money > 20000000
	// 未启动
	cond30 := recent_money[0] > recent_ave_price_1[0] && recent_money[0] > recent_ave_price_2[0] && recent_money[0] < recent_ave_price_3[0]
	// 预盈预增预升，高派息
	cond31 := GetConceptByCode(code, "预盈预增|业绩预升|高派息|独角兽|高送转|基金重仓|QFII|RQFII")
	// 私募持仓
	percent, cond35 := GetSimu(code, "募")
	// 连续放量
	cond34 := RecentBigCount(recent_count, recent_ave_count_2, 1, 5, 4) && WeekLimitedPrice(recent_percent_weekly, 15, 3)
	utils.DoNothing(percent, cond35, cond27, cond31, cond, cond2, cond3, cond30, cond29, cond4, cond28, cond5, cond6, cond7, cond8, cond9, cond10, cond11, cond12, cond13, cond14, cond15, cond16, cond17, cond18, cond19, cond20, cond21, cond22, cond23, cond24, cond25, cond26, recent_count, recent_money, recent_money_kai, recent_money_high, recent_money_low, recent_percent, recent_money_weekly, recent_percent_weekly,
		recent_ave_price_1, recent_ave_price_2, recent_ave_price_3, recent_ave_price_1_weekly, recent_ave_price_2_weekly, recent_ave_count_1, recent_ave_count_2)

	cond_str := ""

	must := cond4 && cond29 && TopRecentDays(recent_percent, 50)

	if must && cond12 && cond13 && cond30 && recent_money[0] > recent_ave_price_1_weekly[0] && cond3 {
		cond_str += "未启动; "
	}

	// 6均线粘合上扬/ 15均线粘合上扬
	if must && ((RecentInRangeAveWithCond(recent_money, recent_ave_price_1, 7, 4) && cond12) || (RecentInRangeAveWithCond(recent_money, recent_ave_price_2, 7, 4) && cond13)) {
		cond_str += "均价粘合上扬; "
	}

	if cond && cond2 && cond5 && cond6 && cond7 {
		cond_str += "放量滞涨; "
	}

	if cond29 && TopRecentDays(recent_percent, 15) && (cond17 || cond18 || cond19 || cond20 || cond28) {
		cond_str += "高位回调; "
	}

	if must && RecentBigCount(recent_count, recent_ave_count_2, 1.8, 10, 6) {
		cond_str += "放巨量; "
	}

	// 量价均线全部上扬
	if must && (cond12 && cond13 && cond14 && cond15 && cond16) && (cond21 && cond22) {
		cond_str += "量价均线上扬; "
	}

	if must && (cond23 || cond24 || cond32 || cond33) && (cond8 || cond9 || cond10 || cond11) {
		cond_str += "金叉股; "
	}

	if must && cond26 {
		cond_str += "5连阳抑价; "
	}

	if must && cond27 && recent_percent[0] > 3 {
		cond_str += "上影洗盘; "
	}

	if recent_percent[0] > 9.94 {
		cond_str += "涨停股; "
	}

	if cond34 {
		cond_str += "连续放量; "
	}

	if cond31 {
		if cond_str != "" {
			cond_str += "优良概念; "
		}
	}

	if cond35 {
		if cond_str != "" {
			cond_str += "私募持仓:" + percent
		}
	}

	if cond_str != "" {
		if debug {
			var code_obj d.Code
			d.Backend.DB.Model(&d.Code{}).Where("code = ?", code).Find(&code_obj)
			var holders []d.Stockholder
			d.Backend.DB.Model(&d.Stockholder{}).Where("code = ?", code).Find(&holders)
			var str string
			for _, k := range holders {
				x := k.Change + ";"
				str += x
			}
			var count int
			d.Backend.DB.Model(&d.FundHoldRank{}).Where("code = ?", code).Count(&count)
			if count > 0 {
				cond_str += strconv.Itoa(count) + "机构持仓; "
			}
			p := d.PredictDebug{Code: code, Date: date, Condition: cond_str, Name: name, Attrs: code_obj.Concept, Holder: str, FundCount: count}
			d.Backend.DB.Save(&p)
		} else {
			p := d.Predict{Code: code, Date: date, Condition: cond_str, Name: name}
			d.Backend.DB.Save(&p)
		}

	}
}

// 开开心心投资，快快乐乐做人！股市人生，无论你赚再多的钱，如果你根本不开心那么一切都变得没有意义。压抑、悲伤、愤怒，只会让自己连续不断犯错误。
//但是如果你任何操作前，把自己当成主力，这样的走势和大环境，我会怎么想，我会怎么做，这就如同反侦察，你洞穿了真相，你能看懂这个上涨是真实的还是为了下跌准备的，这个下跌是真实的卖出，还是为上涨准备
func GetTicketByConditions(code string) {
	count_10_0, count_10_1, count_10_2, count_10_3, count_10_4, count_40_0, count_40_1, count_40_2, count_40_3, count_40_4, ava_6_0, ava_6_1, ava_6_2, ava_6_3, ava_15_0, ava_15_1, ava_15_2, ava_15_3, _, _, ava_6_0_week,
		ava_6_1_week, ava_6_2_week, ava_15_0_week, ava_15_1_week, ava_15_2_week, date, recent_count, recent_money, _, recent_money_high, recent_percent, recent_percent_week := utils.CalcResult(code)
	if len(recent_count) != 40 || len(recent_money) != 40 {
		return
	}
	// 抑价10%(三周滞涨)
	cond := WeekLimitedPrice(recent_percent_week, 10, 3)
	// 抑价2%(两天滞涨)
	cond2 := DayLimitedPrice(recent_percent, 5, 2)
	// 近五天没有涨停
	cond3 := TopRecentDays(recent_percent, 5)
	// 当前价格必须在6日均线上方
	condition := CurrentPriceBiggerAve6(recent_money[0], ava_6_0)

	// 量能-近5天有4天量能站上均线
	condition1, condition2 := Condition(recent_count, count_10_0, count_10_1, count_10_2, count_10_3, count_10_4, count_40_0, count_40_1, count_40_2, count_40_3, count_40_4)
	// 滞涨15%
	if condition1 && condition2 && cond && condition {
		p := d.Predict{Code: code, Date: date, Condition: "放量滞涨"}
		d.Backend.DB.Save(&p)
	}
	// 量能-近5天有2天量能是40线的两倍以上, 两天不涨， 五天内无涨停
	condition4 := Condition2(recent_count, count_40_0, count_40_1, count_40_2, count_40_3, count_40_4)
	if condition4 && condition && cond2 && cond3 {
		p := d.Predict{Code: code, Date: date, Condition: "放巨量"}
		d.Backend.DB.Save(&p)
	}
	// 6日, 15均线必须上扬, 6 周均线上扬
	condition3 := ava_6_0 > ava_6_1 && ava_6_1 > ava_6_2 && ava_6_2 > ava_6_3 && ava_15_0 > ava_15_1 && ava_15_1 > ava_15_2 && ava_15_2 > ava_15_3 && (ava_6_0_week >= ava_6_1_week && ava_6_1_week >= ava_6_2_week)
	if condition3 && condition && cond {
		p := d.Predict{Code: code, Date: date, Condition: "均线上扬"}
		d.Backend.DB.Save(&p)
	}
	// 前日均线金叉
	condition6 := ava_6_0 > ava_15_0 && ava_6_1 > ava_15_1 && ava_6_2 < ava_15_2
	// 昨日均线金叉
	condition7 := ava_6_0 > ava_15_0 && ava_6_1 < ava_15_1
	if condition && (condition6 || condition7) && cond {
		p := d.Predict{Code: code, Date: date, Condition: "日线金叉"}
		d.Backend.DB.Save(&p)
	}
	// 本周线金叉（本周的6均大于15均， 上周的6均小于15均, 并且6均上扬)
	condition8 := ava_6_0_week > ava_15_0_week && ava_6_1_week < ava_15_1_week && ava_6_0_week > ava_6_1_week
	// 上周线金叉
	condition9 := ava_6_0_week > ava_15_0_week && ava_6_1_week > ava_15_1_week && ava_6_2_week < ava_15_2_week && ava_6_0_week > ava_6_1_week
	if condition && (condition8 || condition9) && cond {
		p := d.Predict{Code: code, Date: date, Condition: "周线金叉"}
		d.Backend.DB.Save(&p)
	}
	// 涨停回调(近三天回调，回调吃掉10%以上)
	if (recent_percent[6] > 9.5 && recent_money[6] == recent_money_high[6]) && (recent_money[0]-recent_money[6])/recent_money[6] < -0.095 {
		p := d.Predict{Code: code, Date: date, Condition: "涨停回调"}
		d.Backend.DB.Save(&p)
	}
	if (recent_percent[5] > 9.5 && recent_money[5] == recent_money_high[5]) && (recent_money[0]-recent_money[5])/recent_money[5] < -0.095 {
		p := d.Predict{Code: code, Date: date, Condition: "涨停回调"}
		d.Backend.DB.Save(&p)
	}
	if (recent_percent[4] > 9.5 && recent_money[4] == recent_money_high[4]) && (recent_money[0]-recent_money[4])/recent_money[4] < -0.095 {
		p := d.Predict{Code: code, Date: date, Condition: "涨停回调"}
		d.Backend.DB.Save(&p)
	}

	utils.DoNothing(condition1, condition2, condition3, condition, condition6, condition7, condition8, condition9)
	if (condition1 && condition2) && condition3 && condition && cond {
		p := d.Predict{Code: code, Date: date, Condition: "综合"}
		d.Backend.DB.Save(&p)
	}
}
