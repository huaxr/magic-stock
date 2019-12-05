// @Contact:    huaxinrui
// @Time:       2019/10/25 下午3:28

package check

import (
	"magic/stock/model"

	"github.com/gin-gonic/gin"
)

var (
	// 获取用户组信息
	GET  = "GET"
	POST = "POST"
)

// 参数检查接口
type ParamCheck interface {
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
	JudgeApi(c *gin.Context) *model.AuthResult
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
