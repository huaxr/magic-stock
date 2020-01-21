// @Time:       2019/11/28 下午4:03

package wechat

import (
	"bytes"
	"errors"
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
	log.Println("支付金额:", money)
	nonce_str := w.genOutTrade() // 下单随机值
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
	} else {
		log.Println("调用失败", res)
	}
	return nil, ""
}

func (w *WeChat) H5Pay(ip string) (string, error) {
	nonce_str := strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	out_trade_no, _ := encrypt.MD5Client.Encrypt(nonce_str)
	notify_url := "https://stock.zhixiutec.com/api/callback/" + out_trade_no
	// 回调函数
	scene_info := `{"h5_info": {"type": "Wap", "wap_url": "https://stock.zhixiutec.com/api/h5_pay", "wap_name": "xxx"}}`
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
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", bytes.NewBuffer([]byte(post_data)))
	if err != nil {
		fmt.Println("出现错误1", err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("出现错误", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	xx := H5PayCompile.FindStringSubmatch(string(body))
	if len(xx) == 2 {
		log.Println("回调url为:", xx[1])
		return xx[1], nil
	}
	log.Println("失败", xx)
	return "", errors.New("失败")
}
