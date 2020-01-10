// @Time:       2019/12/2 下午5:24

package crawler

import (
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
)

// num 表示均线值, count 表示取几条数据
func (craw *Crawler) getRecentAvePriceByNum(recent_money []float64, num int, count int) []float64 {
	if len(recent_money) < count+num-1 {
		return nil
	}
	// 均线最大取40
	var ave []float64
	for i := 0; i < count; i++ {
		var tmp float64
		slice := recent_money[i : num+i]
		for _, s := range slice {
			tmp += s
		}
		ave = append(ave, tmp/float64(num))
	}
	return ave
}

func (craw *Crawler) getRecentDailyData(ths []dal.TicketHistory) *model.RecentDailyData {
	var recent_count, recent_money_shou, recent_money_kai, recent_money_high, recent_money_low, recent_percent, recent_amplitude, recent_turn_over []float64
	var date string
	var total_money float64
	for i, th := range ths {
		// 近40天数据临时存留计算用
		if i == 0 {
			date = th.Date
			total_money = th.TotalMoney
		}
		recent_count = append(recent_count, th.TotalCount)
		recent_money_shou = append(recent_money_shou, th.Shou)
		recent_money_kai = append(recent_money_kai, th.Kai)
		recent_money_high = append(recent_money_high, th.High)
		recent_money_low = append(recent_money_low, th.Low)
		recent_percent = append(recent_percent, th.Percent)
		recent_amplitude = append(recent_amplitude, th.Amplitude)    // 振幅
		recent_turn_over = append(recent_turn_over, th.TurnoverRate) // 换手率
	}
	return &model.RecentDailyData{recent_count, recent_money_shou, recent_money_kai, recent_money_high, recent_money_low,
		recent_percent, recent_amplitude, recent_turn_over, date, total_money}
}

func (craw *Crawler) getRecentWeeklyData(thw []dal.TicketHistoryWeekly) *model.RecentWeeklyData {
	var recent_count, recent_money_shou, recent_money_kai, recent_money_high, recent_money_low, recent_percent, recent_amplitude, recent_turn_over []float64
	for _, th := range thw {
		recent_count = append(recent_count, th.TotalCount)
		recent_money_shou = append(recent_money_shou, th.Shou)
		recent_money_kai = append(recent_money_kai, th.Kai)
		recent_money_high = append(recent_money_high, th.High)
		recent_money_low = append(recent_money_low, th.Low)
		recent_percent = append(recent_percent, th.Percent)
		recent_amplitude = append(recent_amplitude, th.Amplitude)    // 振幅
		recent_turn_over = append(recent_turn_over, th.TurnoverRate) // 换手率
	}
	return &model.RecentWeeklyData{recent_count, recent_money_shou, recent_money_kai, recent_money_high, recent_money_low,
		recent_percent, recent_amplitude, recent_turn_over}
}

func (craw *Crawler) CalcResultWithDefined(params *model.Params) *model.CalcResult {
	var ths []dal.TicketHistory
	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ? and date <= ?", params.Code, params.Date).Limit(70).Offset(params.Offset).Order("date desc").Find(&ths)

	var thw []dal.TicketHistoryWeekly
	store.MysqlClient.GetDB().Model(&dal.TicketHistoryWeekly{}).Where("code = ? and date <= ?", params.Code, params.Date).Limit(70).Offset(params.Offset).Order("date desc").Find(&thw)

	recent_daily := craw.getRecentDailyData(ths)
	recent_weekly := craw.getRecentWeeklyData(thw)
	//if recent_daily.CurrDate != params.Date {
	//	return nil
	//}
	if len(recent_daily.RecentCount) != 70 {
		log.Println("交易日次数不足")
		return nil
	}
	var recent_ave_price_1, recent_ave_price_2, recent_ave_price_3, recent_ave_price_4, recent_ave_count_1, recent_ave_count_2 []float64
	var recent_ave_price_weekly_1, recent_ave_price_weekly_2, recent_ave_price_weekly_3, recent_ave_price_weekly_4, recent_ave_count_week_1, recent_ave_count_week_2 []float64
	recent_ave_price_1 = craw.getRecentAvePriceByNum(recent_daily.RecentClose, params.AveragePrice1, 6)
	recent_ave_price_2 = craw.getRecentAvePriceByNum(recent_daily.RecentClose, params.AveragePrice2, 6)
	recent_ave_price_3 = craw.getRecentAvePriceByNum(recent_daily.RecentClose, params.AveragePrice3, 6)
	recent_ave_price_4 = craw.getRecentAvePriceByNum(recent_daily.RecentClose, params.AveragePrice4, 6)

	recent_ave_count_1 = craw.getRecentAvePriceByNum(recent_daily.RecentCount, params.AverageCount1, 6)
	recent_ave_count_2 = craw.getRecentAvePriceByNum(recent_daily.RecentCount, params.AverageCount2, 6)

	recent_ave_price_weekly_1 = craw.getRecentAvePriceByNum(recent_weekly.RecentCloseWeek, params.AveragePrice1, 6)
	recent_ave_price_weekly_2 = craw.getRecentAvePriceByNum(recent_weekly.RecentCloseWeek, params.AveragePrice2, 6)
	recent_ave_price_weekly_3 = craw.getRecentAvePriceByNum(recent_weekly.RecentCloseWeek, params.AveragePrice3, 6)
	recent_ave_price_weekly_4 = craw.getRecentAvePriceByNum(recent_weekly.RecentCloseWeek, params.AveragePrice4, 6)

	recent_ave_count_week_1 = craw.getRecentAvePriceByNum(recent_weekly.RecentCountWeek, params.AverageCount1, 6)
	recent_ave_count_week_2 = craw.getRecentAvePriceByNum(recent_weekly.RecentCountWeek, params.AverageCount2, 6)

	return &model.CalcResult{RecentDailyData: recent_daily,
		RecentWeeklyData:    recent_weekly,
		RecentAverage:       &model.RecentAverage{recent_ave_price_1, recent_ave_price_2, recent_ave_price_3, recent_ave_price_4, recent_ave_count_1, recent_ave_count_2},
		RecentAverageWeekly: &model.RecentAverageWeekly{recent_ave_price_weekly_1, recent_ave_price_weekly_2, recent_ave_price_weekly_3, recent_ave_price_weekly_4, recent_ave_count_week_1, recent_ave_count_week_2},
	}
}
