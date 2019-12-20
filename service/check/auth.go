// @Contact:    huaxinrui
// @Time:       2019/9/27 上午11:52

package check

import (
	"io/ioutil"
	"log"
	"magic/stock/dao"
	"magic/stock/model"
	"magic/stock/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var Authentication AuthenticationIF

func init() {
	tmp := new(authentication)
	if utils.TellEnv() == "loc" {
		tmp.debug = true
	}
	tmp.sync_group_chan = make(chan map[string]interface{}, 1)
	Authentication = tmp
}

type authentication struct {
	debug           bool // 本地环境
	sync_group_chan chan map[string]interface{}
}

func (a *authentication) HttpGetWithToken(url, token string) []byte {
	err := Security.PreventSSRF(url)
	if err != nil {
		return nil
	}
	req, _ := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Add("Authorization", token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil
	}
	return body
}

func (a *authentication) HttpGetWithBasicAuth(url string, key, sec string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(GET, url, nil)
	req.SetBasicAuth(key, sec)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	bodyByte, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodyByte
}

func (a *authentication) GetHeaderToken(c *gin.Context, head string) (string, bool) {
	h := c.Request.Header.Get(head)
	if len(h) > 0 {
		return strings.TrimLeft(h, "Token "), true
	}
	return "", false
}

// Token 认证
func (a *authentication) checkToken(c *gin.Context) (*model.AuthResult, bool) {
	return nil, false
}

// Session 认证
func (a *authentication) checkSession(c *gin.Context) *model.AuthResult {
	user_obj, _ := dao.UserDao.Query("id = ?", []interface{}{888})
	return &model.AuthResult{User: user_obj.UserName, Uid: int(user_obj.ID), Member: user_obj.IsMember, QueryLeft: user_obj.QueryLeft}

	//session := sessions.Default(c)
	//user := session.Get("user")
	//uid := session.Get("uid")
	//if user == nil || uid == nil {
	//	if utils.TellEnv() == "loc" {
	//		user_obj, _ := dao.UserDao.Query("id = ?", []interface{}{888})
	//		return &model.AuthResult{User: user_obj.UserName, Uid: int(user_obj.ID), Member: user_obj.IsMember, QueryLeft: user_obj.QueryLeft}
	//	} else {
	//		return &model.AuthResult{errors.New("登录错误"), "", -1, false, 0}
	//	}
	//}
	//user_obj, _ := dao.UserDao.Query("id = ?", []interface{}{uid})
	//return &model.AuthResult{User: user.(string), Uid: uid.(int), Member: user_obj.IsMember, QueryLeft: user_obj.QueryLeft}
}

func (a *authentication) getDebug() bool {
	return a.debug
}

// token authentication takes precedence over session
func (a *authentication) JudgeApi(c *gin.Context) *model.AuthResult {
	return a.checkSession(c)
}
