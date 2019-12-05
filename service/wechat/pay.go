// @Time:       2019/11/28 下午4:03

package wechat

import (
	"magic/stock/service/conf"
	"magic/stock/service/encrypt"

	"github.com/panghu1024/anypay"
	uuid "github.com/satori/go.uuid"
)

func (w *WeChat) genOutTrade() string {
	nonce_str := uuid.Must(uuid.NewV4()).String()
	token, _ := encrypt.MD5Client.Encrypt(nonce_str)
	return token
}

// https://www.ctolib.com/panghu1024-anypay.html
func (w *WeChat) JSApiPay(openid string, money string) *anypay.WeResJsApi {
	nonce_str := w.genOutTrade()
	config := anypay.WeConfig{
		AppId: WX_APPID,
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
			return &param
		}
	}
	return nil
}
