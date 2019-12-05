// @Time:       2019/11/30 上午10:45

package model

// jsapi 支付response （唤起微信支付）
type WeResJsApi struct {
	TimeStamp string `json:"timeStamp"` // 时间戳
	NonceStr  string `json:"nonceStr"`  // 随机字符串
	Package   string `json:"package"`   // PrepayId 拼接的字符串
	Sign      string `json:"sign"`      // 加密签名
	SignType  string `json:"signType"`
	AppId     string `json:"appId"`
}

// 微信登录
type WxUserInfo struct {
	OpenId     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	City       string `json:"city"`
	Province   string `json:"province"`
	Country    string `json:"country"`
	Headimgurl string `json:"headimgurl"`
}
