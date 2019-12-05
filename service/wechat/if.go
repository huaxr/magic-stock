// @Time:       2019/12/1 上午11:18

package wechat

import (
	"regexp"

	"github.com/panghu1024/anypay"
)

type WechatServiceIF interface {
	GetAccessTokenByCode(code string) (openid, access_token string)
	GetCodeUrl() string
	JSApiPay(openid string, money string) *anypay.WeResJsApi
}

const (
	IP                     = "62.234.65.214"
	PAY_BODY               = "天津市知修科技有限公司"
	WX_APPID               = "wx882db180ce1b9351"
	WX_APPSECRET           = "a23e4b2ac82832f0b5ade8ab80cfcd91"
	WX_MCH                 = "1546048431"
	WX_KEY                 = "11111222223333344444aaaaabbbbb33" //商户KEY
	GrantType              = "authorization_code"
	OpenIdUrl              = "https://api.weixin.qq.com/sns/oauth2/access_token?code=%s&appid=%s&secret=%s&grant_type=%s"
	UserInfoUrl            = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s"
	CodeUrl                = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect`
	PayCallbackXmlResponse = `<xml>
          <return_code><![CDATA[SUCCESS]]></return_code>
          <return_msg><![CDATA[OK]]></return_msg>
        </xml>`
)

var (
	PayCallbackXmlCompile = regexp.MustCompile(`<result_code><!\[CDATA\[(.*?)]]></result_code>`)
)

var WechatGlobal WechatServiceIF

func init() {
	tmp := new(WeChat)
	WechatGlobal = tmp
}

type WeChat struct {
}
