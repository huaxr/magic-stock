// @Time:       2019/12/2 下午5:24

package crawler

import (
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
)

// num 表示均线值, count 表示取几条数据
func (craw *Crawler) getRecentAvePriceByNum(recent_money []float64, num int, count int) []float64 {
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
	return &model.RecentDailyData{recent_count, recent_money_shou, recent_money_kai, recent_money_high, recent_money_low, recent_percent, recent_amplitude, recent_turn_over, date, total_money}
}

func (craw *Crawler) getRecentWeeklyData(thw []dal.TicketHistoryWeekly) *model.RecentWeeklyData {
	var recent_money, recent_percent []float64
	for _, th := range thw {
		// 近40天数据临时存留计算用
		recent_money = append(recent_money, th.Shou)
		recent_percent = append(recent_percent, th.Percent)
	}
	return &model.RecentWeeklyData{recent_money, recent_percent}
}

func (craw *Crawler) CalcResultWithDefined(params *model.Params) *model.CalcResult {
	if params.AveragePrice1 > 6 || params.AveragePrice2 > 15 || params.AveragePrice3 > 30 || params.AverageCount1 > 40 || params.AverageCount2 > 40 {
		panic("argument error")
	}
	var ths []dal.TicketHistory
	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ?", params.Code).Limit(50).Offset(params.Offset).Order("date desc").Find(&ths)

	var thw []dal.TicketHistoryWeekly
	store.MysqlClient.GetDB().Model(&dal.TicketHistoryWeekly{}).Where("code = ?", params.Code).Limit(20).Offset(params.Offset).Order("date desc").Find(&thw)

	recent_daily := craw.getRecentDailyData(ths)
	recent_week := craw.getRecentWeeklyData(thw)

	if len(recent_daily.RecentCount) != 50 || len(recent_week.RecentWeeklyClose) != 20 {
		return nil
	}

	recent_ave_price_1 := craw.getRecentAvePriceByNum(recent_daily.RecentClose, params.AveragePrice1, 20)
	recent_ave_price_2 := craw.getRecentAvePriceByNum(recent_daily.RecentClose, params.AveragePrice2, 20)
	recent_ave_price_3 := craw.getRecentAvePriceByNum(recent_daily.RecentClose, params.AveragePrice3, 20)

	recent_ave_price_1_weekly := craw.getRecentAvePriceByNum(recent_week.RecentWeeklyClose, params.AveragePrice1, 5)
	recent_ave_price_2_weekly := craw.getRecentAvePriceByNum(recent_week.RecentWeeklyClose, params.AveragePrice2, 5)

	recent_ave_count_1 := craw.getRecentAvePriceByNum(recent_daily.RecentCount, params.AverageCount1, 10)
	recent_ave_count_2 := craw.getRecentAvePriceByNum(recent_daily.RecentCount, params.AverageCount2, 10)

	return &model.CalcResult{recent_daily, recent_week, &model.RecentAverage{recent_ave_price_1, recent_ave_price_2, recent_ave_price_3, recent_ave_price_1_weekly, recent_ave_price_2_weekly, recent_ave_count_1, recent_ave_count_2}}
}
