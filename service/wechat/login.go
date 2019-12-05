// @Time:       2019/11/28 下午4:03

package wechat

import (
	"encoding/json"
	"fmt"
	"magic/stock/service/check"
	"magic/stock/service/conf"
	"net/url"
)

func (w *WeChat) GetAccessTokenByCode(code string) (string, string) {
	res := check.Authentication.HttpGetWithToken(fmt.Sprintf(OpenIdUrl, code, WX_APPID, WX_APPSECRET, GrantType), "")
	var response map[string]interface{}
	json.Unmarshal(res, &response)
	openid, access_token := response["openid"], response["access_token"]
	return openid.(string), access_token.(string)
}

func (w *WeChat) GetCodeUrl() string {
	return fmt.Sprintf(CodeUrl, WX_APPID, url.QueryEscape(conf.Config.Host))
}
