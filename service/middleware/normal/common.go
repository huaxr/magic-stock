package normal

import (
	"fmt"
	"magic/stock/model"
	"magic/stock/service/conf"
	"magic/stock/service/wechat"
	"net/http"
	"runtime/debug"
	"strings"

	"magic/stock/service/alarm"
	"magic/stock/service/check"
	"magic/stock/service/jwt"

	"github.com/gin-gonic/gin"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authentication := check.Authentication.JudgeApi(c)
		if authentication.Err != nil {
			// 注意: c.json 会对 & unicode 编码
			token := c.DefaultQuery("token", "")
			//c.String(400, wechat.WechatGlobal.GetCodeUrl())
			c.JSON(200, gin.H{"error_code": 2, "err_msg": "请登录", "data": wechat.WechatGlobal.GetCodeUrl(conf.Config.WxRedirect + "?token=" + token)})
			c.Abort()
			return
		} else {
			c.Set("auth", authentication)
			c.Next()
		}
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_auth, ok := c.Get("auth")
		if !ok {
			panic("middleware 注册顺序有误")
			c.Abort()
			return
		}
		authentication := _auth.(*model.AuthResult)
		if !authentication.Member {
			c.JSON(200, gin.H{"error_code": 1, "err_msg": "没有权限执行此操作", "data": nil})
			c.Abort()
			return
		}
		c.Next()
	}
}

func DebugCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// dynamic using the given origin . when using "*" which will disable cookie by chrome save reasons
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", strings.Join([]string{"content-type", jwt.JWToken}, ","))
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, gin.H{"error_code": 0, "err_msg": nil, "data": "Options Request Success!"})
			c.Abort()
			return
		}
	}
}

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if err := recover(); err != nil {
				userName := "UnKnown"
				_auth, ok := c.Get("auth")
				if ok {
					userName = _auth.(*model.AuthResult).User
				}
				_ = alarm.AlarmClient.Report(nil, fmt.Sprintf("user:%s, recover err: %v, stack: %s, path: %s", userName, err, debug.Stack(), c.Request.URL.Path), false, "")
				c.JSON(200, gin.H{"error_code": 1, "err_msg": "unknown error, please refer: " + alarm.AlarmClient.GetType(), "data": fmt.Sprint(err)})
				c.Abort()
				return
			}
		}(c)
		c.Next()
	}
}
