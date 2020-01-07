// @Time:       2019/12/1 下午3:11

package control

import (
	"errors"
	"fmt"
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/service/adapter"
	"magic/stock/service/check"
	"magic/stock/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/fatih/set.v0"
)

type PredictIF interface {
	Query(where string, args []interface{}) (*dal.Predict, error)
	QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.Predict, error)
	Exist(where string, args []interface{}) bool
	// 获取预测数据列表 post 请求
	PredictList(c *gin.Context)
	// 获取股票详情
	GetDetail(c *gin.Context)
	GetFunds(c *gin.Context)
	// 通过机构code查询机构持仓
	FundHold(c *gin.Context)
	// 通过名称查询流通股东可能存在的其它持仓
	TopHolderHold(c *gin.Context)

	GetConditions(c *gin.Context)
	GetHighConditions(c *gin.Context)
	// 获取基本条件列表
	GetQueryList(c *gin.Context)

	GetFenHong(c *gin.Context)
	GetPeiGuZhuangZeng(c *gin.Context)
	// 获取所有子公司
	GetSubComp(c *gin.Context)
	Response(c *gin.Context, data interface{}, err error, param ...int)
}

var (
	PredictControlGlobal PredictIF
	OrderLimit           = []string{"id", "score", "percent", "price", "fund_count", "sm_count"}
	LocationList         = strings.Split("安徽,北京,福建,甘肃,广东,广西,贵州,海南,河北,河南,黑龙江,湖北,湖南,吉林,江苏,江西,辽宁,内蒙古,宁夏,青海,山东,山西,陕西,上海,四川,天津,西藏,新疆,云南,浙江", ",")
	BelongList           = strings.Split("保险业,玻璃及玻璃制品业,采掘服务业,餐饮业,仓储业,出版业,电机制造业,电力、蒸汽、热水的生产和供应业,电力生产业,电器机械及器材制造业,电信服务业,电子元器件制造业,房地产管理业,房地产开发与经营业,房地产中介服务业,纺织业,非金属矿物制品业,服装及其他纤维制品制造业,港口业,公共设施服务业,公路管理及养护业,公路运输业,管道运输业,广播电视设备制造业,广播电影电视业,贵金属矿采选业,贵金属冶炼业,海洋渔业,航空航天器制造业,航空运输业,黑色金属矿采选业,黑色金属冶炼及压延加工业,化学肥料制造业,化学农药制造业,化学纤维制造业,化学原料及化学制品制造业,计算机及相关设备制造业,计算机软件开发与咨询,计算机应用服务业,交通运输辅助业,交通运输设备制造业,金融信托业,金属制品业,酒精及饮料酒制造业,炼钢业,粮食及饲料加工业,林业,零售业,旅馆业,旅游业,轮胎制造业,毛纺织业,煤气生产和供应业,煤炭采选业,摩托车制造业,能源、材料和机械电子设备批发业,能源批发业,农、林、牧、渔服务业,农业,普通机械制造业,其他传播、文化产业,其他电子设备制造业,其他非金属矿物制品业,其他交通运输业,其他金融业,其他批发业,其他社会服务业,其他通用零部件制造业,其他专用设备制造业,其它,汽车制造业,轻纺工业专用设备制造业,人寿保险业,日用电器制造业,日用化学产品制造业,乳制品制造业,软饮料制造业,商业经纪与代理业,生物药品制造业,生物制品业,石化及其他工业专用设备制造业,石墨及碳素制品业,石油和天然气开采服务业,石油和天然气开采业,石油加工及炼焦业,食品、饮料、烟草和家庭用品批发业,食品加工业,食品制造业,市内公共交通业,输配电及控制设备制造业,水泥制造业,水上运输业,塑料制造业,陶瓷制品业,铁路、公路、隧道、桥梁建筑业,铁路运输设备制造业,铁路运输业,通信服务业,通信及相关设备制造业,通信设备制造业,通用设备制造业,土木工程建筑业,卫生、保健、护理服务业,文教体育用品制造业,稀有稀土金属冶炼业,橡胶制造业,信息传播服务业,畜牧业,沿海运输业,药品及医疗器械零售业,冶金、矿山、机电工业专用设备制造业,医疗器械制造业,医药制造业,仪器仪表及文化、办公用机械制造业,银行业,饮料制造业,印刷业,有色金属矿采选业,有色金属冶炼及压延加工业,渔业,渔业服务业,远洋运输业,造纸及纸制品业,造纸业,照明器具制造业,证券、期货业,制糖业,制造业,中药材及中成药加工业,重有色金属冶炼业,轴承、阀门制造业,专业、科研服务业,专用设备制造业,装修装饰业,自来水的生产和供应业,自行车制造业,租赁服务业", ",")
	ConceptList          = strings.Split("360概念,5G概念,6G概念,AMC概念,ETC概念,GDR概念,IGBT概念,IP概念,LCP概念,MLCC概念,NFC概念,PCB概念,PPP概念,TOF概念,VPN概念,阿里概念,百度概念,博彩概念,超导概念,单车概念,电商概念,抖音概念,独角兽概念,分拆概念,芬太尼概念,风能概念,富时罗素概念,富士康概念,供销社概念,谷歌概念,股期概念,骨传导概念,华为概念,黄金概念,机器人概念,基因概念,激光概念,甲醇概念,进博会概念,京东概念,空客概念,快递概念,沥青概念,脸书概念,麻醉概念,美团概念,纳米概念,铌概念,宁德时代概念,硼墨稀概念,苹果概念,期货概念,前海概念,全息概念,人造肉概念,赛马概念,三沙概念,涉矿概念,陶瓷概念,腾讯概念,体育概念,网易概念,微信概念,蔚来汽车概念,消防概念,小米概念,芯片概念,锌溴概念,信托概念,眼科概念,婴童概念,摘帽概念,征信概念,中芯国际概念,重组概念,猪肉概念,足球概念", ",")
	LabelList            = strings.Split("3D打印,5G,GPU,H股,IPV6,LED,MLED,MSCI中国,OLED,PM2.5,QFII重仓,QLED,SAAS,ST板块,阿尔茨海默,安防,安防服务,安徽国资,白酒,百货O2O,半导体,保险重仓,北斗导航,北京国资,本月解禁,边缘计算,冰雪产业,兵工集团,兵装集团,彩票,参股金融,参股新三板,操作系统,草甘膦,长三角一体化,长株潭,超大盘,超高清,超级电容,超级高铁,超级细菌,超级真菌,超宽带,朝鲜改革,车联网,成渝城市群,成渝特区,充电桩,出口退税,储能,触摸屏,传感器,船舶工业集团,船舶重工集团,创投,创业板壳预期,磁悬浮,次新股,次新开板,大豆,大飞机,大盘,大数据,大唐集团,单晶硅,低碳经济,迪士尼,地理信息,地热能,地下管廊,电池管理,电动物流车,电改,电力物联网,电子发票,电子化学品,电子竞技,电子商务,电子烟,电子支付,定制家居,东北振兴,东亚自贸,动力煤,动物疫苗,短视频,多晶硅,二胎,钒电池,房屋租赁,仿制药,分拆上市,风能,氟化工,福建国资,福建自贸区,复兴号,富勒烯,覆铜板,盖板玻璃,干细胞,甘肃国资,钢铁互联网,港口运输,高出口占比,高送转,高铁,高校背景,高校系,工程机械,工业4.0,工业大麻,工业互联网,供给侧改革,供应链金融,共享经济,共享汽车,股权激励,股权转让,钴镍,固废处理,固态电池,光刻胶,光通信,光学,广电系,广东国资,广西国资,轨道交通,贵州国资,国产乳业,国产软件,国电集团,国防军工,国机集团,国企改革,国资改革,国资入股,海底光缆,海工装备,海绵城市,海南国资,海南自贸,海南自贸区,海上丝路,海水淡化,海峡西岸,含B股,含GDR,含H股,杭州亚运,航空工业集团,航母产业,航天航空,航天军工,航天科工集团,航天科技集团,合同能源,核电,核电核能,黑龙江自贸区,横琴新区,湖南国资,互联金融,互联网保险,互联网家装,互联医疗,沪警示板,沪自贸区,华为电视,华为海思,华为鲲鹏,环境监测,环球影城,黄河三角,黄金股,混改,机器视觉,机械,基金重仓,基因测序,集成电路,家具卫浴,家用电器,嘉兴地区,甲型流感,价值成长,价值品牌,建筑节能,健康中国,江苏国资,降解材料,铰链,教育产业,节能环保,金刚线,金融IC,金融参股,金融改革,金融机具,金融科技,金属铝,京津冀,精准医疗,举牌,聚氨酯,军工电子,军工航天,军民融合,喀什规划区,抗癌,抗癌药物,抗癌治癌,抗流感,抗震防震,壳资源,可燃冰,空气治理,跨境电商,宽带提速,垃圾发电,垃圾分类,兰州自贸区,冷链物流,锂电池,锂硫电池,量子通信,磷化工,零售药店,流媒体,旅游酒店,绿色照明,轮胎,蚂蚁金服,煤化工,镁空气电池,免疫治疗,民爆,民营医院,民营银行,民用航空,内贸规划,内容审核,纳米银线,能源互联,能源纸,农村电商,农村电网,农机,啤酒,平潭实验区,苹果产业链,苹果三星,汽车电子,汽车后市场,汽车金融,汽车零部件,汽车轻量化,禽流感药物,青岛,青蒿素,氢能源,氢燃料,清洁能源,区块链,全面屏,券商相关,券商重仓,燃料电池,染料涂料,人工智能,人脸识别,人脑工程,日韩贸易,融资融券,融资租赁,柔性电子,三元锂电,山西国资,陕甘宁,商业保理,上海本地,上海国资,上海自贸,奢侈品,设计咨询,社保重仓,深ST板,深汕合作区,深圳本地,深圳国资,生态林业,生态农业,生物农药,生物燃料,生物识别,生物医药,生物疫苗,生物育种,生物质能,石家庄,石墨烯,食品安全,食品检测,食药追溯,世界杯,手势识别,手游,数字货币,数字孪生,数字丝绸,数字中国,水产品,水电,水利工程,水利建设,水泥,水域改革,丝绸之路,饲料,送转潜力,胎压监测,太空互联网,太阳能,态势感知,钛白粉,碳捕集,碳纤维,特高压,特色小镇,特斯拉,特种玻璃,特种钢,腾讯云,体外诊断,体育产业,天津国资,天津自贸,天然气,天文观测,铁路基建,通用航空,透明工厂,图们江,土地流转,土壤修复,退市警示,外资背景,皖江区域,网红经济,网络安全,网络切片,网络营销,网络游戏,网络直播,维生素,尾气处理,卫星导航,未股改,文化振兴,污水处理,无感支付,无人机,无人驾驶,无人零售,无锡国资委,无线充电,无线耳机,武汉规划,物联网,西部开发,西藏国资,稀缺资源,稀土永磁,细胞治疗,乡村振兴,消费电子,小贷典当,小金属,小盘,谐振器,锌电池,新版医保,新材料,新基建,新疆建设兵团,新疆振兴,新零售,新能源,新能源车,信托重仓,信息安全,雄安地产,雄安环保,雄安基建,雄安交运,雄安金融,雄安能化,雄安新区,休闲食品,虚拟现实,血液制品,循环经济,沿海发展,盐湖提锂,盐化工,央企50,央企改革,央企金控,养鸡,养老产业,业绩预降,业绩预升,页岩气,液态金属,一带一路,一线蓝筹,医疗美容,医疗器械,医药电商,医用耗材,移动支付,乙二醇,乙肝疫苗,疫苗,引力波,应急管理,影视传媒,影视动漫,油气改革,油气管网,油气勘探,有机硅,余热发电,预盈预增,员工持股,园区开发,远洋运输,粤港澳,云计算,云视频,云印刷,云游戏,在线教育,在线旅游,增持回购,增强现实,债转股,浙江国资,针状焦,振兴沈阳,整车,整体上市,证金汇金,支付牌照,知识产权,职业教育,指纹识别,智慧城市,智慧农业,智慧停车,智慧物流,智慧医疗,智慧政务,智能穿戴,智能电视,智能电网,智能机器,智能家居,智能交通,智能汽车,智能手表,智能手机,智能音箱,中概股回归,中国化工集团,中国建材集团,中韩自贸区,中核工业集团,中化集团,中科院系,中粮集团,中盘,中石化系,中药,中字头,种业,重工装备,重庆国改,舟山自贸区,猪肉,装配建筑,准ST股,资产注入,自贸区,自由贸易港", ",")
	FormList             = strings.Split("民营企业,民营相对控股企业,国有企业,国有相对控股企业,中外合资企业,其它", ",")
)

func init() {
	tmp := new(PredictControl)
	tmp.service = adapter.PredictServiceGlobal
	tmp.response = new(model.HttpResponse)
	PredictControlGlobal = tmp
}

type PredictControl struct {
	service  adapter.PredictServiceIF
	response *model.HttpResponse
}

func (u *PredictControl) Query(where string, args []interface{}) (*dal.Predict, error) {
	return u.service.Query(where, args)
}

func (u *PredictControl) QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.Predict, error) {
	return u.service.QueryAll(where, args, offset, limit, select_only)
}

func (u *PredictControl) Exist(where string, args []interface{}) bool {
	c, _ := u.service.Count(where, args)
	return c > 0
}

func (d *PredictControl) Response(c *gin.Context, data interface{}, err error, param ...int) {
	c.AbortWithStatusJSON(200, d.response.Response(data, err, param...))
}

func (d *PredictControl) getMinMax(da map[string]float64) (min, max float64, is_nil bool) {
	min, ok_min := da["min"]
	max, ok_max := da["max"]
	if !ok_min || !ok_max {
		return 0, 0, true
	}
	if min == -10000 && max == 10000 {
		return 0, 0, true
	}
	return min, max, false
}

func (d *PredictControl) ParseStockPerTicket(param map[string]float64, field string) set.Interface {
	tmp := set.New(set.ThreadSafe)
	min, max, is_nil := d.getMinMax(param)
	if is_nil {
		return tmp
	}
	var codes []Codes
	store.MysqlClient.GetDB().Model(&dal.StockPerTicket{}).Select("code").Where(fmt.Sprintf("%s >= ? and %s <= ?", field, field), min, max).Debug().Scan(&codes)
	for _, i := range codes {
		tmp.Add(i.Code)
	}
	if tmp.Size() > 0 {
		return tmp
	}
	return nil
}

func (d *PredictControl) ParseCount(param map[string]float64, date string, field string) set.Interface {
	tmp := set.New(set.ThreadSafe)
	min, max, is_nil := d.getMinMax(param)
	if is_nil {
		return tmp
	}
	var codes []Codes
	store.MysqlClient.GetDB().Model(&dal.Predict{}).Select("code").Where(fmt.Sprintf("date = ? and %s >= ? and %s <= ?", field, field), date, min, max).Debug().Scan(&codes)
	for _, i := range codes {
		tmp.Add(i.Code)
	}
	if tmp.Size() > 0 {
		return tmp
	}
	return nil
}

func (d *PredictControl) ParseLastDayRange(param map[string]float64, date string, field string) set.Interface {
	tmp := set.New(set.ThreadSafe)
	min, max, is_nil := d.getMinMax(param)
	if is_nil {
		return tmp
	}
	var codes []Codes
	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Select("code").Where(fmt.Sprintf("date = ? and %s >= ? and %s <= ?", field, field), date, min, max).Scan(&codes)
	for _, i := range codes {
		tmp.Add(i.Code)
	}
	if tmp.Size() > 0 {
		return tmp
	}
	return nil
}

type Codes struct {
	Code string
}

func (d *PredictControl) doQueryLeft(authentication *model.AuthResult) error {
	if !authentication.Member {
		if authentication.QueryLeft == 0 {
			return errors.New("查询次数不足")
		} else {
			user_obj, _ := UserControlGlobal.Query("id = ?", []interface{}{authentication.Uid})
			left := user_obj.QueryLeft - 1
			exp := user_obj.Exp + 1
			user_obj.QueryLeft = left
			user_obj.Exp = exp
			err := store.MysqlClient.GetDB().Save(&user_obj).Error
			log.Println("查询次数剩余", authentication.User, authentication.QueryLeft, left, err)
		}
	}
	return nil
}

func (d *PredictControl) PredictList(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	var post model.GetPredicts
	err := c.BindJSON(&post)
	if err != nil {
		d.Response(c, nil, err)
		return
	}

	err = d.doQueryLeft(authentication)
	if err != nil {
		d.Response(c, nil, err)
		return
	}
	offset, limit := check.ParamParse.GetPagination(c)
	// 如果用户提交查询并保存查询结果
	if post.Save {
		//err := adapter.UserServiceGlobal.SaveUserConditions(&post, authentication)
		//if err != nil {
		//	log.Println("保存用户查询数据失败", err)
		//}
	}
	var where_belongs, where_locations, where_concepts, where_forms []string
	var args_belongs, args_locationgs, args_concepts, args_forms []interface{}

	belong_set := set.New(set.ThreadSafe)
	form_set := set.New(set.ThreadSafe)
	location_set := set.New(set.ThreadSafe)
	concept_set := set.New(set.ThreadSafe)
	per_ticket_set1 := set.New(set.ThreadSafe)
	per_ticket_set2 := set.New(set.ThreadSafe)
	per_ticket_set3 := set.New(set.ThreadSafe)
	per_ticket_set4 := set.New(set.ThreadSafe)
	per_ticket_set5 := set.New(set.ThreadSafe)
	per_ticket_set6 := set.New(set.ThreadSafe)

	last_day_set1 := set.New(set.ThreadSafe)
	last_day_set2 := set.New(set.ThreadSafe)
	last_day_set3 := set.New(set.ThreadSafe)
	last_day_set4 := set.New(set.ThreadSafe)
	last_day_set5 := set.New(set.ThreadSafe)

	ability_set1 := set.New(set.ThreadSafe)
	ability_set2 := set.New(set.ThreadSafe)
	ability_set3 := set.New(set.ThreadSafe)
	ability_set4 := set.New(set.ThreadSafe)
	ability_set5 := set.New(set.ThreadSafe)
	ability_set6 := set.New(set.ThreadSafe)
	ability_set7 := set.New(set.ThreadSafe)
	ability_set8 := set.New(set.ThreadSafe)
	ability_set9 := set.New(set.ThreadSafe)
	ability_set10 := set.New(set.ThreadSafe)
	ability_set11 := set.New(set.ThreadSafe)
	ability_set12 := set.New(set.ThreadSafe)
	ability_set13 := set.New(set.ThreadSafe)
	ability_set14 := set.New(set.ThreadSafe)
	ability_set15 := set.New(set.ThreadSafe)
	ability_set16 := set.New(set.ThreadSafe)
	ability_set17 := set.New(set.ThreadSafe)

	other_set1 := set.New(set.ThreadSafe)
	other_set2 := set.New(set.ThreadSafe)
	other_set3 := set.New(set.ThreadSafe)
	other_set4 := set.New(set.ThreadSafe)
	other_set5 := set.New(set.ThreadSafe)
	other_set6 := set.New(set.ThreadSafe)
	other_set7 := set.New(set.ThreadSafe)

	if len(post.Query.Belongs) > 0 {
		var codes []Codes
		for _, i := range post.Query.Belongs {
			where_belongs = append(where_belongs, "belong = ?")
			args_belongs = append(args_belongs, i)
		}
		where_str := strings.Join(where_belongs, " OR ")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_belongs...).Scan(&codes)
		for _, i := range codes {
			belong_set.Add(i.Code)
		}
	}

	if len(post.Query.Forms) > 0 {
		var codes []Codes
		for _, i := range post.Query.Forms {
			where_forms = append(where_forms, "organizational_form = ?")
			args_forms = append(args_forms, i)
		}
		where_str := strings.Join(where_forms, " OR ")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_forms...).Scan(&codes)
		for _, i := range codes {
			form_set.Add(i.Code)
		}
	}

	if len(post.Query.Locations) > 0 {
		var codes []Codes
		for _, i := range post.Query.Locations {
			where_locations = append(where_locations, "location = ?")
			args_locationgs = append(args_locationgs, i)
		}
		where_str := strings.Join(where_locations, " OR ")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_locationgs...).Scan(&codes)
		for _, i := range codes {
			location_set.Add(i.Code)
		}
	}

	if len(post.Query.Concepts) > 0 || len(post.Query.Labels) > 0 {
		var codes []Codes
		arrays := append(post.Query.Concepts, post.Query.Labels...)
		for _, i := range arrays {
			where_concepts = append(where_concepts, "concept like ?")
			args_concepts = append(args_concepts, "%"+i+"%")
		}
		where_str := strings.Join(where_concepts, " OR ")
		store.MysqlClient.GetDB().Model(&dal.Code{}).Select("code").Where(where_str, args_concepts...).Scan(&codes)
		for _, i := range codes {
			concept_set.Add(i.Code)
		}
	}
	per_ticket_set1 = d.ParseStockPerTicket(post.Query.PerTickets.Shouyiafter, "shouyiafter")
	per_ticket_set2 = d.ParseStockPerTicket(post.Query.PerTickets.Jiaquanshouyi, "jiaquanshouyi")
	per_ticket_set3 = d.ParseStockPerTicket(post.Query.PerTickets.Jinzichanafter, "jinzichanafter")
	per_ticket_set4 = d.ParseStockPerTicket(post.Query.PerTickets.Jingyingxianjinliu, "jingyingxianjinliu")
	per_ticket_set5 = d.ParseStockPerTicket(post.Query.PerTickets.Gubengongjijin, "gubengongjijin")
	per_ticket_set6 = d.ParseStockPerTicket(post.Query.PerTickets.Weifenpeilirun, "weifenpeilirun")

	last_day_set1 = d.ParseLastDayRange(post.Query.LastDayRange.LastPercent, post.Date, "percent")
	last_day_set2 = d.ParseLastDayRange(post.Query.LastDayRange.LastAmplitude, post.Date, "amplitude")
	last_day_set3 = d.ParseLastDayRange(post.Query.LastDayRange.LastTurnoverrate, post.Date, "turnover_rate")
	last_day_set4 = d.ParseLastDayRange(post.Query.LastDayRange.LastPrice, post.Date, "shou")
	last_day_set5 = d.ParseLastDayRange(post.Query.LastDayRange.LastNumberRate, post.Date, "number_rate")

	// 盈利能力
	ability_set1 = d.ParseStockPerTicket(post.Query.YlAbility.YlZongzichanlirunlv, "yl_zongzichanlirunlv")
	ability_set2 = d.ParseStockPerTicket(post.Query.YlAbility.YlZhuyingyewulirunlv, "yl_zhuyingyewulirunlv")
	ability_set3 = d.ParseStockPerTicket(post.Query.YlAbility.YlZongzichanjinglirunlv, "yl_zongzichanjinglirunlv")
	ability_set4 = d.ParseStockPerTicket(post.Query.YlAbility.YlYingyelirunlv, "yl_yingyelirunlv")
	ability_set5 = d.ParseStockPerTicket(post.Query.YlAbility.YlXiaoshoujinglilv, "yl_xiaoshoujinglilv")
	ability_set6 = d.ParseStockPerTicket(post.Query.YlAbility.YlGubenbaochoulv, "yl_gubenbaochoulv")
	ability_set7 = d.ParseStockPerTicket(post.Query.YlAbility.YlJingzichanbaochoulv, "yl_jingzichanbaochoulv")
	ability_set8 = d.ParseStockPerTicket(post.Query.YlAbility.YlZichanbaochoulv, "yl_zichanbaochoulv")
	// 成长能力
	ability_set9 = d.ParseStockPerTicket(post.Query.CzAbility.CzZhuyingyewushouruzengzhanglv, "cz_zhuyingyewushouruzengzhanglv")
	ability_set10 = d.ParseStockPerTicket(post.Query.CzAbility.CzJinglirunzengzhanglv, "cz_jinglirunzengzhanglv")
	ability_set11 = d.ParseStockPerTicket(post.Query.CzAbility.CzJingzichanzengzhanglv, "cz_jingzichanzengzhanglv")
	ability_set12 = d.ParseStockPerTicket(post.Query.CzAbility.CzZongzichanzengzhanglv, "cz_zongzichanzengzhanglv")
	// 运营能力
	ability_set13 = d.ParseStockPerTicket(post.Query.YyAbility.YyYingshouzhangkuanzhouzhuanlv, "yy_yingshouzhangkuanzhouzhuanlv")
	ability_set14 = d.ParseStockPerTicket(post.Query.YyAbility.YyCunhuozhouzhuanglv, "yy_cunhuozhouzhuanglv")
	ability_set15 = d.ParseStockPerTicket(post.Query.YyAbility.YyLiudongzichanzhouzhuanglv, "yy_liudongzichanzhouzhuanglv")
	ability_set16 = d.ParseStockPerTicket(post.Query.YyAbility.YyZongzichanzhouzhuanglv, "yy_zongzichanzhouzhuanglv")
	ability_set17 = d.ParseStockPerTicket(post.Query.YyAbility.YyGudongquanyizhouzhuanglv, "yy_gudongquanyizhouzhuanglv")

	// 其他条件
	other_set1 = d.ParseCount(post.Query.Other.SmCount, post.Date, "sm_count")
	other_set2 = d.ParseCount(post.Query.Other.FundCount, post.Date, "fund_count")
	other_set3 = d.ParseCount(post.Query.Other.FhCount, post.Date, "fh_count")
	other_set4 = d.ParseCount(post.Query.Other.SgCount, post.Date, "sg_count")
	other_set5 = d.ParseCount(post.Query.Other.ZzCount, post.Date, "zz_count")
	other_set6 = d.ParseCount(post.Query.Other.PgCount, post.Date, "pg_count")
	other_set7 = d.ParseCount(post.Query.Other.ZfCount, post.Date, "zf_count")

	all_sets := []set.Interface{belong_set, location_set, concept_set, form_set,
		per_ticket_set1, per_ticket_set2, per_ticket_set3, per_ticket_set4, per_ticket_set5, per_ticket_set6,
		last_day_set1, last_day_set2, last_day_set3, last_day_set4, last_day_set5,
		other_set1, other_set2, other_set3, other_set4, other_set5, other_set6, other_set7,
		ability_set1, ability_set2, ability_set3, ability_set4, ability_set5, ability_set6, ability_set7, ability_set8,
		ability_set9, ability_set10, ability_set11, ability_set12, ability_set13, ability_set14, ability_set15, ability_set16, ability_set17}

	used_sets := []set.Interface{}
	for _, i := range all_sets {
		if i == nil {
			d.Response(c, map[string]interface{}{"result": nil, "total": 0}, nil)
			return
		}
		if i.Size() > 0 {
			used_sets = append(used_sets, i)
		}
	}
	var coders []interface{}
	if len(used_sets) == 0 {
		coders = nil
	} else if len(used_sets) == 1 {
		coders = used_sets[0].List()
	} else if len(used_sets) == 2 {
		coders = set.Intersection(used_sets[0], used_sets[1]).List()
	} else if len(used_sets) > 2 {
		coders = set.Intersection(used_sets[0], used_sets[1], used_sets[2:]...).List()
	}
	var predicts []dal.Predict
	var total int
	tmp := store.MysqlClient.GetDB().Model(&dal.Predict{}).Where("date = ?", post.Date)
	for _, i := range post.Query.Predicts {
		tmp = tmp.Where("`condition` regexp ? OR `bad_condition` regexp ? OR `finance` regexp ?", i, i, i)
	}
	tmp = tmp.Where("code IN (?)", coders)
	tmp.Count(&total)

	log.Println(fmt.Sprintf("一共筛选(%d个), 带条件后剩余(%d个)", len(coders), total))

	if !utils.ContainsString(OrderLimit, post.Order) {
		post.Order = "id"
		tmp.Order(fmt.Sprintf("%s asc", post.Order))
	} else {
		tmp.Order(fmt.Sprintf("%s desc", post.Order))
	}
	tmp.Limit(limit).Offset(offset).Find(&predicts)

	var response []model.PredictListResponse
	for _, i := range predicts {
		var coder dal.Code
		store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", i.Code).Find(&coder)
		x := model.PredictListResponse{Name: i.Name, Code: i.Code, Price: i.Price, Percent: i.Percent, Location: coder.Location,
			Form: coder.OrganizationalForm, Belong: coder.Belong, FundCount: i.FundCount, SimuCount: i.SMCount, Conditions: i.Condition,
			BadConditions: i.BadCondition, Finance: i.Finance, Date: i.Date, Score: i.Score,
			FenghongCount: i.FenghongCount, SongguCount: i.SongguCount, ZhuangzengCount: i.ZhuangzengCount, PeiguCount: i.PeiguCount, ZengfaCount: i.ZengfaCount, SubcompCount: i.SubcompCount}
		response = append(response, x)
	}
	d.Response(c, map[string]interface{}{"result": response, "total": total}, nil)
}

func (d *PredictControl) GetDetail(c *gin.Context) {
	date := c.DefaultQuery("date", "")
	code := c.DefaultQuery("code", "")
	if code == "" {
		d.Response(c, nil, errors.New("证券代码空"))
		return
	}
	var coder_obj dal.Code
	err := store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ? or name = ?", code, code).Find(&coder_obj).Error
	if err != nil {
		d.Response(c, nil, errors.New("证券代码不存在"))
		return
	} else {
		code = coder_obj.Code
	}
	if date == "" {
		var x []PredictDate
		store.MysqlClient.GetDB().Model(&dal.Predict{}).Select("distinct(date) as date").Order("date desc").Scan(&x)
		if len(x) > 0 {
			date = x[0].Date
		}
	}
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	err = d.doQueryLeft(authentication)
	if err != nil {
		d.Response(c, nil, err)
		return
	}

	var TicketHistoryTmp, TicketHistory []dal.TicketHistory
	store.MysqlClient.GetDB().Model(&dal.TicketHistory{}).Where("code = ? and date <= ?", code, date).Limit(90).Order("date desc").Find(&TicketHistoryTmp)
	for i := len(TicketHistoryTmp) - 1; i >= 0; i-- {
		TicketHistory = append(TicketHistory, TicketHistoryTmp[i])
	}
	//去掉周线数据
	//var TicketHistoryWeekly []dal.TicketHistoryWeekly
	//store.MysqlClient.GetDB().Model(&dal.TicketHistoryWeekly{}).Where("code = ? and date <= ?", code, date).Limit(40).Order("date asc").Find(&TicketHistoryWeekly)

	var Stockholder []dal.Stockholder
	store.MysqlClient.GetDB().Model(&dal.Stockholder{}).Where("code = ?", code).Find(&Stockholder)

	var Stock dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("code = ?", code).Find(&Stock)

	var Predict dal.Predict
	store.MysqlClient.GetDB().Model(&dal.Predict{}).Where("code = ? and date = ?", code, date).Find(&Predict)

	var StockCashFlow []dal.StockCashFlow
	store.MysqlClient.GetDB().Model(&dal.StockCashFlow{}).Where("code = ?", code).Order("date asc").Find(&StockCashFlow)

	var StockLiabilities []dal.StockLiabilities
	store.MysqlClient.GetDB().Model(&dal.StockLiabilities{}).Where("code = ?", code).Order("date asc").Find(&StockLiabilities)

	var StockProfit []dal.StockProfit
	store.MysqlClient.GetDB().Model(&dal.StockProfit{}).Where("code = ?", code).Order("date asc").Find(&StockProfit)

	var PerTickets dal.StockPerTicket
	store.MysqlClient.GetDB().Model(&dal.StockPerTicket{}).Where("code = ?", code).Find(&PerTickets)

	var response = map[string][]model.Signal{}
	response["mg"] = append(response["mg"], model.Signal{"每股未分配利润(元)", PerTickets.Weifenpeilirun})
	response["mg"] = append(response["mg"], model.Signal{"每股资本公积金(元)", PerTickets.Gubengongjijin})
	response["mg"] = append(response["mg"], model.Signal{"每股经营性现金流(元)", PerTickets.Jingyingxianjinliu})
	response["mg"] = append(response["mg"], model.Signal{"调整后每股净资产(元)", PerTickets.Jinzichanafter})
	response["mg"] = append(response["mg"], model.Signal{"每股加权收益(元)", PerTickets.Jiaquanshouyi})
	response["mg"] = append(response["mg"], model.Signal{"调整后每股收益(元)", PerTickets.Shouyiafter})

	response["yy"] = append(response["yy"], model.Signal{"股东权益周转率(次)", PerTickets.YyGudongquanyizhouzhuanglv})
	response["yy"] = append(response["yy"], model.Signal{"总资产周转率(次)", PerTickets.YyZongzichanzhouzhuanglv})
	response["yy"] = append(response["yy"], model.Signal{"流动资产周转率(次)", PerTickets.YyLiudongzichanzhouzhuanglv})
	response["yy"] = append(response["yy"], model.Signal{"存货周转率(次)", PerTickets.YyCunhuozhouzhuanglv})
	response["yy"] = append(response["yy"], model.Signal{"应收账款周转率(次)", PerTickets.YyYingshouzhangkuanzhouzhuanlv})

	response["cz"] = append(response["cz"], model.Signal{"总资产增长率(%)", PerTickets.CzZongzichanzengzhanglv})
	response["cz"] = append(response["cz"], model.Signal{"净资产增长率(%)", PerTickets.CzJingzichanzengzhanglv})
	response["cz"] = append(response["cz"], model.Signal{"净利润增长率(%)", PerTickets.CzJinglirunzengzhanglv})
	response["cz"] = append(response["cz"], model.Signal{"主营业务收入增长率(%)", PerTickets.CzZhuyingyewushouruzengzhanglv})

	response["yl"] = append(response["yl"], model.Signal{"资产报酬率(%))", PerTickets.YlZichanbaochoulv})
	response["yl"] = append(response["yl"], model.Signal{"净资产报酬率(%)", PerTickets.YlJingzichanbaochoulv})
	response["yl"] = append(response["yl"], model.Signal{"股本报酬率(%)", PerTickets.YlGubenbaochoulv})
	response["yl"] = append(response["yl"], model.Signal{"销售净利率(%)", PerTickets.YlXiaoshoujinglilv})
	response["yl"] = append(response["yl"], model.Signal{"营业利润率(%)", PerTickets.YlYingyelirunlv})
	response["yl"] = append(response["yl"], model.Signal{"总资产净利润率(%)", PerTickets.YlZongzichanjinglirunlv})
	response["yl"] = append(response["yl"], model.Signal{" 主营业务利润率(%)", PerTickets.YlZhuyingyewulirunlv})
	response["yl"] = append(response["yl"], model.Signal{"总资产利润率(%)", PerTickets.YlZongzichanlirunlv})

	var _response model.StockDetail
	_response.TicketHistory = TicketHistory
	_response.Stockholder = Stockholder
	_response.Stock = Stock
	_response.Predict = Predict
	_response.StockCashFlow = StockCashFlow
	_response.StockLiabilities = StockLiabilities
	_response.StockProfit = StockProfit
	//response.TicketHistoryWeekly = TicketHistoryWeekly
	_response.PerTicket = response
	d.Response(c, _response, nil)
}

func (d *PredictControl) GetFunds(c *gin.Context) {
	offset, limit := check.ParamParse.GetPagination(c)
	code := c.DefaultQuery("code", "")
	if code == "" {
		d.Response(c, nil, errors.New("证券代码为空"))
		return
	}
	var StockFund []dal.StockFund
	var total int
	tmp := store.MysqlClient.GetDB().Model(&dal.StockFund{}).Where("code = ?", code)
	tmp.Count(&total)
	tmp.Offset(offset).Limit(limit).Order("percent_jingzhi desc").Find(&StockFund)
	d.Response(c, StockFund, nil, total)
}

func (d *PredictControl) FundHold(c *gin.Context) {
	offset, limit := check.ParamParse.GetPagination(c)
	code := c.DefaultQuery("fund_code", "")
	if code == "" {
		d.Response(c, nil, errors.New("机构代码为空"))
		return
	}
	var StockFund []dal.StockFund
	var total int
	tmp := store.MysqlClient.GetDB().Model(&dal.StockFund{}).Where("fund_code = ?", code)
	tmp.Count(&total)
	tmp.Offset(offset).Limit(limit).Order("percent_jingzhi desc").Find(&StockFund)
	d.Response(c, StockFund, nil, total)
}

func (d *PredictControl) TopHolderHold(c *gin.Context) {
	holder := c.DefaultQuery("holder_name", "")
	if holder == "" {
		d.Response(c, nil, errors.New("查询用户为空"))
		return
	}
	var Stockholder []dal.Stockholder
	store.MysqlClient.GetDB().Model(&dal.Stockholder{}).Where("holder_name = ?", holder).Find(&Stockholder)
	d.Response(c, Stockholder, nil)
}

type PredictDate struct {
	Date string `json:"date"`
}

type Belongs struct {
	Belong string `json:"date"`
}

type Locations struct {
	Id       string `json:"id"`
	Location string `json:"location"`
}

type OrganizationalForms struct {
	OrganizationalForm string `json:"organizational_form"`
}

type Concepts struct {
	Name string `json:"name"`
}

func (d *PredictControl) GetPredictDates(c *gin.Context) {
	var x []PredictDate
	store.MysqlClient.GetDB().Model(&dal.Predict{}).Select("distinct(date) as date").Order("date desc").Scan(&x)
	d.Response(c, x, nil)
}

func (d *PredictControl) GetConditions(c *gin.Context) {
	var x []dal.Conditions
	response := map[string][]string{}
	store.MysqlClient.GetDB().Model(&dal.Conditions{}).Find(&x)
	for _, i := range x {
		response[i.Type] = append(response[i.Type], i.Name)
	}
	d.Response(c, response, nil)
}

func (d *PredictControl) GetHighConditions(c *gin.Context) {
	var x []dal.HighConditions
	response := map[string][]dal.HighConditions{}
	store.MysqlClient.GetDB().Model(&dal.HighConditions{}).Find(&x)
	for _, i := range x {
		response[i.Tag] = append(response[i.Tag], i)
	}
	d.Response(c, response, nil)
}

func (d *PredictControl) GetQueryList(c *gin.Context) {
	response := map[string][]string{"concept": ConceptList, "label": LabelList,
		"belong": BelongList, "location": LocationList, "forms": FormList}
	d.Response(c, response, nil)
}

type FindRes struct {
	Count float64 `json:"count"`
	Date  string  `json:"date"`
}

func (d *PredictControl) GetFenHong(c *gin.Context) {
	typ := c.DefaultQuery("type", "")
	code := c.DefaultQuery("code", "")
	if code == "" || typ == "" {
		d.Response(c, nil, errors.New("code/type 为空"))
		return
	}
	var response []FindRes
	var err error
	switch typ {
	case "fh":
		err = store.MysqlClient.GetDB().Model(&dal.StockFengHong{}).Select("pai_xi as count, date").Where("code = ? and pai_xi > 0", code).Scan(&response).Error
	case "zz":
		err = store.MysqlClient.GetDB().Model(&dal.StockFengHong{}).Select("zhuang_zeng as count, date").Where("code = ? and zhuang_zeng > 0", code).Scan(&response).Error
	case "sg":
		err = store.MysqlClient.GetDB().Model(&dal.StockFengHong{}).Select("song_gu as count, date").Where("code = ? and song_gu > 0", code).Scan(&response).Error
	}
	log.Println(err)
	d.Response(c, response, nil)
}

func (d *PredictControl) GetPeiGuZhuangZeng(c *gin.Context) {
	typ := c.DefaultQuery("type", "")
	code := c.DefaultQuery("code", "")
	if code == "" || typ == "" {
		d.Response(c, nil, errors.New("code/type 为空"))
		return
	}
	switch typ {
	case "pg":
		var peigu []dal.StockPeiGu
		store.MysqlClient.GetDB().Model(&dal.StockPeiGu{}).Where("code = ?", code).Find(&peigu)
		d.Response(c, peigu, nil)
		return
	case "zf":
		var peigu []dal.StockZengFa
		store.MysqlClient.GetDB().Model(&dal.StockZengFa{}).Where("code = ?", code).Find(&peigu)
		d.Response(c, peigu, nil)

		return
	}
	d.Response(c, nil, errors.New("type 类型错误"))
	return
}

func (d *PredictControl) GetSubComp(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	offset, limit := check.ParamParse.GetPagination(c)
	if code == "" {
		d.Response(c, nil, errors.New("code 为空"))
		return
	}
	var subs []dal.StockSubCompany
	store.MysqlClient.GetDB().Model(&dal.StockSubCompany{}).Where("code = ?", code).Offset(offset).Limit(limit).Find(&subs)
	d.Response(c, subs, nil)
}
