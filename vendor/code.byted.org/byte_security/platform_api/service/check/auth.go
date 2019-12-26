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

	"code.byted.org/byte_security/platform_api/dao/d_system"

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

func (a *authentication) parseQuery(query_str string) map[string][]string {
	if query_str == "" {
		return nil
	}
	mapper := map[string][]string{}
	query_part := strings.TrimSpace(strings.TrimLeft(query_str, "?"))
	quert_parts := strings.Split(query_part, "&")
	for _, i := range quert_parts {
		key_value := strings.Split(i, "=")
		if len(key_value) != 2 {
			return nil
		}
		key := key_value[0]
		value := key_value[1]
		values := strings.Split(value, ",")
		mapper[key] = values
	}
	return mapper
}

func (a *authentication) compareMapper(mapper_must, mapper_received map[string][]string) error {
	if mapper_must == nil || mapper_received == nil {
		return errors.New("该token需要指定query")
	}
	// mapper_must 中的key必须存在
	for k, v := range mapper_must {
		if _, ok := mapper_received[k]; !ok {
			var keys []string
			for i, _ := range mapper_must {
				keys = append(keys, i)
			}
			// path 中的数据必须一致才能认证通过
			return errors.New(fmt.Sprintf("错误的token query key, 必须拥有 %s", strings.Join(keys, ",")))
		}
		value := mapper_received[k]
		for _, i := range value {
			if !utils.ContainsString(v, i) {
				return errors.New(fmt.Sprintf("错误的token query [%s] value, 只能允许 %s", k, strings.Join(v, ",")))
			}
		}
	}
	return nil
}

func (a *authentication) checkToken(c *gin.Context) (*models.AuthResult, bool) {
	t, ok1 := a.GetHeaderToken(c, "Authorization")
	faasURL, ok2 := a.GetHeaderToken(c, "FAASURL")
	if ok1 {
		t = strings.TrimLeft(t, "Token ")
		token, err := s_auth.TokenServiceGlobal.Query("token = ?", []interface{}{t})
		if err != nil {
			return &models.AuthResult{Err: errors.New("token 无效")}, false
		}
		// 如果使用了指定参数鉴权 e.g 漏洞业务类型为ies等 business_id=1,2,3&type=a
		if len(token.MustQuery) > 0 {
			log.Println(token.MustQuery, c.Request.URL.RawQuery)
			err := a.compareMapper(a.parseQuery(token.MustQuery), a.parseQuery(c.Request.URL.RawQuery))
			if err != nil {
				return &models.AuthResult{Err: err}, false
			}
		}
		if c.Request.URL.Path == token.Path {
			userName := token.Owner
			if ok2 { // for different robot token
				robot, _ := d_system.RobotDao.GetRobotInfoByURL(faasURL)
				if len(robot.Owner) > 0 {
					userName = robot.Owner
				}
			}
			user, _ := s_auth.UserServiceGlobal.Query("user_name = ?", []interface{}{userName})
			userEmail := user.Email
			if len(userEmail) == 0 {
				userEmail, _ = s_auth.UserServiceGlobal.GetUserEmailFromPeople(userName)
			}
			return &models.AuthResult{User: userName, Email: userEmail, Uid: int(user.ID), Admin: true}, true
		} else {
			return &models.AuthResult{Err: errors.New("token 鉴权失败")}, false
		}
	}
	return nil, false
}

func (a *authentication) checkSession(c *gin.Context) *models.AuthResult {
	session := sessions.Default(c)
	user := session.Get("user")
	group := session.Get("group")
	uid := session.Get("uid")
	email := session.Get("email")
	if user == nil || group == nil || uid == nil {
		return &models.AuthResult{Err: errors.New("登录错误")}
	}
	u, err := s_auth.CasbinServiceGlobal.Count("p_type = ? and v0 = ? and v1 = ?", []interface{}{"g", user.(string), "admin"})
	if err != nil {
		return &models.AuthResult{Err: errors.New("adapter 异常")}
	}
	if email == nil {
		email = s_auth.UserServiceGlobal.GetUserEmail(&auth.User{
			UserName: user.(string),
		})
	}
	return &models.AuthResult{User: user.(string), Group: group.(string), Email: email.(string), Uid: uid.(int), Admin: u > 0}
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
	// if using token auth, just go
	_, ok := a.GetHeaderToken(c, "Authorization")
	if ok {
		res, _ := a.checkToken(c)
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
	var response models.PeopleRspInfo
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

func (a *authentication) GetKaNiRolesUrl(user_email string) string {
	return fmt.Sprintf("https://ei.byted.org/ratak/user/%s/roles/?need_tags=0&tag=", user_email)
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

func (a *authentication) JudgeHasKaNiAdminPerm(perms []string) bool {
	if utils.ContainsString(perms, "all") {
		return true
	}
	//if utils.ContainsString(perms, "admin") {
	//	return true
	//}
	return false
}

func (a *authentication) JudgeHasKaNiPerm(typ string, perms []string) bool {
	if a.JudgeHasKaNiAdminPerm(perms) {
		return true
	}
	if utils.ContainsString(perms, typ) {
		return true
	}
	return false
}

func (a *authentication) JudgeHasKaNiRole(authentication *models.AuthResult, role_name string) bool {
	user, err := s_auth.UserServiceGlobal.Query("id = ?", []interface{}{authentication.Uid})
	if err != nil {
		log.Println("user not exist")
		return false
	}
	email := s_auth.UserServiceGlobal.GetUserEmail(user)
	U := a.GetKaNiHasRoleUrl(email, role_name)
	res := a.HttpGetWithBasicAuth(U, KaNiApp.GetBasicId(), KaNiApp.GetBasicSecret())
	log.Println(string(res))
	var response map[string]interface{}
	json.Unmarshal(res, &response)
	ok_res, ok := response["ok"]
	if !ok {
		return false
	}
	return ok_res.(bool)
}

func (a *authentication) GetKaNiRoles(user_email string) []string {
	U := a.GetKaNiRolesUrl(user_email)
	res := a.HttpGetWithBasicAuth(U, KaNiApp.GetBasicId(), KaNiApp.GetBasicSecret())
	log.Println(string(res))
	var response []map[string]interface{}
	json.Unmarshal(res, &response)
	var result []string
	for _, i := range response {
		result = append(result, i["key"].(string))
	}
	return result
}
