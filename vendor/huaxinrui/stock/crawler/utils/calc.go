// @Contact:    huaxinrui
// @Time:       2019/7/5 下午4:59

package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	d "huaxinrui/stock/dao"
)

func DoNothing(x ...interface{}) {

}

func getRecentDailyData(ths []d.TicketHistory) ([]float64,[]float64,[]float64,[]float64,[]float64,[]float64, string, float64){
	var recent_count,recent_money,recent_money_kai,recent_money_high,recent_money_low,recent_percent []float64
	var date string
	var total_money float64
	for i, th := range ths {
		// 近40天数据临时存留计算用
		if i == 0 {
			date = th.Date
			total_money = th.TotalMoney
		}
		recent_count = append(recent_count, th.TotalCount)
		recent_money = append(recent_money, th.Shou)
		recent_money_kai = append(recent_money_kai, th.Kai)
		recent_money_high = append(recent_money_high, th.High)
		recent_money_low = append(recent_money_low, th.Low)
		recent_percent = append(recent_percent, th.Percent)
	}
	return recent_count,recent_money,recent_money_kai,recent_money_high,recent_money_low,recent_percent, date, total_money
}

func getRecentWeeklyData(thw []d.TicketHistoryWeekly) ([]float64,[]float64){
	var recent_money, recent_percent []float64
	for _, th := range thw {
		// 近40天数据临时存留计算用
		recent_money = append(recent_money, th.Shou)
		recent_percent = append(recent_percent, th.Percent)
	}
	return recent_money, recent_percent
}

// num 表示均线值, count 表示取几条数据
func getRecentAvePriceByNum(recent_money []float64, num int, count int) []float64{
	// 均线最大取40
	var ave []float64
	for i := 0; i < count; i ++ {
		var tmp float64
		slice := recent_money[i: num + i]
		for _, s := range slice {
			tmp += s
		}
		ave = append(ave, tmp / float64(num))
	}
	return ave
}

// 用户自定义均线参数计算, ave_price_1, ave_price_2, ave_price_3 默认三条价格均线， 2条两年能均线
// 需要返回最近的价格， 近几日均线等数据
func CalcResultWithDefined(code string, offset int,  ave_price_1, ave_price_2, ave_price_3 int, ave_count_1, ave_count_2 int) ([]float64,[]float64,[]float64,[]float64,[]float64,[]float64,[]float64,[]float64,[]float64,[]float64,[]float64,[]float64,[]float64,[]float64,[]float64, string, float64){

	if ave_price_1 > 6 || ave_price_2 > 15 || ave_price_3 > 30 || ave_count_1 > 40 || ave_count_2 > 40{
		panic("argument error")
	}

	var ths []d.TicketHistory
	d.Backend.DB.Model(&d.TicketHistory{}).Where("code = ?", code).Limit(50).Offset(offset).Order("date desc").Find(&ths)

	var thw []d.TicketHistoryWeekly
	d.Backend.DB.Model(&d.TicketHistoryWeekly{}).Where("code = ?", code).Limit(20).Offset(offset).Order("date desc").Find(&thw)

	recent_count,recent_money,recent_money_kai,recent_money_high,recent_money_low,recent_percent, date, total_money := getRecentDailyData(ths)
	recent_money_weekly, recent_percent_weekly := getRecentWeeklyData(thw)

	if len(recent_count) != 50 || len(recent_money_weekly) != 20 {
		return []float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{},[]float64{}, "", 0
	}

	recent_ave_price_1 := getRecentAvePriceByNum(recent_money, ave_price_1, 20)
	recent_ave_price_2 := getRecentAvePriceByNum(recent_money, ave_price_2, 20)
	recent_ave_price_3 := getRecentAvePriceByNum(recent_money, ave_price_3, 20)

	recent_ave_price_1_weekly := getRecentAvePriceByNum(recent_money_weekly, ave_price_1, 5)
	recent_ave_price_2_weekly := getRecentAvePriceByNum(recent_money_weekly, ave_price_2, 5)

	recent_ave_count_1 := getRecentAvePriceByNum(recent_count, ave_count_1, 10)
	recent_ave_count_2 := getRecentAvePriceByNum(recent_count, ave_count_2, 10)

	return recent_count,recent_money,recent_money_kai,recent_money_high,recent_money_low,recent_percent, recent_money_weekly, recent_percent_weekly,
		recent_ave_price_1, recent_ave_price_2, recent_ave_price_3, recent_ave_price_1_weekly, recent_ave_price_2_weekly, recent_ave_count_1, recent_ave_count_2, date, total_money
}


func Calc1(doc *goquery.Document, penny int) float64 {
	var total_count float64
	for i := 0; i <= penny; i++ {
		for j := 1; j <= 3; j++ {
			x := doc.Find(fmt.Sprintf("#datalist > tbody > tr:nth-child(%d) > td:nth-child(%d)", i, j)).Text()
			if len(x) == 0 {
				continue
			}
			if j == 2 {
				m, _ := strconv.ParseFloat(x, 64)
				t, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", m/float64(10000)), 64)
				total_count += t
			}
		}
	}
	return total_count
}

func Calc2(doc *goquery.Document, tt float64, penny int) [][]float64 {
	var tmps [][]float64
	for i := 0; i <= penny; i++ {
		var tmp []float64
		for j := 1; j <= 3; j++ {
			x := doc.Find(fmt.Sprintf("#datalist > tbody > tr:nth-child(%d) > td:nth-child(%d)", i, j)).Text()
			if len(x) == 0 {
				continue
			}

			if j == 1 {
				m, _ := strconv.ParseFloat(x, 64)
				t, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", m), 64)
				tmp = append(tmp, t)
			}

			if j == 2 {
				m, _ := strconv.ParseFloat(x, 64)
				t, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", m/float64(10000)), 64)
				percent, _ := strconv.ParseFloat(fmt.Sprintf("%.5f", t/tt), 64)
				tmp = append(tmp, percent)
			}
		}

		if len(tmp) > 0 {
			tmps = append(tmps, tmp)
		}
	}
	return tmps
}