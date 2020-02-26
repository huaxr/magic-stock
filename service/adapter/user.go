// @Contact:    huaxinrui
// @Time:       2019/9/26 下午9:25

package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/dao"
	"magic/stock/model"
	"magic/stock/service/check"
	"magic/stock/service/wechat"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	uuid "github.com/satori/go.uuid"

	"github.com/panghu1024/anypay"
)

type UserServiceIF interface {
	Create(event *dal.User) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*dal.User, error)
	QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.User, error)
	Count(where string, args []interface{}) (int, error)
	CreateUserIfNotExist(user *dal.User, token string) (us *dal.User, err error)
	LoginWx(code, token string) (*dal.User, error)
	PayWxJsAPi(authentication *model.AuthResult, post *model.SpendType) (*anypay.WeResJsApi, error)
	PayWxH5(c *gin.Context) (string, error)
	SaveUserConditions(query *model.GetPredicts, auth *model.AuthResult) error
	EditUserConditions(query *model.EditPredicts, auth *model.AuthResult) error
	DeleteUserConditions(id int, auth *model.AuthResult) error
}

var UserServiceGlobal UserServiceIF

func init() {
	tmp := new(UserService)
	tmp.dao = dao.UserDao
	UserServiceGlobal = tmp
}

type UserService struct {
	dao dao.UserDaoIF
}

func (m *UserService) Create(app *dal.User) error {
	return m.dao.Create(app)
}

func (m *UserService) Delete(id int) error {
	return m.dao.Delete(id)
}

func (m *UserService) Update(id int, ma map[string]interface{}) error {
	return m.dao.Update(id, ma)
}

func (m *UserService) Query(where string, args []interface{}) (*dal.User, error) {
	return m.dao.Query(where, args)
}

func (m *UserService) QueryAll(where string, args []interface{}, offset, limit int, select_only string) (*[]dal.User, error) {
	return m.dao.QueryAll(where, args, offset, limit, select_only)
}

func (m *UserService) Count(where string, args []interface{}) (int, error) {
	return m.dao.Count(where, args)
}

func (u *UserService) CreateUserIfNotExist(user *dal.User, token string) (us *dal.User, err error) {
	user_obj, err := u.Query("open_id = ?", []interface{}{user.OpenId})
	if err != nil {
		store.MysqlClient.GetDB().Save(user)
		user_obj = user

		if token != "" {
			user_obj2, err := u.Query("share_token = ?", []interface{}{token})
			if err == nil {
				user_obj2.QueryLeft += 200
				store.MysqlClient.GetDB().Save(user_obj2)
				// 保存拉新记录
				share_record := dal.UserShare{ShareUserId: int(user_obj.ID), BeShareId: int(user_obj2.ID)}
				store.MysqlClient.GetDB().Model(&share_record)
			} else {
				log.Println("该token不存在", token)
			}
		}
	}
	return user_obj, err
}

func (u *UserService) LoginWx(code, token string) (*dal.User, error) {
	login_response, err := wechat.WechatGlobal.GetAccessTokenByCode(code)
	if err != nil {
		log.Println("GetAccessTokenByCode error:", err)
		return nil, err
	}
	res := check.Authentication.HttpGetWithToken(fmt.Sprintf(wechat.UserInfoUrl, login_response.AccessToken, login_response.Openid), "")
	//log.Println("微信返回:", string(res))
	var user_info model.WxUserInfo
	err = json.Unmarshal(res, &user_info)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}
	avatar := user_info.Headimgurl
	sex := user_info.Sex
	city := user_info.City
	province := user_info.Province
	country := user_info.Country
	username := user_info.Nickname
	openid := user_info.OpenId
	if openid == "" {
		return nil, errors.New("未知错误")
	}
	uid := uuid.Must(uuid.NewV4()).String()
	user := dal.User{OpenId: openid, UserName: username, Avatar: avatar, Sex: sex, City: city, Province: province, Country: country, MemberExpireTime: time.Now(), QueryLeft: 20, ShareToken: uid[0:23]}
	obj, err := u.CreateUserIfNotExist(&user, token)
	return obj, nil
}

func (u *UserService) getMoney(typ int) *dal.Price {
	var price dal.Price
	store.MysqlClient.GetDB().Model(&dal.Price{}).Where("id = ?", typ).Find(&price)
	return &price
}

func (u *UserService) PayWxJsAPi(authentication *model.AuthResult, post *model.SpendType) (*anypay.WeResJsApi, error) {
	price := u.getMoney(post.Id)
	user, _ := u.Query("id = ?", []interface{}{authentication.Uid})
	payment, NonceStr := wechat.WechatGlobal.JSApiPay(user.OpenId, strconv.Itoa(price.Spend))
	if payment == nil {
		return nil, errors.New("唤起支付调用失败")
	}
	pay_record := dal.Pay{UserId: authentication.Uid, Spend: price.Spend, PaySuccess: false, OrderId: NonceStr, Type: price.Type, Count: price.Count, Extra: post.Content}
	store.MysqlClient.GetDB().Save(&pay_record)
	return payment, nil
}

func (u *UserService) PayWxH5(c *gin.Context) (string, error) {
	//_auth, _ := c.Get("auth")
	//authentication := _auth.(*model.AuthResult)
	//user, _ := u.Query("id = ?", []interface{}{authentication.Uid})
	return wechat.WechatGlobal.H5Pay(c.ClientIP())
}

func (m *UserService) genName() string {
	return "自定义条件-" + time.Now().Format("20060102150405")
}

func (u *UserService) SaveUserConditions(post *model.GetPredicts, authentication *model.AuthResult) error {
	var count int
	store.MysqlClient.GetDB().Model(&dal.UserConditions{}).Where("user_id = ?", authentication.Uid).Count(&count)
	if count <= 10 {
		res, _ := json.Marshal(post.Query)
		uc := dal.UserConditions{Name: u.genName(), UserId: authentication.Uid, Conditions: res}
		err := store.MysqlClient.GetDB().Save(&uc).Error
		return err
	} else {
		return errors.New("用户条件多为10个, 改条件将不再保存")
	}
}

func (u *UserService) EditUserConditions(query *model.EditPredicts, auth *model.AuthResult) error {
	var uc dal.UserConditions
	err := store.MysqlClient.GetDB().Model(&dal.UserConditions{}).Where("id = ? and user_id = ?", query.Id, auth.Uid).Find(&uc).Error
	if err != nil {
		return err
	}
	res, _ := json.Marshal(query.Query)
	uc.Name = query.Name
	uc.Conditions = res
	store.MysqlClient.GetDB().Save(&uc)
	return nil
}

func (u *UserService) DeleteUserConditions(id int, auth *model.AuthResult) error {
	store.MysqlClient.GetDB().Delete(&dal.UserConditions{}, "id = ? and user_id = ?", id, auth.Uid)
	return nil
}
