// @Time:       2019/12/2 下午5:24

package crawler

import (
	"fmt"
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/utils"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/parnurzeal/gorequest"
)

// 爬取上证所有股票代码和名称(注意因为关联， 千万不要多次执行)
func (craw *Crawler) GetAllTicketCode() {
	for i := 10000; i <= 19999; i++ {
		fmt.Println(i)
		_, body, _ := gorequest.New().Get("http://qt.gtimg.cn/q=ff_sz00" + strconv.Itoa(i)[1:]).End()
		if len(body) == 0 {
			continue
		}
		x := strings.Split(string(body), "~")
		if len(x) == 1 {
			continue
		}
		name := utils.ConvertToString(x[12], "gbk", "utf-8")
		if name == "" {
			continue
		}
		dh := dal.Code{Code: x[0][len(x[0])-6:], Name: name}
		err := store.MysqlClient.GetDB().Save(&dh).Error

		fmt.Println(err, x[0][len(x[0])-6:], name)
	}

	for i := 600000; i <= 609999; i++ {
		_, body, _ := gorequest.New().Get("http://qt.gtimg.cn/q=ff_sh" + strconv.Itoa(i)).End()
		if len(body) == 0 {
			continue
		}
		x := strings.Split(string(body), "~")
		if len(x) == 1 {
			continue
		}
		name := utils.ConvertToString(x[12], "gbk", "utf-8")
		if name == "" {
			continue
		}
		dh := dal.Code{Code: x[0][len(x[0])-6:], Name: name}
		store.MysqlClient.GetDB().Save(&dh)
		fmt.Println(x[0][len(x[0])-6:], name)
	}

	for i := 688001; i <= 688050; i++ {
		// 科创板
		_, body, _ := gorequest.New().Get("http://qt.gtimg.cn/q=sh" + strconv.Itoa(i)).End()
		if len(body) == 0 {
			continue
		}
		x := strings.Split(string(body), "~")
		if len(x) == 1 {
			continue
		}
		name := utils.ConvertToString(x[1], "gbk", "utf-8")
		if name == "" {
			continue
		}
		dh := dal.Code{Code: x[2], Name: name}
		store.MysqlClient.GetDB().Save(&dh)
		fmt.Println(x[0][len(x[0])-6:], name)
	}

	for i := 300001; i <= 300900; i++ { // 创业板
		_, body, _ := gorequest.New().Get("http://qt.gtimg.cn/q=ff_sz" + strconv.Itoa(i)).End()
		if len(body) == 0 {
			continue
		}
		x := strings.Split(string(body), "~")
		if len(x) == 1 {
			continue
		}
		name := utils.ConvertToString(x[12], "gbk", "utf-8")
		if name == "" {
			continue
		}
		dh := dal.Code{Code: x[0][len(x[0])-6:], Name: name}
		store.MysqlClient.GetDB().Save(&dh)
		fmt.Println(x[0][len(x[0])-6:], name)
	}
}

// 根据历史生成周表(每个周五的收盘价) 执行一次就行
func (craw *Crawler) GenerateWeekHistory(code string) {
	dates := []string{"2019-11-29", "2019-11-22", "2019-11-15", "2019-11-08", "2019-11-01", "2019-10-25", "2019-10-18", "2019-10-11", "2019-09-30",
		"2019-09-27", "2019-09-20", "2019-09-12", "2019-09-06", "2019-08-30", "2019-08-23", "2019-08-16", "2019-08-09", "2019-08-02",
		"2019-07-26", "2019-07-19", "2019-07-12", "2019-07-05", "2019-06-28", "2019-06-21", "2019-06-14", "2019-06-06", "2019-05-31",
		"2019-05-24", "2019-05-17", "2019-05-10", "2019-04-30", "2019-04-26", "2019-04-19", "2019-04-12", "2019-04-04"}

	for _, i := range dates {
		var th dal.TicketHistory
		err := store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ? and date = ?", code, i).Find(&th).Error
		if err != nil {
			fmt.Println(err, i, code)
			continue
		}
		weekly := dal.TicketHistoryWeekly{Code: th.Code, Date: th.Date, Name: th.Name, Shou: th.Shou}
		store.MysqlClient.GetDB().Save(&weekly)
	}
}

// 获取单个股票的历史记录
func (craw *Crawler) GetSignalTicket(code, name string, proxy bool) error {
	var doc *goquery.Document
	for year := 2019; year <= 2019; year++ { // 1992 2019
		for ji := 2; ji <= 4; ji++ { // 1 4
			if !proxy {
				doc, _ = craw.NewDocument(fmt.Sprintf("http://quotes.money.163.com/trade/lsjysj_%s.html?year=%d&season=%d", code, year, ji))
			} else {
				doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://quotes.money.163.com/trade/lsjysj_%s.html?year=%d&season=%d", code, year, ji))
			}
			for i := 73; i >= 1; i-- { //for i := 73; i >= 1 ; i --
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
					store.MysqlClient.GetDB().Save(&dh)
					fmt.Println(code, name, date, kai, high, low, shou, zhangdiee, percent, tc, tm, zhenfu, huanshou)

				}
			}
		}
	}
	return nil
}

// 股票的所属概念信息 记录到 code 表中
// ex.g.经纬辉开 本月解禁,股权激励,小盘,苹果三星,证金汇金,特高压,新基建,电力物联网,高出口占比,苹果产业链,
func (craw *Crawler) GetAllTicketCodeConcept(code dal.Code, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://vip.stock.finance.sina.com.cn/corp/go.php/vCI_CorpOtherInfo/stockid/%s/menu_num/2.phtml", code.Code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://vip.stock.finance.sina.com.cn/corp/go.php/vCI_CorpOtherInfo/stockid/%s/menu_num/2.phtml", code.Code))
	}
	concept := ""
	for i := 3; i <= 20; i++ {
		x := doc.Find(fmt.Sprintf("#con02-0 > table:nth-child(2) > tbody > tr:nth-child(%d) > td:nth-child(1)", i)).Text()
		r := utils.ConvertToString(x, "gbk", "utf-8")
		if r == "" {
			continue
		}
		concept += r + ","
	}
	fmt.Println(code.ID, code.Name, concept)
	code.Concept = concept
	store.MysqlClient.GetDB().Save(&code)
}

// 拓展股票公司简介信息
func (craw *Crawler) GetAllTicketCodeInfo(code dal.Code, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://quotes.money.163.com/f10/gszl_%s.html#01f01", code.Code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://quotes.money.163.com/f10/gszl_%s.html#01f01", code.Code))
	}
	CompanyName := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(3) > td:nth-child(2)")).Text())
	OrganizationalForm := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(1) > td:nth-child(2)")).Text())
	Location := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(1) > td:nth-child(4)")).Text())
	Address := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(2) > td:nth-child(4)")).Text())
	NetAddress := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(9) > td:nth-child(2)")).Text())
	MajorBusinesses := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(11) > td:nth-child(2)")).Text())
	BusinessScope := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(12) > td:nth-child(2)")).Text())
	EstablishmentTime := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_r_01 > table > tbody > tr:nth-child(1) > td:nth-child(2)")).Text())
	ListingDate := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_r_01 > table > tbody > tr:nth-child(2) > td:nth-child(2)")).Text())
	fmt.Println(code.ID, code.Name, CompanyName, OrganizationalForm, Location, Address, NetAddress, MajorBusinesses, BusinessScope, EstablishmentTime, ListingDate)
	code.CompanyName = CompanyName
	code.OrganizationalForm = OrganizationalForm
	code.Location = Location
	code.Address = Address
	code.NetAddress = NetAddress
	code.MajorBusinesses = MajorBusinesses
	code.BusinessScope = BusinessScope
	code.EstablishmentTime = EstablishmentTime
	code.ListingDate = ListingDate
	store.MysqlClient.GetDB().Save(&code)
}

// 爬取所有股票的所属板块信息 记录到 code 表中
func (craw *Crawler) GetAllTicketCodeBelong(code dal.Code, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://vip.stock.finance.sina.com.cn/corp/go.php/vCI_CorpOtherInfo/stockid/%s/menu_num/2.phtml", code.Code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://vip.stock.finance.sina.com.cn/corp/go.php/vCI_CorpOtherInfo/stockid/%s/menu_num/2.phtml", code.Code))

	}
	x := doc.Find(fmt.Sprintf("#con02-0 > table:nth-child(1) > tbody > tr:nth-child(3) > td:nth-child(1)")).Text()
	r := utils.ConvertToString(x, "gbk", "utf-8")
	fmt.Println(code.Name, r)
	code.Belong = r
	store.MysqlClient.GetDB().Save(&code)
}

// 网易api获得十大流通股东（带新进变化趋势）
func (craw *Crawler) GetTopStockholder(code, namer string, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://quotes.money.163.com/f10/gdfx_%s.html#01d01", code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://quotes.money.163.com/f10/gdfx_%s.html#01d01", code))

	}

	for i := 1; i <= 10; i++ { // 2 4 5
		name := strings.TrimSpace(doc.Find(fmt.Sprintf("#ltdateTable > table > tbody > tr:nth-child(%d) > td.td_text", i)).Text())
		count := strings.TrimSpace(doc.Find(fmt.Sprintf("#ltdateTable > table > tbody > tr:nth-child(%d) > td:nth-child(3)", i)).Text())
		percent := strings.TrimSpace(doc.Find(fmt.Sprintf("#ltdateTable > table > tbody > tr:nth-child(%d) > td:nth-child(2)", i)).Text())

		change1 := strings.TrimSpace(doc.Find(fmt.Sprintf("#ltdateTable > table > tbody > tr:nth-child(%d) > td.cGreen", i)).Text())
		change2 := strings.TrimSpace(doc.Find(fmt.Sprintf("#ltdateTable > table > tbody > tr:nth-child(%d) > td.cRed", i)).Text())
		change3 := strings.TrimSpace(doc.Find(fmt.Sprintf("#ltdateTable > table > tbody > tr:nth-child(%d) > td:nth-child(4)", i)).Text())
		var change string
		if len(change1) > 0 {
			change = change1
		}
		if len(change2) > 0 {
			change = change2
		}
		if len(change3) > 0 {
			change = change3
		}
		fmt.Println(name, count, percent, change1, change2, change3)
		h := dal.Stockholder{Code: code, Name: namer, HolderName: name, Count: count, Percent: percent, Change: change}
		store.MysqlClient.GetDB().Save(&h)
	}
}

func (craw *Crawler) GetStockProfit(code string, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://vip.stock.finance.sina.com.cn/corp/go.php/vFD_ProfitStatement/stockid/%s/ctrl/part/displaytype/4.phtml", code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://vip.stock.finance.sina.com.cn/corp/go.php/vFD_ProfitStatement/stockid/%s/ctrl/part/displaytype/4.phtml", code))
	}
	for i := 2; i <= 6; i++ {
		date := strings.TrimSpace(doc.Find(fmt.Sprintf("#ProfitStatementNewTable0 > tbody > tr:nth-child(1) > td:nth-child(%d)", i)).Text())
		shouyi := strings.TrimSpace(doc.Find(fmt.Sprintf("#ProfitStatementNewTable0 > tbody > tr:nth-child(5) > td:nth-child(%d)", i)).Text())
		sy, err := strconv.ParseFloat(strings.Replace(shouyi, ",", "", -1), 64)
		if err != nil {
			sy = 0
		}
		chengben := strings.TrimSpace(doc.Find(fmt.Sprintf("#ProfitStatementNewTable0 > tbody > tr:nth-child(5) > td:nth-child(%d)", i)).Text())
		cb, err := strconv.ParseFloat(strings.Replace(chengben, ",", "", -1), 64)
		if err != nil {
			cb = 0
		}
		jinglirun := strings.TrimSpace(doc.Find(fmt.Sprintf("#ProfitStatementNewTable0 > tbody > tr:nth-child(23) > td:nth-child(%d)", i)).Text())
		hlr, err := strconv.ParseFloat(strings.Replace(jinglirun, ",", "", -1), 64)
		if err != nil {
			hlr = 0
		}
		profit := dal.StockProfit{Code: code, GrossTradingIncome: sy, TotalOperatingCost: cb, NetProfit: hlr, Date: date}
		store.MysqlClient.GetDB().Save(&profit)
		log.Println("公司利润表", profit)
	}

}

func (craw *Crawler) GetStockLiabilities(code string, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vFD_BalanceSheet/stockid/%s/ctrl/part/displaytype/4.phtml", code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vFD_BalanceSheet/stockid/%s/ctrl/part/displaytype/4.phtml", code))
	}
	for i := 2; i <= 6; i++ {
		date := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(1) > td:nth-child(%d)", i)).Text())
		// 流动资产
		liudongzichan := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(23) > td:nth-child(%d)", i)).Text())
		ldzc, err := strconv.ParseFloat(strings.Replace(liudongzichan, ",", "", -1), 64)
		if err != nil {
			ldzc = 0
		}
		// 非流动资产
		feiliudongzichan := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(47) > td:nth-child(%d)", i)).Text())
		fldzc, err := strconv.ParseFloat(strings.Replace(feiliudongzichan, ",", "", -1), 64)
		if err != nil {
			fldzc = 0
		}
		// 资产总计
		zichanziji := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(48) > td:nth-child(%d)", i)).Text())
		zctotal, err := strconv.ParseFloat(strings.Replace(zichanziji, ",", "", -1), 64)
		if err != nil {
			zctotal = 0
		}
		// 流动负债
		liudongfuzhai := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(68) > td:nth-child(%d)", i)).Text())
		ldfztotal, err := strconv.ParseFloat(strings.Replace(liudongfuzhai, ",", "", -1), 64)
		if err != nil {
			ldfztotal = 0
		}
		// 非流动负债
		feiliudongfuzhai := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(81) > td:nth-child(%d)", i)).Text())
		fldfztotal, err := strconv.ParseFloat(strings.Replace(feiliudongfuzhai, ",", "", -1), 64)
		if err != nil {
			fldfztotal = 0
		}
		// 负债总数
		fuzhaitotal := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(82) > td:nth-child(%d)", i)).Text())
		fztotal, err := strconv.ParseFloat(strings.Replace(fuzhaitotal, ",", "", -1), 64)
		if err != nil {
			fztotal = 0
		}
		profit := dal.StockLiabilities{Code: code, CurrentAssets: ldzc, NotCurrentAssets: fldzc, TotalAssets: zctotal, CurrentLiabilities: ldfztotal, NotCurrentLiabilities: fldfztotal, TotalLiabilities: fztotal, Date: date}
		store.MysqlClient.GetDB().Save(&profit)
		log.Println("资产负债表", code, profit)
	}
}

func (craw *Crawler) GetStockCashFlow(code string, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vFD_CashFlow/stockid/%s/ctrl/part/displaytype/4.phtml", code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vFD_CashFlow/stockid/%s/ctrl/part/displaytype/4.phtml", code))
	}
	for i := 2; i <= 6; i++ {
		date := strings.TrimSpace(doc.Find(fmt.Sprintf("#ProfitStatementNewTable0 > tbody > tr:nth-child(1) > td:nth-child(%d)", i)).Text())

		// 经营活动产生的现金流量净额
		jingyinghuodong := strings.TrimSpace(doc.Find(fmt.Sprintf("#ProfitStatementNewTable0 > tbody > tr:nth-child(13) > td:nth-child(%d)", i)).Text())
		jy, err := strconv.ParseFloat(strings.Replace(jingyinghuodong, ",", "", -1), 64)
		if err != nil {
			jy = 0
		}
		// 投资活动产生的现金流量净额
		touzihuodong := strings.TrimSpace(doc.Find(fmt.Sprintf("#ProfitStatementNewTable0 > tbody > tr:nth-child(26) > td:nth-child(%d)", i)).Text())
		tz, err := strconv.ParseFloat(strings.Replace(touzihuodong, ",", "", -1), 64)
		if err != nil {
			tz = 0
		}
		// 募集活动产生的现金流量金额
		mujihuodong := strings.TrimSpace(doc.Find(fmt.Sprintf("#ProfitStatementNewTable0 > tbody > tr:nth-child(39) > td:nth-child(%d)", i)).Text())
		mj, err := strconv.ParseFloat(strings.Replace(mujihuodong, ",", "", -1), 64)
		if err != nil {
			mj = 0
		}
		// 期末现金及现金等价物余额
		qimo := strings.TrimSpace(doc.Find(fmt.Sprintf("#ProfitStatementNewTable0 > tbody > tr:nth-child(43) > td:nth-child(%d)", i)).Text())
		qm, err := strconv.ParseFloat(strings.Replace(qimo, ",", "", -1), 64)
		if err != nil {
			qm = 0
		}
		cash_flow := dal.StockCashFlow{Code: code, ManageCashFlow: jy, InvestCashFlow: tz, FundraisingCashFlow: mj, CashRemain: qm, Date: date}
		store.MysqlClient.GetDB().Save(&cash_flow)
		log.Println("现金流量表", code, cash_flow)
	}
}

// 周线表计算涨跌幅
func (craw *Crawler) CalcPercentTicketWeekly(code string) {
	var th []dal.TicketHistoryWeekly
	store.MysqlClient.GetDB().Model(&dal.TicketHistoryWeekly{}).Where("code = ?", code).Order("date asc").Find(&th)
	var tmp dal.TicketHistoryWeekly

	for j, i := range th {
		if j == 0 {
			i.Percent = 0
			store.MysqlClient.GetDB().Save(&i)
			tmp = i
			continue
		}
		fmt.Println(i.Shou, tmp.Shou, i.Name, fmt.Sprintf("%.2f", (i.Shou-tmp.Shou)*100/tmp.Shou))
		i.Percent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", (i.Shou-tmp.Shou)*100/tmp.Shou), 64)
		store.MysqlClient.GetDB().Save(&i)
		tmp = i
	}

}

//// 计算涨跌幅， 放在history percent 点 (只需执行一次就把历史数据跑完)
//func (craw *Crawler) CalcPercentForAllTicket(code string) {
//	var th []dal.TicketHistory
//	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ?", code).Order("date asc").Find(&th)
//	var tmp dal.TicketHistory
//
//	for j, i := range th {
//		if j == 0 {
//			i.Percent = 0
//			store.MysqlClient.GetDB().Save(&i)
//			tmp = i
//			continue
//		}
//		fmt.Println(i.Shou, tmp.Shou, i.Name, fmt.Sprintf("%.2f", (i.Shou-tmp.Shou)*100/tmp.Shou))
//		i.Percent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", (i.Shou-tmp.Shou)*100/tmp.Shou), 64)
//		store.MysqlClient.GetDB().Save(&i)
//		tmp = i
//	}
//}
