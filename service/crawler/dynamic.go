// @Time:       2019/12/2 下午2:24

package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/utils"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/parnurzeal/gorequest"
)

func dealPanic(doc *goquery.Document, i int, selector string) (str string, err error) {
	defer func() {
		if errs := recover(); errs != nil {
			err = errors.New("err")
		}
	}()
	str = doc.Find(fmt.Sprintf(selector, i)).Text()
	return
}

// 每天收盘执行一次, 收集所有股票的当天股价和成交量
func (craw *Crawler) GetAllTicketTodayDetail(code, name, today string, proxy bool) error {
	var doc *goquery.Document
	for year := 2019; year <= 2019; year++ { // 马上 2020 了
		for ji := 4; ji <= 4; ji++ { // 这里需要动态改
			if !proxy {
				doc, _ = craw.NewDocument(fmt.Sprintf("http://quotes.money.163.com/trade/lsjysj_%s.html?year=%d&season=%d", code, year, ji))
			} else {
				doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://quotes.money.163.com/trade/lsjysj_%s.html?year=%d&season=%d", code, year, ji))
			}
			for i := 1; i >= 1; i-- { //for i := 73; i >= 1 ; i --
				// 一季度按照72个交易日来算  // 因为第一栏是中文标题
				text := ""
				// 获取日期不要用
				//x := doc.Find(fmt.Sprintf("body > div.area > div.inner_box > table > tbody > tr:nth-child(%d) > td:nth-child(1)", i)).Text()
				x, err := dealPanic(doc, i, "body > div.area > div.inner_box > table > tbody > tr:nth-child(%d) > td:nth-child(1)")
				if err != nil {
					// 空指针错误
					return err
				}
				if len(x) == 0 {
					continue
				}
				if x != today {
					continue
				}
				text += strings.TrimSpace(x) + ","
				for j := 2; j <= 11; j++ {
					doc.Find(fmt.Sprintf("body > div.area > div.inner_box > table > tbody > tr:nth-child(%d) > td:nth-child(%d)", i, j)).Each(func(xxoo int, selection *goquery.Selection) {
						t := selection.Text()
						if len(t) != 0 {
							t := strings.Replace(t, ",", "", -1)
							text += t + ","
						}
					})
				}
				if len(text) != 0 {
					tmp := strings.Split(text, ",")
					date := tmp[0]
					kai, _ := strconv.ParseFloat(tmp[1], 64)
					high, _ := strconv.ParseFloat(tmp[2], 64)
					low, _ := strconv.ParseFloat(tmp[3], 64)
					shou, _ := strconv.ParseFloat(tmp[4], 64)
					zhangdiee, _ := strconv.ParseFloat(tmp[5], 64)
					percent, _ := strconv.ParseFloat(tmp[6], 64)
					tc, _ := strconv.ParseFloat(strings.Replace(tmp[7], ",", "", -1), 64)
					tm, _ := strconv.ParseFloat(strings.Replace(tmp[8], ",", "", -1), 64)
					zhenfu, _ := strconv.ParseFloat(tmp[9], 64)
					huanshou, _ := strconv.ParseFloat(tmp[10], 64)
					dh := dal.TicketHistory{Code: code, Name: name, Kai: kai, High: high, Low: low, Shou: shou, TotalCount: tc, TotalMoney: tm, Date: date,
						Percent: percent, Change: zhangdiee, Amplitude: zhenfu, TurnoverRate: huanshou} //, Percent:p}
					if utils.TellEnv() == "loc" {
						err := store.MysqlClient.GetOnlineDB().Save(&dh).Error
						if err != nil {
							log.Println("写入线上错误")
						}
					}
					store.MysqlClient.GetDB().Save(&dh)
					fmt.Println(code, name, date, kai, high, low, shou, zhangdiee, percent, tc, tm, zhenfu, huanshou)

				}
			}
		}
	}
	return nil
}

// 把今日的收盘价 加入到周线均价的表中
func (craw *Crawler) AddTodayShouToWeek(code, last_week, last_day_to_delete, today string) {
	if len(last_day_to_delete) != 0 {
		store.MysqlClient.GetDB().Exec("delete from ticket_history_week where date = ?", last_day_to_delete)
	}
	var th dal.TicketHistory
	err := store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ? and date = ?", code, today).Find(&th).Error
	if err != nil {
		log.Println("数据不存在", code, today)
	}

	var xx dal.TicketHistoryWeekly
	err = store.MysqlClient.GetDB().Model(&dal.TicketHistoryWeekly{}).Where("code = ? and date = ?", code, last_week).Find(&xx).Error
	if err != nil {
		log.Println("数据不存在2", code, today)
	}

	p, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (th.Shou-xx.Shou)*100/xx.Shou), 64)

	weekly := dal.TicketHistoryWeekly{Code: th.Code, Date: th.Date, Name: th.Name, Shou: th.Shou, Percent: p}
	if utils.TellEnv() == "loc" {
		err := store.MysqlClient.GetOnlineDB().Save(&weekly).Error
		if err != nil {
			log.Println("写入线上错误")
		}
	}
	store.MysqlClient.GetDB().Save(&weekly)
}

func msToTime(ms string) (string, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return "", err
	}

	tm := time.Unix(0, msInt*int64(time.Millisecond))
	return tm.Format("2006-01-02"), nil
}

// 获取短期回报基金排行前300名
func (craw *Crawler) GetFundRanks() {
	store.MysqlClient.GetDB().Exec("truncate table magic_stock_fund_rank")
	var m = map[string]int{"股票型": 1, "指数型": 4, "混合型": 3, "QDII": 6, "LOF": 7}
	for k, v := range m {
		_, body, _ := gorequest.New().Get(fmt.Sprintf("https://fund.jrj.com.cn/json/netrank/open?type=%d&mana=0&limit=500&page=1&sort=6&order=1&vname=fundranklist", v)).End()
		if len(body) == 0 {
			return
		}
		body = strings.TrimLeft(body, "var fundranklist=")
		var m model.Result
		json.Unmarshal([]byte(body), &m)

		for _, i := range m.List {
			t, _ := msToTime(strconv.FormatInt(i.Time, 10))
			var count int
			store.MysqlClient.GetDB().Model(&dal.FundRank{}).Where("fund_code = ? and time = ?", i.FundCode, t).Count(&count)
			if count > 0 {
				continue
			}
			x := dal.FundRank{Type: k, FundCode: i.FundCode, FundSName: i.FundSName, LastMonth: i.LastMonth, LastWeek: i.LastWeek, Last3Month: i.Last3Month,
				Last6Month: i.Last6Month, SinceBase: i.SinceBase, ThisYear: i.ThisYear, LastYear: i.LastYear, Last2Year: i.Last2Year, Last3Year: i.Last3Year, Time: t}
			store.MysqlClient.GetDB().Save(&x)
		}
	}
}

// 通过上面获取重仓股
func (craw *Crawler) GetFundHighHold(date string) {
	store.MysqlClient.GetDB().Exec("truncate table magic_stock_fund_hold")
	var ranks []dal.FundRank
	store.MysqlClient.GetDB().Model(&dal.FundRank{}).Where("time = ?", date).Find(&ranks)

	tmp := map[string]bool{}

	for _, r := range ranks {
		if _, ok := tmp[r.FundCode]; ok {
			continue
		} else {
			tmp[r.FundCode] = true
		}
		_, body, _ := gorequest.New().Get(fmt.Sprintf("https://fund.jrj.com.cn/archives,%s.shtml", r.FundCode)).End()
		if len(body) == 0 {
			fmt.Println("NOT SUPPORT1,", r.FundCode)
			continue
		}
		x := strings.Split(body, "var fundjjzcg = ")
		if len(x) != 2 {
			x = strings.Split(body, "var fundjjzcg=")
			if len(x) != 2 {
				fmt.Println("NOT SUPPORT2,", r.FundCode)
			}
		}
		xx := strings.Split(x[1], "</script>")
		xxx := xx[0]

		var m []model.Hold
		json.Unmarshal([]byte(xxx), &m)

		for _, i := range m {
			var code dal.Code
			store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", i.Code).Find(&code)
			a := dal.FundHoldRank{Type: r.Type, FundCode: r.FundCode, FundSName: r.FundSName, Code: i.Code, Name: code.Name, Percent: i.Percent, Time: date}
			store.MysqlClient.GetDB().Save(&a)
		}
	}
	store.MysqlClient.GetDB().Exec("delete from magic_stock_fund_hold where name=''")
}

//// 每天收盘执行一次, 收集所有股票的当天股价和成交量  http://hq.sinajs.cn/list=sz000001
//func (craw *Crawler) GetAllTicketTodayDetail(code, name, last_day, today string) bool {
//	coder := code
//	if strings.HasPrefix(code, "6") {
//		code = "sh" + code //上海
//	} else {
//		code = "sz" + code //深圳
//	}
//	_, body, _ := gorequest.New().Get("http://hq.sinajs.cn/list=" + code).End()
//	if len(body) == 0 {
//		log.Println("获取新浪数据失败")
//		return false
//	}
//	x := strings.Split(string(body), ",")
//	if len(x) == 1 {
//		log.Println("数据解析错误")
//		return false
//	}
//
//	jinkai, _ := strconv.ParseFloat(x[1], 64)
//	jinshou, _ := strconv.ParseFloat(x[3], 64)
//	jinhigh, _ := strconv.ParseFloat(x[4], 64)
//	jinlow, _ := strconv.ParseFloat(x[5], 64)
//	t_count, _ := strconv.ParseFloat(x[8], 64)
//	t_money, _ := strconv.ParseFloat(x[9], 64)
//	date := x[30]
//	if date != today {
//		log.Println("时间错误")
//		return false
//	}
//	var xx dal.TicketHistory
//	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ? and date = ?", coder, last_day).Find(&xx)
//	p, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (jinshou-xx.Shou)*100/xx.Shou), 64)
//
//	th := dal.TicketHistory{Date: date, Name: name, Code: coder, Kai: jinkai, High: jinhigh,
//		Shou: jinshou, Low: jinlow, TotalCount: t_count, TotalMoney: t_money, Percent: p}
//	store.MysqlClient.GetDB().Create(&th)
//	return true
//}
