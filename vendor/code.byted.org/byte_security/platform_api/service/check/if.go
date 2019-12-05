// @Contact:    huaxinrui
// @Time:       2019/10/25 下午3:28

package check

import (
	"context"

	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/platform_api/models"
	"github.com/gin-gonic/gin"
)

var (
	// 获取用户组信息
	GET         = "GET"
	POST        = "POST"
	URL         = "https://ee.byted.org/ratak/employees/%s/groups/"
	GROUPURL    = "https://open.byted.org/people/employee/?email=%s@bytedance.com"
	PEOPLETOKEN = "Basic MTExOjU5QTFEMTlENUJGRjQ2QzdCRkVENDQ2QzBCNUU1RjY2"
)

// 参数检查接口
type ParamCheck interface {
	ParseSearchLikeWhere(c *gin.Context, searchFields ...string) (where string, args []interface{})
	ParseSearchEqualWhere(c *gin.Context, searchFields ...string) (where string, args []interface{})
	ParseParamsToWhereArgs(c *gin.Context, allowed []string, using_blurry bool) (string, []interface{})
	GetParamsBlurry(c *gin.Context, allowed []string) (string, []interface{})
	GetParamsSpecific(c *gin.Context, allowed []string) (string, []interface{})
	GetPagination(c *gin.Context) (offset int, limit int)
	CheckValue(k string, detail map[string]interface{}, t interface{}) interface{}
	CheckLevel(l string) int
}

// 认证检查接口
type AuthenticationIF interface {
	HttpGetWithToken(url, token string) []byte
	HttpGetWithBasicAuth(url string, key, sec string) []byte
	GetHeaderToken(c *gin.Context, head string) (string, bool)
	JudgeApi(c *gin.Context) *models.AuthResult
	TellAdmin(auth *models.AuthResult) error
	// 通过ticket拿用户名
	GetUserNameBySSOTicket(c *gin.Context) (*models.SsoRes, error)
	// 本地环境通过用户名获取组
	GetGroupByNameLoc(ctx context.Context, name string) map[string]interface{}
	// 异步/同步 boe tce环境通过用户名获取组 （异步并发下可能出现问题）
	GetGroupByNameAsync(ctx context.Context, name string) map[string]interface{}
	GetGroupByName(ctx context.Context, name string) map[string]interface{}
	GetLoginUrl() string
	GetUserInfo(name string) (user *auth.User, department string)

	// 查询资源权限
	GetKaNiUrl(user string) string
	GetKaNi(user string) map[string][]string
	QueryEmployeeHasPerm(user, typ, perm string) bool
	QueryEmployeeAllPerm(user, typ string) []string
	JudgeHasKaNiPerm(typ string, perms []string) bool

	// 查询角色权限
	GetKaNiHasRoleUrl(user_email, role_name string) string
	// 判断用户是否存在 role_name 的角色
	JudgeHasKaNiRole(user_email, role_name string) bool
}

// 安全检测接口
type SecurityIF interface {
	PreventFileAnyUpload(suffix string) error
	PreventXSS(content string) string
	PreventSSRF(url string) error
	PreventXXE(url string) error
	PreventSystemCommand(exec string) error
	PreventSQLI(field ...string) error
	// ...
}
