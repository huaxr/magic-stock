// @Time:       2019/11/28 下午4:03

package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"magic/stock/model"
	"magic/stock/service/check"
	"magic/stock/service/conf"
	"net/url"
)

func (w *WeChat) GetAccessTokenByCode(code string) (*model.AccessTokenResponse, error) {
	res := check.Authentication.HttpGetWithToken(fmt.Sprintf(OpenIdUrl, code, STOCK_WX_APPID, STOCK_WX_APPSECRET, GrantType), "")
	log.Println("GetAccessTokenByCode", string(res))
	var response model.AccessTokenResponse
	err := json.Unmarshal(res, &response)
	if err != nil {
		return nil, errors.New("GetAccessTokenByCode 登录失败")
	}
	return &response, nil
}

func (w *WeChat) GetCodeUrl() string {
	return fmt.Sprintf(CodeUrl, STOCK_WX_APPID, url.QueryEscape(conf.Config.WxRedirect))
}
