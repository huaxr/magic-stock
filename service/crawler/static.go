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
	"time"

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

// 获取单个股票的历史记录
// 我要得出 60 季度线的所有数据 则需要从 2004 年开始拿所有数据
func (craw *Crawler) GetSignalTicket(code dal.Code, proxy bool) error {
	var doc *goquery.Document
	for year := 1991; year <= 2020; year++ { // 1992 2019
		for ji := 1; ji <= 4; ji++ { // 1 4
			if !proxy {
				doc, _ = craw.NewDocument(fmt.Sprintf("http://quotes.money.163.com/trade/lsjysj_%s.html?year=%d&season=%d", code.Code, year, ji))
			} else {
				doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://quotes.money.163.com/trade/lsjysj_%s.html?year=%d&season=%d", code.Code, year, ji))
			}
			for i := 80; i >= 1; i-- { //for i := 73; i >= 1 ; i --
				// 一季度按照72个交易日来算  // 因为第一栏是中文标题
				text := ""

				x, err := dealPanic(doc, i, "body > div.area > div.inner_box > table > tbody > tr:nth-child(%d) > td:nth-child(1)")
				if err != nil {
					log.Println("出现错误了, 正在重试")
					time.Sleep(5 * time.Second)
					i += 1
					continue
				}
				// 获取日期不要用
				//x := doc.Find(fmt.Sprintf("body > div.area > div.inner_box > table > tbody > tr:nth-child(%d) > td:nth-child(1)", i)).Text()
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

					if code.ID <= 1500 {
						dh := dal.HistoryALL1{Code: code.Code, Name: code.Name, Kai: kai, High: high, Low: low, Shou: shou, TotalCount: tc, TotalMoney: tm, Date: date,
							Percent: percent, Change: zhangdiee, Amplitude: zhenfu, TurnoverRate: huanshou} //, Percent:p}
						store.MysqlClient.GetDB().Save(&dh)
					} else {
						dh := dal.HistoryALL2{Code: code.Code, Name: code.Name, Kai: kai, High: high, Low: low, Shou: shou, TotalCount: tc, TotalMoney: tm, Date: date,
							Percent: percent, Change: zhangdiee, Amplitude: zhenfu, TurnoverRate: huanshou} //, Percent:p}
						store.MysqlClient.GetDB().Save(&dh)
					}
				}
			}
		}
	}
	log.Println("爬取完毕", code.Name)
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
	Location := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(1) > td:nth-child(4)")).Text())
	Address := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(2) > td:nth-child(4)")).Text())
	NetAddress := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(9) > td:nth-child(2)")).Text())
	MajorBusinesses := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(11) > td:nth-child(2)")).Text())
	BusinessScope := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_l_01 > table > tbody > tr:nth-child(12) > td:nth-child(2)")).Text())
	EstablishmentTime := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_r_01 > table > tbody > tr:nth-child(1) > td:nth-child(2)")).Text())
	ListingDate := strings.TrimSpace(doc.Find(fmt.Sprintf("body > div.area > div.col_r_01 > table > tbody > tr:nth-child(2) > td:nth-child(2)")).Text())
	fmt.Println(code.ID, code.Name, CompanyName, Location, Address, NetAddress, MajorBusinesses, BusinessScope, EstablishmentTime, ListingDate)
	code.CompanyName = CompanyName
	code.Location = Location
	code.Address = Address
	code.NetAddress = NetAddress
	code.MajorBusinesses = MajorBusinesses
	code.BusinessScope = BusinessScope
	code.EstablishmentTime = EstablishmentTime
	code.ListingDate = ListingDate
	store.MysqlClient.GetDB().Save(&code)
}

// 拓展股票公司简介信息(用新浪api更牛)
func (craw *Crawler) GetAllTicketCodeInfo2(code dal.Code, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://vip.stock.finance.sina.com.cn/corp/go.php/vCI_CorpInfo/stockid/%s.phtml", code.Code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://vip.stock.finance.sina.com.cn/corp/go.php/vCI_CorpInfo/stockid/%s.phtml", code.Code))
	}
	OrganizationalForm := strings.TrimSpace(doc.Find(fmt.Sprintf("#comInfo1 > tbody > tr:nth-child(6) > td:nth-child(4)")).Text())
	OrganizationalForm = utils.ConvertToString(OrganizationalForm, "gbk", "utf-8")
	InstitutionalType := strings.TrimSpace(doc.Find(fmt.Sprintf("#comInfo1 > tbody > tr:nth-child(6) > td:nth-child(2)")).Text())
	HistoryNames := strings.TrimSpace(doc.Find(fmt.Sprintf("#comInfo1 > tbody > tr:nth-child(17) > td.ccl")).Text())
	InstitutionalType = utils.ConvertToString(InstitutionalType, "gbk", "utf-8")
	HistoryNames = utils.ConvertToString(HistoryNames, "gbk", "utf-8")
	fmt.Println(code.ID, code.Name, OrganizationalForm, HistoryNames, InstitutionalType)
	code.OrganizationalForm = OrganizationalForm
	code.InstitutionalType = InstitutionalType
	code.HistoryNames = HistoryNames
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

// 每股指标 营业能力 运营能力
func (craw *Crawler) GetStockPerTicket(code string, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vFD_FinancialGuideLine/stockid/%s/displaytype/4.phtml", code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vFD_FinancialGuideLine/stockid/%s/displaytype/4.phtml", code))
	}
	for i := 2; i <= 2; i++ {
		date := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(1) > td:nth-child(%d)", i)).Text())

		tanbo := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(3) > td:nth-child(%d)", i)).Text())
		tanbo1, err := strconv.ParseFloat(strings.Replace(tanbo, "--", "", -1), 64)
		if err != nil {
			tanbo1 = 0
		}
		jiaquanshouyi := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(4) > td:nth-child(%d)", i)).Text())
		jiaquanshouyi1, err := strconv.ParseFloat(strings.Replace(jiaquanshouyi, "--", "", -1), 64)
		if err != nil {
			jiaquanshouyi1 = 0
		}
		shouyi := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(5) > td:nth-child(%d)", i)).Text())
		shouyi1, err := strconv.ParseFloat(strings.Replace(shouyi, "--", "", -1), 64)
		if err != nil {
			shouyi1 = 0
		}
		jizichan_front := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(7) > td:nth-child(%d)", i)).Text())
		jizichan_front1, err := strconv.ParseFloat(strings.Replace(jizichan_front, "--", "", -1), 64)
		if err != nil {
			jizichan_front1 = 0
		}
		jizichan_after := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(8) > td:nth-child(%d)", i)).Text())
		jizichan_after1, err := strconv.ParseFloat(strings.Replace(jizichan_after, "--", "", -1), 64)
		if err != nil {
			jizichan_after1 = 0
		}
		jingyingxianjinliu := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(9) > td:nth-child(%d)", i)).Text())
		jingyingxianjinliu1, err := strconv.ParseFloat(strings.Replace(jingyingxianjinliu, "--", "", -1), 64)
		if err != nil {
			jingyingxianjinliu1 = 0
		}
		zibengongjijin := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(10) > td:nth-child(%d)", i)).Text())
		zibengongjijin1, err := strconv.ParseFloat(strings.Replace(zibengongjijin, "--", "", -1), 64)
		if err != nil {
			zibengongjijin1 = 0
		}
		weifenpeilirun := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(11) > td:nth-child(%d)", i)).Text())
		weifenpeilirun1, err := strconv.ParseFloat(strings.Replace(weifenpeilirun, "--", "", -1), 64)
		if err != nil {
			weifenpeilirun1 = 0
		}

		// 盈利能力
		YlZongzichanlirunlv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(14) > td:nth-child(%d)", i)).Text())
		YlZongzichanlirunlv1, err := strconv.ParseFloat(strings.Replace(YlZongzichanlirunlv, "--", "", -1), 64)
		if err != nil {
			YlZongzichanlirunlv1 = 0
		}
		YlZhuyingyewulirunlv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(15) > td:nth-child(%d)", i)).Text())
		YlZhuyingyewulirunlv1, err := strconv.ParseFloat(strings.Replace(YlZhuyingyewulirunlv, "--", "", -1), 64)
		if err != nil {
			YlZhuyingyewulirunlv1 = 0
		}

		YlZongzichanjinglirunlv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(16) > td:nth-child(%d)", i)).Text())
		YlZongzichanjinglirunlv1, err := strconv.ParseFloat(strings.Replace(YlZongzichanjinglirunlv, "--", "", -1), 64)
		if err != nil {
			YlZongzichanjinglirunlv1 = 0
		}

		YlYingyelirunlv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(18) > td:nth-child(%d)", i)).Text())
		YlYingyelirunlv1, err := strconv.ParseFloat(strings.Replace(YlYingyelirunlv, "--", "", -1), 64)
		if err != nil {
			YlYingyelirunlv1 = 0
		}

		YlXiaoshoujinglilv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(20) > td:nth-child(%d)", i)).Text())
		YlXiaoshoujinglilv1, err := strconv.ParseFloat(strings.Replace(YlXiaoshoujinglilv, "--", "", -1), 64)
		if err != nil {
			YlXiaoshoujinglilv1 = 0
		}

		YlGubenbaochoulv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(21) > td:nth-child(%d)", i)).Text())
		YlGubenbaochoulv1, err := strconv.ParseFloat(strings.Replace(YlGubenbaochoulv, "--", "", -1), 64)
		if err != nil {
			YlGubenbaochoulv1 = 0
		}

		YlJingzichanbaochoulv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(22) > td:nth-child(%d)", i)).Text())
		YlJingzichanbaochoulv1, err := strconv.ParseFloat(strings.Replace(YlJingzichanbaochoulv, "--", "", -1), 64)
		if err != nil {
			YlJingzichanbaochoulv1 = 0
		}

		YlZichanbaochoulv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(23) > td:nth-child(%d)", i)).Text())
		YlZichanbaochoulv1, err := strconv.ParseFloat(strings.Replace(YlZichanbaochoulv, "--", "", -1), 64)
		if err != nil {
			YlZichanbaochoulv1 = 0
		}
		CzZhuyingyewushouruzengzhanglv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(35) > td:nth-child(%d)", i)).Text())
		CzZhuyingyewushouruzengzhanglv1, err := strconv.ParseFloat(strings.Replace(CzZhuyingyewushouruzengzhanglv, "--", "", -1), 64)
		if err != nil {
			CzZhuyingyewushouruzengzhanglv1 = 0
		}

		Czjinglirunzengzhanglv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(36) > td:nth-child(%d)", i)).Text())
		Czjinglirunzengzhanglv1, err := strconv.ParseFloat(strings.Replace(Czjinglirunzengzhanglv, "--", "", -1), 64)
		if err != nil {
			Czjinglirunzengzhanglv1 = 0
		}

		CzJingzichanzengzhanglv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(37) > td:nth-child(%d)", i)).Text())
		CzJingzichanzengzhanglv1, err := strconv.ParseFloat(strings.Replace(CzJingzichanzengzhanglv, "--", "", -1), 64)
		if err != nil {
			CzJingzichanzengzhanglv1 = 0
		}

		CzZongzichanzengzhanglv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(38) > td:nth-child(%d)", i)).Text())
		CzZongzichanzengzhanglv1, err := strconv.ParseFloat(strings.Replace(CzZongzichanzengzhanglv, "--", "", -1), 64)
		if err != nil {
			CzZongzichanzengzhanglv1 = 0
		}

		YyYingshouzhangkuanzhouzhuanlv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(40) > td:nth-child(%d)", i)).Text())
		YyYingshouzhangkuanzhouzhuanlv1, err := strconv.ParseFloat(strings.Replace(YyYingshouzhangkuanzhouzhuanlv, "--", "", -1), 64)
		if err != nil {
			YyYingshouzhangkuanzhouzhuanlv1 = 0
		}
		YyCunhuozhouzhuanglv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(43) > td:nth-child(%d)", i)).Text())
		YyCunhuozhouzhuanglv1, err := strconv.ParseFloat(strings.Replace(YyCunhuozhouzhuanglv, "--", "", -1), 64)
		if err != nil {
			YyCunhuozhouzhuanglv1 = 0
		}
		YyLiudongzichanzhouzhuanglv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(47) > td:nth-child(%d)", i)).Text())
		YyLiudongzichanzhouzhuanglv1, err := strconv.ParseFloat(strings.Replace(YyLiudongzichanzhouzhuanglv, "--", "", -1), 64)
		if err != nil {
			YyLiudongzichanzhouzhuanglv1 = 0
		}
		YyZongzichanzhouzhuanglv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(45) > td:nth-child(%d)", i)).Text())
		YyZongzichanzhouzhuanglv1, err := strconv.ParseFloat(strings.Replace(YyZongzichanzhouzhuanglv, "--", "", -1), 64)
		if err != nil {
			YyZongzichanzhouzhuanglv1 = 0
		}

		YyGudongquanyizhouzhuanglv := strings.TrimSpace(doc.Find(fmt.Sprintf("#BalanceSheetNewTable0 > tbody > tr:nth-child(49) > td:nth-child(%d)", i)).Text())
		YyGudongquanyizhouzhuanglv1, err := strconv.ParseFloat(strings.Replace(YyGudongquanyizhouzhuanglv, "--", "", -1), 64)
		if err != nil {
			YyGudongquanyizhouzhuanglv1 = 0
		}

		per_ticket := dal.StockPerTicket{Code: code, Tanboshouyi: tanbo1, Jiaquanshouyi: jiaquanshouyi1, Shouyiafter: shouyi1, Jinzichanfront: jizichan_front1,
			Jinzichanafter: jizichan_after1, Jingyingxianjinliu: jingyingxianjinliu1, Gubengongjijin: zibengongjijin1, Weifenpeilirun: weifenpeilirun1, YlZhuyingyewulirunlv: YlZhuyingyewulirunlv1,
			YlZongzichanlirunlv: YlZongzichanlirunlv1, YlZongzichanjinglirunlv: YlZongzichanjinglirunlv1, YlYingyelirunlv: YlYingyelirunlv1, YlXiaoshoujinglilv: YlXiaoshoujinglilv1,
			YlGubenbaochoulv: YlGubenbaochoulv1, YlJingzichanbaochoulv: YlJingzichanbaochoulv1, YlZichanbaochoulv: YlZichanbaochoulv1, CzZhuyingyewushouruzengzhanglv: CzZhuyingyewushouruzengzhanglv1,
			CzJinglirunzengzhanglv: Czjinglirunzengzhanglv1, CzJingzichanzengzhanglv: CzJingzichanzengzhanglv1, CzZongzichanzengzhanglv: CzZongzichanzengzhanglv1, YyYingshouzhangkuanzhouzhuanlv: YyYingshouzhangkuanzhouzhuanlv1,
			YyCunhuozhouzhuanglv: YyCunhuozhouzhuanglv1, YyLiudongzichanzhouzhuanglv: YyLiudongzichanzhouzhuanglv1, YyZongzichanzhouzhuanglv: YyZongzichanzhouzhuanglv1, YyGudongquanyizhouzhuanglv: YyGudongquanyizhouzhuanglv1,
			Date: date}
		store.MysqlClient.GetDB().Save(&per_ticket)
	}
}

// 分红 配股 增发 （这些属于股票的历史数据 可以每三个月爬一次？ 2019-12-31 号爬取）
func (craw *Crawler) GetProfitSharingAndStockOwnership(code string, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vISSUE_ShareBonus/stockid/%s.phtml", code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vISSUE_ShareBonus/stockid/%s.phtml", code))
	}
	for i := 1; i <= 100; i++ {
		a := doc.Find(fmt.Sprintf("#sharebonus_1 > tbody > tr:nth-child(%d) > td:nth-child(1)", i)).Text()
		if len(a) == 0 {
			break
		}
		b := doc.Find(fmt.Sprintf("#sharebonus_1 > tbody > tr:nth-child(%d) > td:nth-child(2)", i)).Text()
		c := doc.Find(fmt.Sprintf("#sharebonus_1 > tbody > tr:nth-child(%d) > td:nth-child(3)", i)).Text()
		d := doc.Find(fmt.Sprintf("#sharebonus_1 > tbody > tr:nth-child(%d) > td:nth-child(4)", i)).Text()
		e := doc.Find(fmt.Sprintf("#sharebonus_1 > tbody > tr:nth-child(%d) > td:nth-child(5)", i)).Text()
		e = utils.ConvertToString(e, "gbk", "utf-8")
		log.Println(code, a, b, c, d, e)

		x := dal.StockFengHong{Code: code, Date: a, SongGu: b, ZhuangZeng: c, PaiXi: d, Process: e}
		store.MysqlClient.GetDB().Save(&x)
	}
	for i := 1; i <= 100; i++ {
		f := doc.Find(fmt.Sprintf("#sharebonus_2 > tbody > tr:nth-child(%d) > td:nth-child(1)", i)).Text()
		if len(f) == 0 {
			break
		}
		g := doc.Find(fmt.Sprintf("#sharebonus_2 > tbody > tr:nth-child(%d) > td:nth-child(2)", i)).Text()
		h := doc.Find(fmt.Sprintf("#sharebonus_2 > tbody > tr:nth-child(%d) > td:nth-child(3)", i)).Text()
		j := doc.Find(fmt.Sprintf("#sharebonus_2 > tbody > tr:nth-child(%d) > td:nth-child(4)", i)).Text()
		log.Println(code, f, g, h, j)
		x := dal.StockPeiGu{Code: code, Date: f, Count: g, Price: h, Number: j}
		store.MysqlClient.GetDB().Save(&x)
	}

}

// 获取增发数据
func (craw *Crawler) GetZengFa(code string, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vISSUE_AddStock/stockid/%s.phtml", code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vISSUE_AddStock/stockid/%s.phtml", code))
	}
	for i := 0; i <= 100; i++ {
		a := doc.Find(fmt.Sprintf("#addStock%d > thead > tr > th", i)).Text()
		a = utils.ConvertToString(a, "gbk", "utf-8")
		if len(a) == 0 {
			break
		}
		a = a[len(a)-10 : len(a)]
		b := doc.Find(fmt.Sprintf("#addStock%d > tbody > tr:nth-child(1) > td:nth-child(2)", i)).Text()
		b = utils.ConvertToString(b, "gbk", "utf-8")
		c := doc.Find(fmt.Sprintf("#addStock%d > tbody > tr:nth-child(2) > td:nth-child(2)", i)).Text()
		c = utils.ConvertToString(c, "gbk", "utf-8")
		d := doc.Find(fmt.Sprintf("#addStock%d > tbody > tr:nth-child(3) > td:nth-child(2)", i)).Text()
		d = utils.ConvertToString(d, "gbk", "utf-8")
		e := doc.Find(fmt.Sprintf("#addStock%d > tbody > tr:nth-child(4) > td:nth-child(2)", i)).Text()
		e = utils.ConvertToString(e, "gbk", "utf-8")
		f := doc.Find(fmt.Sprintf("#addStock%d > tbody > tr:nth-child(5)> td:nth-child(2)", i)).Text()
		f = utils.ConvertToString(f, "gbk", "utf-8")
		log.Println(code, a, b, c, d, e)

		x := dal.StockZengFa{Code: code, Date: a, Way: b, Price: c, AllPrice: d, CostPrice: e, AllCount: f}
		store.MysqlClient.GetDB().Save(&x)
	}

}

// 子公司记录， 每半年出一次报表用的
func (craw *Crawler) GetSubCompany(code string, proxy bool) {
	var doc *goquery.Document
	if !proxy {
		doc, _ = craw.NewDocument(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vCO_HoldingCompany/stockid/%s.phtml", code))
	} else {
		doc, _ = craw.NewDocumentWithProxy(fmt.Sprintf("http://money.finance.sina.com.cn/corp/go.php/vCO_HoldingCompany/stockid/%s.phtml", code))
	}
	//#holdingcompany > tbody:nth-child(2) > tr:nth-child(15) > td:nth-child(4)
	//#holdingcompany > tbody:nth-child(2) > tr:nth-child(1) > td.head
	for i := 1; i <= 10000; i++ {
		a := doc.Find(fmt.Sprintf("#holdingcompany > tbody:nth-child(2) > tr:nth-child(%d) > td:nth-child(1)", i)).Text()
		a = utils.ConvertToString(a, "gbk", "utf-8")
		if len(a) == 0 {
			break
		}
		b := doc.Find(fmt.Sprintf("#holdingcompany > tbody:nth-child(2) > tr:nth-child(%d) > td:nth-child(2)", i)).Text()
		b = utils.ConvertToString(b, "gbk", "utf-8")
		c := doc.Find(fmt.Sprintf("#holdingcompany > tbody:nth-child(2) > tr:nth-child(%d) > td:nth-child(3)", i)).Text()
		c = utils.ConvertToString(c, "gbk", "utf-8")
		d := doc.Find(fmt.Sprintf("#holdingcompany > tbody:nth-child(2) > tr:nth-child(%d) > td:nth-child(4)", i)).Text()
		d = utils.ConvertToString(d, "gbk", "utf-8")
		log.Println(code, a, b, c, d)
		x := dal.StockSubCompany{Code: code, Name: a, Relation: b, Percent: c, Type: d}
		store.MysqlClient.GetDB().Save(&x)
	}
}

type Date struct {
	Date string
}

type Res struct {
	Number float64
}

func (craw *Crawler) GetWeekDays(code dal.Code, date1, date2 string) {
	var model interface{}
	if code.ID <= 1500 {
		model = &dal.HistoryALL1{}
	} else {
		model = &dal.HistoryALL2{}
	}
	var dates []Date
	var res []string
	store.MysqlClient.GetDB().Model(model).Select("date").Where("code = ?", code.Code).Where("date >= ? and date < ?", date1, date2).Scan(&dates)
	for _, i := range dates {
		res = append(res, i.Date)
	}
	if len(res) == 0 {
		return
	}
	x, xx := utils.GetWeekPair(res)
	var last Res
	for i := 0; i <= len(x)/5-1; i++ {
		tmp := strings.Replace(x[5*i:5*(i+1)], "0", "", -1)
		if len(tmp) == 0 {
			continue
		}
		res := xx[0:len(tmp)]
		xx = xx[len(tmp):len(xx)]

		var high, low, kai, shou, liang, turnover_rate Res
		var p, a float64 // percent 振幅

		store.MysqlClient.GetDB().Model(model).Select("max(high) as number").Where("code = ?", code.Code).Where("date >= ? and date <= ?", res[0][1], res[len(res)-1][1]).Scan(&high)
		store.MysqlClient.GetDB().Model(model).Select("min(low) as number").Where("code = ?", code.Code).Where("date >= ? and date <= ?", res[0][1], res[len(res)-1][1]).Scan(&low)
		store.MysqlClient.GetDB().Model(model).Select("kai as number").Where("code = ?", code.Code).Where("date = ?", res[0][1]).Scan(&kai)
		store.MysqlClient.GetDB().Model(model).Select("shou as number").Where("code = ?", code.Code).Where("date = ?", res[len(res)-1][1]).Scan(&shou)
		store.MysqlClient.GetDB().Model(model).Select("sum(total_count) as number").Where("code = ?", code.Code).Where("date >= ? and date <= ?", res[0][1], res[len(res)-1][1]).Scan(&liang)
		store.MysqlClient.GetDB().Model(model).Select("sum(turnover_rate) as number").Where("code = ?", code.Code).Where("date >= ? and date <= ?", res[0][1], res[len(res)-1][1]).Scan(&turnover_rate)
		if i == 0 {
			p = 0
			a = 0
			last = shou
		} else {
			percent := (shou.Number - last.Number) / last.Number
			amplitude := (high.Number - low.Number) / last.Number
			p, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", percent*100), 64)
			a, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", amplitude*100), 64)
		}
		change, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", shou.Number-last.Number), 64)
		turnover, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", turnover_rate.Number), 64)
		log.Println(high.Number, low.Number, kai.Number, shou.Number, liang.Number, turnover, p, a, change, res[len(res)-1][1], code.Code)
		last = shou
		h := dal.TicketHistoryWeeklyALL{Code: code.Code, Name: code.Name, Date: res[len(res)-1][1].(string), Kai: kai.Number, Shou: shou.Number, High: high.Number, Low: low.Number,
			TotalCount: liang.Number, Percent: p, Change: change, Amplitude: a, TurnoverRate: turnover}
		store.MysqlClient.GetDB().Save(&h)
	}
}

func (craw *Crawler) GetMonthDays(code dal.Code) {
	var model interface{}
	if code.ID <= 1500 {
		model = &dal.HistoryALL1{}
	} else {
		model = &dal.HistoryALL2{}
	}
	var year = []string{"1991", "1992", "1993", "1994", "1995", "1996", "1997", "1998", "1999",
		"2000", "2001", "2002", "2003", "2004", "2005", "2006", "2007", "2008", "2009", "2010",
		"2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020"}
	var month = []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	var last Res
	for j := 0; j <= 29; j++ {
		var date1, date2 string
		for i := 0; i <= 11; i++ {
			if i == 11 {
				if j+1 > 29 {
					continue
				}
				date1 = year[j] + "-" + month[i]
				date2 = year[j+1] + "-" + month[0]
			} else {
				date1 = year[j] + "-" + month[i]
				date2 = year[j] + "-" + month[i+1]
			}
			if date1 == "2020-01" {
				break
			}
			var c int
			store.MysqlClient.GetDB().Model(model).Select("max(date) as date").Where("code = ?", code.Code).Where("date >= ? and date < ?", date1, date2).Count(&c)
			if c == 0 {
				continue
			}
			var da_min, da_max Date
			store.MysqlClient.GetDB().Model(model).Select("max(date) as date").Where("code = ?", code.Code).Where("date >= ? and date < ?", date1, date2).Scan(&da_max)

			store.MysqlClient.GetDB().Model(model).Select("min(date) as date").Where("code = ?", code.Code).Where("date >= ? and date < ?", date1, date2).Scan(&da_min)

			var high, low, kai, shou, liang, turnover_rate Res
			var p, a float64 // percent 振幅

			store.MysqlClient.GetDB().Model(model).Select("max(high) as number").Where("code = ?", code.Code).Where("date >= ? and date < ?", da_min.Date, da_max.Date).Scan(&high)
			store.MysqlClient.GetDB().Model(model).Select("min(low) as number").Where("code = ?", code.Code).Where("date >= ? and date < ?", da_min.Date, da_max.Date).Scan(&low)
			store.MysqlClient.GetDB().Model(model).Select("kai as number").Where("code = ?", code.Code).Where("date = ?", da_min.Date).Scan(&kai)
			store.MysqlClient.GetDB().Model(model).Select("shou as number").Where("code = ?", code.Code).Where("date = ?", da_max.Date).Scan(&shou)
			store.MysqlClient.GetDB().Model(model).Select("sum(total_count) as number").Where("code = ?", code.Code).Where("date >= ? and date < ?", da_min.Date, da_max.Date).Scan(&liang)
			store.MysqlClient.GetDB().Model(model).Select("sum(turnover_rate) as number").Where("code = ?", code.Code).Where("date >= ? and date < ?", da_min.Date, da_max.Date).Scan(&turnover_rate)

			if i == 0 && j == 0 {
				p = 0
				a = 0
				last = shou
			} else {
				percent := (shou.Number - last.Number) / last.Number
				amplitude := (high.Number - low.Number) / last.Number
				p, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", percent*100), 64)
				a, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", amplitude*100), 64)
			}
			change, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", shou.Number-last.Number), 64)
			turnover, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", turnover_rate.Number), 64)
			log.Println(high.Number, low.Number, kai.Number, shou.Number, liang.Number, turnover, p, a, change, da_max.Date, code.ID, code.Code, code.Name)
			last = shou
			h := dal.TicketHistoryMonthAll{Code: code.Code, Name: code.Name, Date: da_max.Date, Kai: kai.Number, Shou: shou.Number, High: high.Number, Low: low.Number,
				TotalCount: liang.Number, Percent: p, Change: change, Amplitude: a, TurnoverRate: turnover}
			store.MysqlClient.GetDB().Save(&h)
		}
	}
}

// 只获取今日的周线
func (craw *Crawler) GetWeekDay(code dal.Code, date1, date2 string, last_day_delete string) {
	if last_day_delete != "" {
		store.MysqlClient.GetDB().Exec("delete from magic_stock_history_week where date = ?", last_day_delete)
	}
	var model = &dal.TicketHistory{}
	var dates []Date
	var res []string
	store.MysqlClient.GetDB().Model(model).Select("date").Where("code = ?", code.Code).Where("date >= ? and date <= ?", date1, date2).Scan(&dates)
	for _, i := range dates {
		res = append(res, i.Date)
	}
	if len(res) == 0 {
		return
	}
	x, xx := utils.GetWeekPair(res)
	log.Println(x, xx)
	var last Res
	for i := 0; i <= len(x)/5-1; i++ {
		tmp := strings.Replace(x[5*i:5*(i+1)], "0", "", -1)
		if len(tmp) == 0 {
			continue
		}
		res := xx[0:len(tmp)]
		xx = xx[len(tmp):len(xx)]

		var high, low, kai, shou, liang, turnover_rate Res
		var p, a float64 // percent 振幅

		store.MysqlClient.GetDB().Model(model).Select("max(high) as number").Where("code = ?", code.Code).Where("date >= ? and date <= ?", res[0][1], res[len(res)-1][1]).Scan(&high)
		store.MysqlClient.GetDB().Model(model).Select("min(low) as number").Where("code = ?", code.Code).Where("date >= ? and date <= ?", res[0][1], res[len(res)-1][1]).Scan(&low)
		store.MysqlClient.GetDB().Model(model).Select("kai as number").Where("code = ?", code.Code).Where("date = ?", res[0][1]).Scan(&kai)
		store.MysqlClient.GetDB().Model(model).Select("shou as number").Where("code = ?", code.Code).Where("date = ?", res[len(res)-1][1]).Scan(&shou)
		store.MysqlClient.GetDB().Model(model).Select("sum(total_count) as number").Where("code = ?", code.Code).Where("date >= ? and date <= ?", res[0][1], res[len(res)-1][1]).Scan(&liang)
		store.MysqlClient.GetDB().Model(model).Select("sum(turnover_rate) as number").Where("code = ?", code.Code).Where("date >= ? and date <= ?", res[0][1], res[len(res)-1][1]).Scan(&turnover_rate)
		if i == 0 {
			p = 0
			a = 0
			last = shou
		} else {
			percent := (shou.Number - last.Number) / last.Number
			amplitude := (high.Number - low.Number) / last.Number
			p, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", percent*100), 64)
			a, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", amplitude*100), 64)
		}
		change, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", shou.Number-last.Number), 64)
		turnover, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", turnover_rate.Number), 64)
		log.Println(high.Number, low.Number, kai.Number, shou.Number, liang.Number, turnover, p, a, change, res[len(res)-1][1], code.Code)
		last = shou
		h := dal.TicketHistoryWeekly{Code: code.Code, Name: code.Name, Date: res[len(res)-1][1].(string), Kai: kai.Number, Shou: shou.Number, High: high.Number, Low: low.Number,
			TotalCount: liang.Number, Percent: p, Change: change, Amplitude: a, TurnoverRate: turnover}
		store.MysqlClient.GetDB().Save(&h)
	}
}
func (craw *Crawler) GetMonthDay(code dal.Code, today string, last_month_day string, last_day_delete string) {
	if last_day_delete != "" {
		store.MysqlClient.GetDB().Exec("delete from magic_stock_history_month where date = ?", last_day_delete)

	}
	model := &dal.TicketHistory{}
	var c int
	store.MysqlClient.GetDB().Model(model).Select("max(date) as date").Where("code = ?", code.Code).Where("date > ? and date <= ?", last_month_day, today).Count(&c)
	if c == 0 {
		return
	}
	var da_min, da_max Date
	store.MysqlClient.GetDB().Model(model).Select("max(date) as date").Where("code = ?", code.Code).Where("date > ? and date <= ?", last_month_day, today).Scan(&da_max)
	store.MysqlClient.GetDB().Model(model).Select("min(date) as date").Where("code = ?", code.Code).Where("date > ? and date <= ?", last_month_day, today).Scan(&da_min)

	var high, low, kai, shou, liang, turnover_rate Res
	var p, a float64 // percent 振幅
	store.MysqlClient.GetDB().Model(model).Select("max(high) as number").Where("code = ?", code.Code).Where("date > ? and date <= ?", da_min.Date, da_max.Date).Scan(&high)
	store.MysqlClient.GetDB().Model(model).Select("min(low) as number").Where("code = ?", code.Code).Where("date > ? and date <= ?", da_min.Date, da_max.Date).Scan(&low)
	store.MysqlClient.GetDB().Model(model).Select("kai as number").Where("code = ?", code.Code).Where("date = ?", da_min.Date).Scan(&kai)
	store.MysqlClient.GetDB().Model(model).Select("shou as number").Where("code = ?", code.Code).Where("date = ?", da_max.Date).Scan(&shou)
	store.MysqlClient.GetDB().Model(model).Select("sum(total_count) as number").Where("code = ?", code.Code).Where("date > ? and date <= ?", da_min.Date, da_max.Date).Scan(&liang)
	store.MysqlClient.GetDB().Model(model).Select("sum(turnover_rate) as number").Where("code = ?", code.Code).Where("date > ? and date <=?", da_min.Date, da_max.Date).Scan(&turnover_rate)

	var last dal.TicketHistoryMonth
	store.MysqlClient.GetDB().Model(&dal.TicketHistoryMonth{}).Where("code = ? and date = ?", code.Code, last_month_day).Find(&last)

	percent := (shou.Number - last.Shou) / last.Shou
	amplitude := (high.Number - low.Number) / last.Shou
	p, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", percent*100), 64)
	a, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", amplitude*100), 64)

	change, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", shou.Number-last.Shou), 64)
	turnover, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", turnover_rate.Number), 64)
	log.Println(high.Number, low.Number, kai.Number, shou.Number, liang.Number, turnover, p, a, change, da_max.Date, code.ID, code.Code, code.Name)
	h := dal.TicketHistoryMonth{Code: code.Code, Name: code.Name, Date: da_max.Date, Kai: kai.Number, Shou: shou.Number, High: high.Number, Low: low.Number,
		TotalCount: liang.Number, Percent: p, Change: change, Amplitude: a, TurnoverRate: turnover}
	store.MysqlClient.GetDB().Save(&h)
}
