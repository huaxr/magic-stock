// @Contact:    huaxinrui
// @Time:       2019/9/27 上午11:52

package check

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"code.byted.org/byte_security/dal/auth"

	"code.byted.org/byte_security/platform_api/models"
	"code.byted.org/byte_security/platform_api/service/conf"
	sessions "code.byted.org/byte_security/platform_api/service/middleware/session"
	"code.byted.org/byte_security/platform_api/service/s_auth"
	"code.byted.org/byte_security/platform_api/utils"
	"code.byted.org/inf/infsecc"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
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
		t, err := infsecc.GetToken(false) // os.Getenv(TokenStringEnv)
		if err == nil {
			req.Header.Add("X-Dps-Token", t) // TCE env enable
		}
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
	req.Header.Add("x-app-version", "v2")
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

func (a *authentication) checkToken(c *gin.Context) (*models.AuthResult, bool) {
	t, ok := a.GetHeaderToken(c, "Authorization")
	if ok {
		t = strings.TrimLeft(t, "Token ")
		token, err := s_auth.TokenServiceGlobal.Query("token = ?", []interface{}{t})
		if err != nil {
			return &models.AuthResult{errors.New("token 无效"), "", "", 0, false, nil}, false
		}
		if c.Request.URL.Path == token.Path {
			user, _ := s_auth.UserServiceGlobal.Query("user_name = ?", []interface{}{token.Owner})
			return &models.AuthResult{nil, token.Owner, "", int(user.ID), true, nil}, true
		} else {
			return &models.AuthResult{errors.New("token 鉴权失败"), "", "", 0, false, nil}, false
		}
	}
	return nil, false
}

func (a *authentication) checkSession(c *gin.Context) *models.AuthResult {
	session := sessions.Default(c)
	user := session.Get("user")
	group := session.Get("group")
	uid := session.Get("uid")
	if user == nil || group == nil || uid == nil {
		//if a.getDebug() {
		//	// whether change the admin flag to the config file .
		//	// cause change the admin status will be frequently when debug it
		//	return &models.AuthResult{nil, "guorui.jerry", "288980000064975441", 2, false, nil}
		//}
		return &models.AuthResult{errors.New("登录错误"), "", "", 0, false, nil}
	}
	u, err := s_auth.CasbinServiceGlobal.Count("p_type = ? and v0 = ? and v1 = ?", []interface{}{"g", user.(string), "admin"})
	if err != nil {
		return &models.AuthResult{errors.New("adapter 异常"), "", "", 0, false, nil}
	}
	return &models.AuthResult{nil, user.(string), group.(string), uid.(int), u > 0, nil}
}

func (a *authentication) GetHeaderToken(c *gin.Context, head string) (string, bool) {
	h := c.Request.Header.Get(head)
	if len(h) > 0 {
		return strings.TrimLeft(h, "Token "), true
	}
	return "", false
}

func (a *authentication) getDebug() bool {
	return a.debug
}

// token authentication takes precedence over session
func (a *authentication) JudgeApi(c *gin.Context) *models.AuthResult {
	res, ok := a.checkToken(c)
	if ok && res != nil {
		return res
	}
	return a.checkSession(c)
}

func (a *authentication) TellAdmin(auth *models.AuthResult) error {
	if auth.Err != nil {
		return errors.New(auth.Err.Error())
	}
	if !auth.Admin {
		return errors.New("没有权限操作")
	}
	return nil
}

func (a *authentication) GetUserNameBySSOTicket(c *gin.Context) (*models.SsoRes, error) {
	//if ticket == "" {
	//	return nil
	//}
	ticket := c.DefaultQuery("ticket", "")
	requests := gorequest.New()
	_, body, _ := requests.Get(conf.Config.SSO + "/cas/serviceValidate?ticket=" + ticket + "&service=" + conf.Config.Host + "/api/auth/login").End()
	//fmt.Println(body)
	if strings.Count(body, "cas:user>") != 2 {
		return nil, errors.New("can't find cas:user>: " + body)
	} else {
		u, _ := regexp.Compile("<cas:user>(.*?)</cas:user>")
		e, _ := regexp.Compile("<cas:employee_id>(.*?)</cas:employee_id>")
		n, _ := regexp.Compile("<cas:name>(.*?)</cas:name>")
		user := u.FindString(body)
		employee_id := e.FindString(body)
		name := n.FindString(body)
		if len(user) < 10 {
			return nil, errors.New(user)
		}
		user = strings.Replace(user, "<cas:user>", "", -1)
		user = strings.Replace(user, "</cas:user>", "", -1)
		if user == "" {
			return nil, errors.New("error user:" + user)
		}
		employee_id = strings.Replace(employee_id, "<cas:employee_id>", "", -1)
		employee_id = strings.Replace(employee_id, "</cas:employee_id>", "", -1)
		name = strings.Replace(name, "<cas:name>", "", -1)
		name = strings.Replace(name, "</cas:name>", "", -1)
		return &models.SsoRes{User: user, EmployeeId: employee_id, Name: name}, nil
	}
}

func (a *authentication) syncGroupLoc(name string) {
	var grouper models.GroupQuery
	client := &http.Client{}
	req, _ := http.NewRequest(GET, fmt.Sprintf(GROUPURL, name), nil)
	req.Header.Add("Authorization", "Basic MTExOjU5QTFEMTlENUJGRjQ2QzdCRkVENDQ2QzBCNUU1RjY2")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(body, &grouper)
	if err != nil {
		panic(err)
	}
	if grouper.Success {
		a.sync_group_chan <- grouper.Employees[0]
	}
}

func (a *authentication) GetGroupByNameLoc(ctx context.Context, name string) map[string]interface{} {
	go a.syncGroupLoc(name)
	select {
	case x := <-a.sync_group_chan:
		return x
	case <-ctx.Done():
		return nil
	}
}

type group_query struct {
	Employees []map[string]interface{}
	Success   bool
}

func (a *authentication) getGroup(name string) group_query {
	token, err := infsecc.GetToken(false)
	if err != nil {
		log.Println(err)
	}
	var grouper group_query

	client := &http.Client{}
	req, _ := http.NewRequest(GET, fmt.Sprintf(GROUPURL, name), nil)
	req.Header.Add("Authorization", "Basic MTExOjU5QTFEMTlENUJGRjQ2QzdCRkVENDQ2QzBCNUU1RjY2")
	req.Header.Add("X-Dps-Token", token)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &grouper)
	if err != nil {
		log.Println(err)
	}
	return grouper
}

// 通过sdk调用获取
func (a *authentication) asyncGroup(name string) {
	grouper := a.getGroup(name)
	if grouper.Success {
		a.sync_group_chan <- grouper.Employees[0]
	}
}

func (a *authentication) GetGroupByNameAsync(ctx context.Context, name string) map[string]interface{} {
	go a.asyncGroup(name)

	select {
	case x := <-a.sync_group_chan:
		return x

	case <-ctx.Done():
		//fmt.Println(ctx.Err()) // prints "context deadline exceeded"
		return nil
	}
}

func (a *authentication) GetLoginUrl() string {
	return conf.Config.SSO + "/cas/login?service=" + conf.Config.Host + "/api/auth/login"
}

func (a *authentication) GetGroupByName(ctx context.Context, name string) map[string]interface{} {
	grouper := a.getGroup(name)
	return grouper.Employees[0]
}

func (a *authentication) GetUserInfo(name string) (*auth.User, string) {
	res := Authentication.HttpGetWithToken(fmt.Sprintf(GROUPURL, name), PEOPLETOKEN)
	//fmt.Println(string(res))
	var response models.Info
	json.Unmarshal(res, &response)
	if response.Success && len(response.Employees) > 0 {
		i := response.Employees[0]
		log.Println(i.Username, i.Name, i.Id, i.Department.Id, i.AvatarUrl, i.Leader.Email)
		var l string
		leader := strings.Split(i.Leader.Email, "@")
		if len(leader) > 0 {
			l = leader[0]
		}
		u := auth.User{UserName: i.Username, RealName: i.Name, Email: i.Email, UserNum: strconv.Itoa(i.Id), GroupId: i.Department.Id, AvatarUrl: i.AvatarUrl, Leader: l}
		return &u, i.Department.Name
	}
	return nil, ""
}

func (a *authentication) GetKaNiUrl(user string) string {
	return "https://ei.byted.org/ratak/employees/" + user + "/permissions/?employee_type=0"
}

func (a *authentication) GetKaNiHasRoleUrl(user_email, role_name string) string {
	return fmt.Sprintf("https://ei.byted.org/ratak/user/%s/role/%s/", user_email, role_name)
}

func (a *authentication) GetKaNi(user string) map[string][]string {
	U := a.GetKaNiUrl(user)
	res := a.HttpGetWithBasicAuth(U, KaNiApp.GetBasicId(), KaNiApp.GetBasicSecret())
	var response map[string][]string
	json.Unmarshal(res, &response)
	return response
}

func (a *authentication) QueryEmployeeHasPerm(user, typ, perm string) bool {
	resource := new(Resource)
	resource.Result = a.GetKaNi(user)
	return resource.HasPerm(typ, perm)
}

func (a *authentication) QueryEmployeeAllPerm(user, typ string) []string {
	response := a.GetKaNi(user)
	result, ok := response[typ]
	if !ok {
		return nil
	}
	return result
}

func (a *authentication) JudgeHasKaNiPerm(typ string, perms []string) bool {
	if utils.ContainsString(perms, "all") {
		return true
	}
	if utils.ContainsString(perms, "admin") {
		return true
	}
	if utils.ContainsString(perms, typ) {
		return true
	}
	return false
}

func (a *authentication) JudgeHasKaNiRole(user_email, role_name string) bool {
	U := a.GetKaNiHasRoleUrl(user_email, role_name)
	res := a.HttpGetWithBasicAuth(U, KaNiApp.GetBasicId(), KaNiApp.GetBasicSecret())
	var response map[string]interface{}
	json.Unmarshal(res, &response)
	return response["ok"].(bool)
}
