// @Time:       2019/11/28 下午4:03

package wechat

import (
	"fmt"
	"io/ioutil"
	"log"
	"magic/stock/service/conf"
	"magic/stock/service/encrypt"
	"net/http"
	"strings"

	"github.com/panghu1024/anypay"
	uuid "github.com/satori/go.uuid"
)

func (w *WeChat) genOutTrade() string {
	nonce_str := uuid.Must(uuid.NewV4()).String()
	token, _ := encrypt.MD5Client.Encrypt(nonce_str)
	return token
}

// https://www.ctolib.com/panghu1024-anypay.html
func (w *WeChat) JSApiPay(openid string, money string) (*anypay.WeResJsApi, string) {
	nonce_str := w.genOutTrade()
	log.Println("下单随机值", nonce_str)
	config := anypay.WeConfig{
		AppId: STOCK_WX_APPID,
		MchId: WX_MCH,
		Key:   WX_KEY,
	}
	payment := anypay.NewWePay(config)               //创建实例
	res := payment.UnifiedOrder(anypay.WeOrderParam{ //创建订单
		Body:           PAY_BODY,
		OutTradeNo:     nonce_str,
		TotalFee:       money, // 单位分 字符串
		SpbillCreateIp: IP,    //务必替换成获取的真实IP
		NotifyUrl:      conf.Config.Host + "/api/callback/" + nonce_str,
		TradeType:      "JSAPI",
		Openid:         openid, // JSAPI方式此参数必传
	})

	//结果判断
	if res.Status == 1 { //调用成功
		order := res.Data.(anypay.WeResOrder)
		//生成前端支付参数
		resParam := payment.JsApiParam(order.PrepayId)
		if resParam.Status == 1 {
			param := resParam.Data.(anypay.WeResJsApi)
			return &param, nonce_str
		}
	}
	return nil, ""
}

func (w *WeChat) H5Pay(ip string) {
	nonce_str := uuid.Must(uuid.NewV4()).String()
	out_trade_no, _ := encrypt.MD5Client.Encrypt(nonce_str)
	notify_url := "https://stock.zhixiutec.com/api/callback/" + out_trade_no
	// 回调函数
	scene_info := `{"h5_info": {"type": "Wap", "wap_url": "https://www.payme.com/api/h5_pay", "wap_name": "xxx"}}`
	signA := fmt.Sprintf("appid=%s&body=%s&mch_id=%s&nonce_str=%s&notify_url=%s&out_trade_no=%s&scene_info=%s&spbill_create_ip=%s&total_fee=%s&trade_type=MWEB",
		STOCK_WX_APPID, "知修科技", WX_MCH, nonce_str, notify_url, out_trade_no, scene_info, ip, "1")
	strSignTmp := signA + "&key=" + WX_KEY
	token, _ := encrypt.MD5Client.Encrypt(strSignTmp)
	sign := strings.ToUpper(token)
	path := signA + "&sign=" + sign
	post_data := "<xml>"
	for _, i := range strings.Split(path, "&") {
		xml1, xml2 := strings.Split(i, "=")[0], strings.Split(i, "=")[1]
		post_data = post_data + "<" + xml1 + ">" + xml2 + "</" + xml1 + ">"
	}
	post_data = post_data + "</xml>"
	log.Println("H5支付Post data", post_data)
	//headers = {'Content-Type': 'binary'}
	//	# 解决post_data 中文编码问题
	//	url = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	//	res = requests.post(url, data=post_data.encode(), headers=headers, verify=False)
	//	# 提交订单信息
	//	# res.text.encode('utf-8')
	//	print(res.text.encode('latin_1').decode('utf8'))
	//	pattern = re.compile("<mweb_url><!\[CDATA\[(.*?)]]></mweb_url")
	//
	//	redicrt_url = pattern.findall(res.text)[0]
	//	# 匹配微信回调函数，调用微信app进行支付
	//	# self.redirect(redicrt_url)
	//	print("the url is url", redicrt_url)
	//
	//	return redicrt_url, out_trade_no
	req, _ := http.NewRequest("GET", "https://api.mch.weixin.qq.com/pay/unifiedorder", nil)
	req.Header.Add("Content-Type", "binary")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("H5 微信端返回", string(body))
	xx := H5PayCompile.FindStringSubmatch(string(body))
	log.Println(xx)
}
