// @Time:       2019/11/28 下午3:35
// 微信果园api
package control

import (
	"errors"
	"io/ioutil"
	"log"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/service/adapter"
	sessions "magic/stock/service/middleware/session"
	"magic/stock/service/wechat"

	"github.com/gin-gonic/gin"
)

type UserIF interface {
	Query(where string, args []interface{}) (*dal.User, error)
	Exist(where string, args []interface{}) bool
	GetUserInfo(c *gin.Context)
	LoginByWeChat(c *gin.Context)
	// 充值会员
	PayByWeChat(c *gin.Context)
	TradeCallBack(c *gin.Context)
	Response(c *gin.Context, data interface{}, err error)
}

var UserControlGlobal UserIF

func init() {
	tmp := new(UserControl)
	tmp.service = adapter.UserServiceGlobal
	tmp.response = new(model.HttpResponse)
	UserControlGlobal = tmp
}

type UserControl struct {
	service  adapter.UserServiceIF
	response *model.HttpResponse
}

func (u *UserControl) Query(where string, args []interface{}) (*dal.User, error) {
	return u.service.Query(where, args)
}

func (u *UserControl) Exist(where string, args []interface{}) bool {
	c, _ := u.service.Count(where, args)
	return c > 0
}

func (d *UserControl) GetUserInfo(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	user, err := d.Query("user_id = ?", []interface{}{authentication.Uid})
	d.Response(c, user, err)
}

func (d *UserControl) LoginByWeChat(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		d.Response(c, nil, errors.New("code为空"))
	}
	user, err := adapter.UserServiceGlobal.LoginWx(code)
	if err != nil {
		d.Response(c, nil, err)
	}
	session := sessions.Default(c)
	session.Set("user", user.UserName)
	session.Set("uid", int(user.ID))
	session.Save()
	d.Response(c, "登录成功", nil)
}

func (d *UserControl) PayByWeChat(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	// 30元
	res, err := adapter.UserServiceGlobal.PayWx(authentication)
	if err != nil {
		d.Response(c, nil, err)
	}
	response := model.WeResJsApi{TimeStamp: res.TimeStamp, NonceStr: res.NonceStr, Package: res.Package, Sign: res.Sign, SignType: "MD5", AppId: wechat.WX_APPID}
	d.Response(c, response, nil)
}

func (d *UserControl) TradeCallBack(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	res := wechat.PayCallbackXmlCompile.FindStringSubmatch(string(body))
	if res[1] == "SUCCESS" {
		order_id := c.Param("order_id")
		adapter.PayServiceGlobal.UpdatePaySuccessAndGenerateIndent(order_id)
		c.Data(200, "text/xml", []byte(wechat.PayCallbackXmlResponse))
		return
	} else {
		log.Println("微信回调错误")
	}
}

func (d *UserControl) Response(c *gin.Context, data interface{}, err error) {
	c.AbortWithStatusJSON(200, d.response.Response(data, err))
}
