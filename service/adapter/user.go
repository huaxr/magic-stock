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

	"github.com/panghu1024/anypay"
)

type UserServiceIF interface {
	Create(event *dal.User) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*dal.User, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]dal.User, error)
	Count(where string, args []interface{}) (int, error)
	CreateUserIfNotExist(user *dal.User) (us *dal.User, err error)
	LoginWx(code string) (*dal.User, error)
	PayWx(authentication *model.AuthResult) (*anypay.WeResJsApi, error)
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

func (m *UserService) QueryAll(where string, args []interface{}, offset, limit int) (*[]dal.User, error) {
	return m.dao.QueryAll(where, args, offset, limit)
}

func (m *UserService) Count(where string, args []interface{}) (int, error) {
	return m.dao.Count(where, args)
}

func (u *UserService) CreateUserIfNotExist(user *dal.User) (us *dal.User, err error) {
	user_obj, err := u.Query("user_name = ?", []interface{}{user.UserName})
	if err != nil {
		store.MysqlClient.GetDB().Save(user)
		user_obj = user
	}
	return user_obj, err
}

func (u *UserService) LoginWx(code string) (*dal.User, error) {
	login_response, err := wechat.WechatGlobal.GetAccessTokenByCode(code)
	if err != nil {
		return nil, err
	}
	res := check.Authentication.HttpGetWithToken(fmt.Sprintf(wechat.UserInfoUrl, login_response.AccessToken, login_response.Openid), "")
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
	user := dal.User{OpenId: openid, UserName: username, Avatar: avatar, Sex: sex, City: city, Province: province, Country: country, IsMember: false}
	log.Println(user)
	obj, _ := u.CreateUserIfNotExist(&user)
	return obj, nil
}

func (u *UserService) PayWx(authentication *model.AuthResult) (*anypay.WeResJsApi, error) {
	user, _ := u.Query("id = ?", []interface{}{authentication.Uid})
	payment := wechat.WechatGlobal.JSApiPay(user.OpenId, strconv.Itoa(int(1)))
	if payment == nil {
		return nil, errors.New("唤起支付调用失败")
	}
	pay_record := dal.Pay{UserId: authentication.Uid, Spend: 30, PaySuccess: false, OrderId: payment.NonceStr}
	store.MysqlClient.GetDB().Save(&pay_record)
	return payment, nil
}

func (m *UserService) genName() string {
	return "自定义条件-" + time.Now().Format("20060102150405")
}

func (u *UserService) SaveUserConditions(post *model.GetPredicts, authentication *model.AuthResult) error {
	res, _ := json.Marshal(post.Query)
	uc := dal.UserConditions{Name: u.genName(), UserId: authentication.Uid, Conditions: res}
	err := store.MysqlClient.GetDB().Save(&uc).Error
	return err
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
