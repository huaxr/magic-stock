// @Time:       2019/11/28 下午3:35
package control

import (
	"errors"
	"io/ioutil"
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"magic/stock/service/adapter"
	"magic/stock/service/check"
	"magic/stock/service/conf"
	sessions "magic/stock/service/middleware/session"
	"magic/stock/service/wechat"
	"time"

	"github.com/gin-gonic/gin"
)

type UserIF interface {
	Query(where string, args []interface{}) (*dal.User, error)
	Exist(where string, args []interface{}) bool
	GetUserInfo(c *gin.Context)
	GetUserToken(c *gin.Context)
	JudgeIsMember(c *gin.Context)
	LoginByWeChat(c *gin.Context)
	LogOut(c *gin.Context)
	// 充值 (h5 和 jsapi)
	PayByWeChatJsApi(c *gin.Context)
	PayByWeChatH5(c *gin.Context)
	// 分享
	GetWxSign(c *gin.Context)
	TradeCallBack(c *gin.Context)

	GetConditions(c *gin.Context)
	GetInvite(c *gin.Context)
	GetDemands(c *gin.Context)
	EditUserConditions(c *gin.Context)
	DeleteUserConditions(c *gin.Context)
	// 提需求
	SubmitDemand(c *gin.Context)
	AddStock(c *gin.Context)
	MySelect(c *gin.Context)
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
	user, err := d.Query("id = ?", []interface{}{authentication.Uid})
	d.Response(c, user, err)
}

func (d *UserControl) GetUserToken(c *gin.Context) {
	authentication := check.Authentication.JudgeApi(c)
	if authentication.Err != nil {
		d.Response(c, "", nil)
		return
	}
	user, _ := d.Query("id = ?", []interface{}{authentication.Uid})
	d.Response(c, user.ShareToken, nil)
	return
}

func (d *UserControl) JudgeIsMember(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	user, err := d.Query("id = ?", []interface{}{authentication.Uid})
	if user.MemberExpireTime.After(time.Now()) {
		d.Response(c, true, err)
		return
	} else {
		d.Response(c, false, err)
	}
}

func (d *UserControl) LoginByWeChat(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	token := c.DefaultQuery("token", "")
	if code == "" {
		d.Response(c, nil, errors.New("code为空"))
		return
	}
	user, err := adapter.UserServiceGlobal.LoginWx(code, token)
	if err != nil {
		d.Response(c, nil, err)
		return
	}
	if user.OpenId != "" {
		session := sessions.Default(c)
		session.Set("user", user.UserName)
		session.Set("open_id", user.OpenId)
		session.Set("uid", int(user.ID))
		session.Save()
		log.Println("登录成功")
		c.Redirect(302, conf.Config.Host)
		return
	} else {
		d.Response(c, nil, errors.New("未知错误"))
	}
}

func (d *UserControl) PayByWeChatJsApi(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)

	var post model.SpendType
	err := c.BindJSON(&post)

	log.Println("充值post请求:", post)

	res, err := adapter.UserServiceGlobal.PayWxJsAPi(authentication, &post)
	if err != nil {
		d.Response(c, nil, err)
		return
	}
	response := model.WeResJsApi{TimeStamp: res.TimeStamp, NonceStr: res.NonceStr, Package: res.Package, Sign: res.Sign, SignType: "MD5", AppId: wechat.STOCK_WX_APPID}
	d.Response(c, response, nil)
}

func (d *UserControl) PayByWeChatH5(c *gin.Context) {
	web_url, err := adapter.UserServiceGlobal.PayWxH5(c)
	if err != nil {
		d.Response(c, nil, err)
		return
	}
	c.JSON(200, gin.H{"error_code": 0, "err_msg": "", "url": web_url})
}

func (d *UserControl) TradeCallBack(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	res := wechat.PayCallbackXmlCompile.FindStringSubmatch(string(body))
	log.Println("支付回调", res, res[1])
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

func (d *UserControl) EditUserConditions(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	var post model.EditPredicts
	err := c.BindJSON(&post)
	if err != nil {
		d.Response(c, nil, err)
		return
	}

	err = d.service.EditUserConditions(&post, authentication)
	if err != nil {
		d.Response(c, nil, err)
		return
	}
	d.Response(c, "success", err)
}

func (d *UserControl) DeleteUserConditions(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	var post model.DeletePredicts
	err := c.BindJSON(&post)
	if err != nil {
		d.Response(c, nil, err)
		return
	}
	d.service.DeleteUserConditions(post.Id, authentication)
	d.Response(c, "success", nil)
}

func (d *UserControl) GetConditions(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	var ucs []dal.UserConditions
	store.MysqlClient.GetDB().Model(&dal.UserConditions{}).Where("user_id = ?", authentication.Uid).Find(&ucs)
	d.Response(c, ucs, nil)
}

type Res struct {
	BeShareId int
}

type Users struct {
	UserName  string
	Avatar    string
	CreatedAt time.Time
}

func (d *UserControl) GetInvite(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	offset, limit := check.ParamParse.GetPagination(c)
	var ids []Res
	store.MysqlClient.GetDB().Model(&dal.UserShare{}).Select("be_share_id").Where("share_user_id = ?", authentication.Uid).Offset(offset).Limit(limit).Scan(&ids)

	var tmp []int
	for _, i := range ids {
		tmp = append(tmp, i.BeShareId)
	}
	var users []Users
	store.MysqlClient.GetDB().Model(&dal.User{}).Select("user_name, avatar, created_at").Where("id in (?)", tmp).Scan(&users)
	d.Response(c, users, nil)
}

func (d *UserControl) SubmitDemand(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	var post model.SubmitDemand
	err := c.BindJSON(&post)
	if err != nil {
		d.Response(c, nil, err)
		return
	}
	de := dal.UserDemands{UserId: authentication.Uid, Content: post.Content, Type: post.Type}
	store.MysqlClient.GetDB().Save(&de)
	d.Response(c, "提交成功", nil)
}

func (d *UserControl) GetDemands(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	offset, limit := check.ParamParse.GetPagination(c)
	var demands []dal.UserDemands
	store.MysqlClient.GetDB().Model(&dal.UserDemands{}).Where("user_id = ?", authentication.Uid).Offset(offset).Limit(limit).Find(&demands)
	d.Response(c, demands, nil)
}

func (d *UserControl) LogOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(302, "/")
}

func (d *UserControl) AddStock(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	var post model.AddStock
	err := c.BindJSON(&post)
	if err != nil {
		d.Response(c, nil, err)
		return
	}
	s := dal.UserSelect{UserId: authentication.Uid, Code: post.Code, Name: post.Name, Price: post.Price}
	store.MysqlClient.GetDB().Save(&s)
	d.Response(c, "添加成功", nil)
}

func (d *UserControl) MySelect(c *gin.Context) {
	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)
	var post model.AddStock
	err := c.BindJSON(&post)
	if err != nil {
		d.Response(c, nil, err)
		return
	}
	s := dal.UserSelect{UserId: authentication.Uid, Code: post.Code, Name: post.Name, Price: post.Price}
	store.MysqlClient.GetDB().Save(&s)
	d.Response(c, "添加成功", nil)
}
