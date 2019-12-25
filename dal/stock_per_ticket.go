// @Time:       2019/12/2 下午4:15

package dal

// 每股指标
type StockPerTicket struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Code string `gorm:"index" json:"code"`

	// 每股指标
	Tanboshouyi        float64 `json:"tanboshouyi"`        // 摊薄每股收益(元)
	Jiaquanshouyi      float64 `json:"jiaquanshouyi"`      // 加权每股收益(元) 使用加权平均法来计算每股收益，这样才可以更准确、更合理地反映公司客观的盈利能力。
	Jinzichanfront     float64 `json:"jinzichanfront"`     // 每股净资产_调整前(元)
	Shouyiafter        float64 `json:"shouyiafter"`        // 每股收益_调整后(元)  指扣除与主营业务无关的一次性损益后的净利润除以总股本得出的每股收益。
	Jinzichanafter     float64 `json:"jinzichanafter"`     // 每股净资产_调整后(元) 这一指标反映每股股票所拥有的资产现值。每股净资产越高，股东拥有的每股资产价值越多；
	Jingyingxianjinliu float64 `json:"jingyingxianjinliu"` // 每股经营性现金流(元) 即每股经营活动产生的现金流量净额
	Gubengongjijin     float64 `json:"gubengongjijin"`     // 每股资本公积金(元) 资本公积金是指从公司的利润以外的收入中提取的一种公积金。其主要来源有股票溢价收入，财产重估增值，以及接受捐赠资产等。每股资本公积金=资本公积金/总股本
	Weifenpeilirun     float64 `json:"weifenpeilirun"`     // 每股未分配利润(元) 每股未分配利润越多，不仅表明该公司盈利能力强，也意味着该公司未来分红、送股的能力强

	// 盈利能力
	YlZongzichanlirunlv     float64 `json:"yl_zongzichanlirunlv"`     // 总资产利润率(%)  总资产利润率=利润总额/资产平均总额×100% 可用来说明企业运用其全部资产获取利润的能力。
	YlZhuyingyewulirunlv    float64 `json:"yl_zhuyingyewulirunlv"`    // 主营业务利润率(%) 主营利润率，是公司主业所产生的利润率。比如公司主业是房地产，那么经营房地产所产生的利润，与主营业务收入的比率，就是主营利润率。
	YlZongzichanjinglirunlv float64 `json:"yl_zongzichanjinglirunlv"` // 总资产净利润率(%) // 又称总资产收益率，是企业净利润总额与企业资产平均总额的比率，即过去所说的资金利润率
	YlYingyelirunlv         float64 `json:"yl_yingyelirunlv"`         // 营业利润率(%)  // 营业利润率是企业付清一切帐项后剩下的金额称为利润。在会计学上，利润可分为毛利
	YlXiaoshoujinglilv      float64 `json:"yl_xiaoshoujinglilv"`      // 销售净利率(%) 是净利润占销售收入的百分比。 该指标反映每一元销售收入带来的净利润的多少，表示销售收入的收益水平。
	YlGubenbaochoulv        float64 `json:"yl_gubenbaochoulv"`        // 股本报酬率(%)  股本报酬率是指公司税后利润与其股本的比率，表明公司股本总额中平均每百元股本所获得的纯利润。
	YlJingzichanbaochoulv   float64 `json:"yl_jingzichanbaochoulv"`   // 净资产报酬率(%) // 该指标反映股东权益的收益水平，用以衡量公司运用自有资本的效率。指标值越高，说明投资带来的收益越高。该指标体现了自有资本获得净收益的能力。
	YlZichanbaochoulv       float64 `json:"yl_zichanbaochoulv"`       // 资产报酬率(%) 用以评价企业运用全部资产的总体获利能力，是评价企业资产运营效益的重要指标。

	// 成长能力
	CzZhuyingyewushouruzengzhanglv float64 `json:"cz_zhuyingyewushouruzengzhanglv"` // 主营业务收入增长率(%)
	CzJinglirunzengzhanglv         float64 `json:"cz_jinglirunzengzhanglv"`         // 净利润增长率(%)  净利润增长率是指企业当期净利润比上期净利润的增长幅度，指标值越大代表企业盈利能力越强。
	CzJingzichanzengzhanglv        float64 `json:"cz_jingzichanzengzhanglv"`        // 净资产增长率(%)
	CzZongzichanzengzhanglv        float64 `json:"cz_zongzichanzengzhanglv"`        // 总资产增长率(%)

	// 运营能力
	YyYingshouzhangkuanzhouzhuanlv float64 `json:"yy_yingshouzhangkuanzhouzhuanlv"` // 应收账款周转率(次) 应收账款周转率越高越好，应收账示周转率高，表明收账迅速，账龄较短；
	YyCunhuozhouzhuanglv           float64 `json:"yy_cunhuozhouzhuanglv"`           // 存货周转率(次) 存货周转率越高，表明企业存货资产变现能力越强，存货及占用在存货上的资金周转速度越快。
	YyLiudongzichanzhouzhuanglv    float64 `json:"yy_liudongzichanzhouzhuanglv"`    // 流动资产周转率(次) 该指标越高，说明企业流动资产的利用效率越好。
	YyZongzichanzhouzhuanglv       float64 `json:"yy_zongzichanzhouzhuanglv"`       // 总资产周转率(次) 总资产周转率越高，说明企业销售能力越强,资产投资的效益越好
	YyGudongquanyizhouzhuanglv     float64 `json:"yy_gudongquanyizhouzhuanglv"`     // 股东权益周转率(次)  指标说明公司运用所有制的资产的效率。 该比率越高，表明所有者资产的运用效率高，营运能力强

	RankCaiwu string `json:"rank_caiwu"` // 财务状况打标
	Date      string `json:"date"`
}

func (StockPerTicket) TableName() string {
	return "magic_stock_per_ticket"
}
