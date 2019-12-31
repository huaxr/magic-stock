// @Time:       2019/12/31 下午1:59

package captcha

import "github.com/gin-gonic/gin"

type CaptchaIF interface {
	NewCaptcha() string
	Check(c *gin.Context) bool
}

var CaptchaGlobal CaptchaIF

type Captcha struct {
}

func init() {
	tmp := new(Captcha)
	CaptchaGlobal = tmp
}
