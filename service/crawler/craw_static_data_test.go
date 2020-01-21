// @Time:       2020/1/6 下午1:11

package crawler

import (
	"fmt"
	"magic/stock/core/store"
	"magic/stock/dal"
	"strings"
	"testing"
)

// 从股票概念中抽出详细的概念保存在表中
func TestGetConcept(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)
	xx := map[string]bool{}
	for _, i := range code {
		x := strings.Split(strings.TrimRight(i.Concept, ","), ",")
		for _, j := range x {
			xx[j] = true
		}
	}

	for i, _ := range xx {
		if len(i) == 0 {
			continue
		}
		if strings.Contains(i, "概念") {
			c := dal.StockConcept{Name: i}
			store.MysqlClient.GetDB().Save(&c)
		} else {
			c := dal.StockLabels{Name: i}
			store.MysqlClient.GetDB().Save(&c)
		}
	}
}

// 上面获取到了个股的财务指标 以及 各项能力后 对按照大小的区间 对所有股打标
func TestCalcCaiWuForPreTicket(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)
	for _, i := range code {
		CalcCaiWuForPreTicket(i.Code)
	}
}

// 对每个字段自动打标吧 哈哈哈 注意： 需要大量的算力  差不多2小时左右
func CalcCaiWuForPreTicket(code string) {
	var per dal.StockPerTicket
	store.MysqlClient.GetDB().Model(&dal.StockPerTicket{}).Where("code = ?", code).Find(&per)
	fields := []string{"weifenpeilirun", "gubengongjijin", "jinzichanafter", "jiaquanshouyi", "yl_zongzichanlirunlv", "yl_gubenbaochoulv", "yl_jingzichanbaochoulv", "yl_zichanbaochoulv",
		"cz_zhuyingyewushouruzengzhanglv", "cz_jinglirunzengzhanglv", "cz_jingzichanzengzhanglv",
		"yy_zongzichanzhouzhuanglv", "yy_gudongquanyizhouzhuanglv"}
	var message string
	for _, field := range fields {
		var all []dal.StockPerTicket
		store.MysqlClient.GetDB().Model(&dal.StockPerTicket{}).Order(fmt.Sprintf("%s desc", field)).Find(&all)
		high := all[745]
		middle := all[1863]
		low := all[2608]
		bad := all[len(all)-1]
		switch field {
		case "weifenpeilirun":
			if per.Weifenpeilirun >= high.Weifenpeilirun {
				message += "未分配利润高; "
			} else if per.Weifenpeilirun >= middle.Weifenpeilirun {
				message += "未分配利润一般; "
			} else if per.Weifenpeilirun >= low.Weifenpeilirun {
				message += "未分配利润偏低; "
			} else if per.Weifenpeilirun >= bad.Weifenpeilirun {
				message += "未分配利润较差; "
			}
		case "gubengongjijin":
			if per.Gubengongjijin >= high.Gubengongjijin {
				message += "股本公积金高; "
			} else if per.Gubengongjijin >= middle.Gubengongjijin {
				message += "股本公积金一般; "
			} else if per.Gubengongjijin >= low.Gubengongjijin {
				message += "股本公积金偏低; "
			} else if per.Gubengongjijin >= bad.Gubengongjijin {
				message += "股本公积金较差; "
			}
		case "jinzichanafter":
			if per.Jinzichanafter >= high.Jinzichanafter {
				message += "每股净资产高; "
			} else if per.Jinzichanafter >= middle.Jinzichanafter {
				message += "每股净资产一般; "
			} else if per.Jinzichanafter >= low.Jinzichanafter {
				message += "每股净资产偏低; "
			} else if per.Jinzichanafter >= bad.Jinzichanafter {
				message += "每股净资产较差; "
			}
		case "jiaquanshouyi":
			if per.Jiaquanshouyi >= high.Jiaquanshouyi {
				message += "加权每股收益高; "
			} else if per.Jiaquanshouyi >= middle.Jiaquanshouyi {
				message += "加权每股收益一般; "
			} else if per.Jiaquanshouyi >= low.Jiaquanshouyi {
				message += "加权每股收益偏低; "
			} else if per.Jiaquanshouyi >= bad.Jiaquanshouyi {
				message += "加权每股收益较差; "
			}
		case "yl_zongzichanlirunlv":
			if per.YlZongzichanlirunlv >= high.YlZongzichanlirunlv {
				message += "总资产利润率高; "
			} else if per.YlZongzichanlirunlv >= middle.YlZongzichanlirunlv {
				message += "总资产利润率一般; "
			} else if per.YlZongzichanlirunlv >= low.YlZongzichanlirunlv {
				message += "总资产利润率偏低; "
			} else if per.YlZongzichanlirunlv >= bad.YlZongzichanlirunlv {
				message += "总资产利润率较差; "
			}
		case "yl_gubenbaochoulv":
			if per.YlGubenbaochoulv >= high.YlGubenbaochoulv {
				message += "股本报酬率高; "
			} else if per.YlGubenbaochoulv >= middle.YlGubenbaochoulv {
				message += "股本报酬率一般; "
			} else if per.YlGubenbaochoulv >= low.YlGubenbaochoulv {
				message += "股本报酬率偏低; "
			} else if per.YlGubenbaochoulv >= bad.YlGubenbaochoulv {
				message += "股本报酬率较差; "
			}
		case "yl_jingzichanbaochoulv":
			if per.YlJingzichanbaochoulv >= high.YlJingzichanbaochoulv {
				message += "净资产报酬率高; "
			} else if per.YlJingzichanbaochoulv >= middle.YlJingzichanbaochoulv {
				message += "净资产报酬率率一般; "
			} else if per.YlJingzichanbaochoulv >= low.YlJingzichanbaochoulv {
				message += "净资产报酬率偏低; "
			} else if per.YlJingzichanbaochoulv >= bad.YlJingzichanbaochoulv {
				message += "净资产报酬率较差; "
			}
		case "yl_zichanbaochoulv":
			if per.YlZichanbaochoulv >= high.YlZichanbaochoulv {
				message += "资产报酬率高; "
			} else if per.YlZichanbaochoulv >= middle.YlZichanbaochoulv {
				message += "资产报酬率率一般; "
			} else if per.YlZichanbaochoulv >= low.YlZichanbaochoulv {
				message += "资产报酬率偏低; "
			} else if per.YlZichanbaochoulv >= bad.YlZichanbaochoulv {
				message += "资产报酬率较差; "
			}
		case "cz_zhuyingyewushouruzengzhanglv":
			if per.CzZhuyingyewushouruzengzhanglv >= high.CzZhuyingyewushouruzengzhanglv {
				message += "主营业务收入增长率高; "
			} else if per.CzZhuyingyewushouruzengzhanglv >= middle.CzZhuyingyewushouruzengzhanglv {
				message += "主营业务收入增长率一般; "
			} else if per.CzZhuyingyewushouruzengzhanglv >= low.CzZhuyingyewushouruzengzhanglv {
				message += "主营业务收入增长率偏低; "
			} else if per.CzZhuyingyewushouruzengzhanglv >= bad.CzZhuyingyewushouruzengzhanglv {
				message += "主营业务收入增长率较差; "
			}
		case "cz_jinglirunzengzhanglv":
			if per.CzJinglirunzengzhanglv >= high.CzJinglirunzengzhanglv {
				message += "净利润增长率高; "
			} else if per.CzJinglirunzengzhanglv >= middle.CzJinglirunzengzhanglv {
				message += "净利润增长率一般; "
			} else if per.CzJinglirunzengzhanglv >= low.CzJinglirunzengzhanglv {
				message += "净利润增长率偏低; "
			} else if per.CzJinglirunzengzhanglv >= bad.CzJinglirunzengzhanglv {
				message += "净利润增长率较差; "
			}
		case "cz_jingzichanzengzhanglv":
			if per.CzJingzichanzengzhanglv >= high.CzJingzichanzengzhanglv {
				message += "净资产增长率高; "
			} else if per.CzJingzichanzengzhanglv >= middle.CzJingzichanzengzhanglv {
				message += "净资产增长率一般; "
			} else if per.CzJingzichanzengzhanglv >= low.CzJingzichanzengzhanglv {
				message += "净资产增长率偏低; "
			} else if per.CzJingzichanzengzhanglv >= bad.CzJingzichanzengzhanglv {
				message += "净资产增长率较差; "
			}
		case "yy_zongzichanzhouzhuanglv":
			if per.YyZongzichanzhouzhuanglv >= high.YyZongzichanzhouzhuanglv {
				message += "企业销售能力/投资收益高; "
			} else if per.YyZongzichanzhouzhuanglv >= middle.YyZongzichanzhouzhuanglv {
				message += "企业销售能力/投资收益一般; "
			} else if per.YyZongzichanzhouzhuanglv >= low.YyZongzichanzhouzhuanglv {
				message += "企业销售能力/投资收益偏低; "
			} else if per.YyZongzichanzhouzhuanglv >= bad.YyZongzichanzhouzhuanglv {
				message += "企业销售能力/投资收益较差; "
			}
		case "yy_gudongquanyizhouzhuanglv":
			if per.YyGudongquanyizhouzhuanglv >= high.YyGudongquanyizhouzhuanglv {
				message += "资产效率/营运能力高; "
			} else if per.YyGudongquanyizhouzhuanglv >= middle.YyGudongquanyizhouzhuanglv {
				message += "资产效率/营运能力一般; "
			} else if per.YyGudongquanyizhouzhuanglv >= low.YyGudongquanyizhouzhuanglv {
				message += "资产效率/营运能力偏低; "
			} else if per.YyGudongquanyizhouzhuanglv >= bad.YyGudongquanyizhouzhuanglv {
				message += "资产效率/营运能力较差; "
			}
		}
	}
	per.RankCaiwu = message
	store.MysqlClient.GetDB().Save(&per)
}

func TestMultiQuery(t *testing.T) {
	var c []dal.Predict
	err := store.MysqlClient.GetDB().Model(&dal.Predict{}).
		Where("date = ?", "2020-01-03").
		//Where("`condition` regexp ?", "高位回调").
		//Where("`condition` regexp ?", "金叉").
		Where("`condition` regexp ?", "近期60日均线与收盘价黏合").
		Find(&c).Error
	fmt.Println(err)
	for _, i := range c {
		fmt.Println(i.Code)
	}
}

// 从concepts中拿出盘口信息
func TestGetTape(t *testing.T) {
	store.MysqlClient.GetOnlineDB().Model(&dal.Code{}).Where("concept like ?", "%超大盘%").Update("tape", "超大盘")
	store.MysqlClient.GetOnlineDB().Model(&dal.Code{}).Where("concept like ?", "%大盘%").Update("tape", "大盘")
	store.MysqlClient.GetOnlineDB().Model(&dal.Code{}).Where("concept like ?", "%中盘%").Update("tape", "中盘")
	store.MysqlClient.GetOnlineDB().Model(&dal.Code{}).Where("concept like ?", "%小盘%").Update("tape", "小盘")
}

// 更新历史数据存在 low=0的情况
func TestUpdateLow(t *testing.T) {
	var history []dal.TicketHistoryMonth
	store.MysqlClient.GetDB().Model(&dal.TicketHistoryMonth{}).Where("low = ?", 0).Find(&history)
	for _, i := range history {
		i.Low = i.Shou
		i.High = i.Shou
		store.MysqlClient.GetDB().Save(&i)
	}
}
