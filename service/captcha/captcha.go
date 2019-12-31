// @Time:       2019/12/31 下午2:00

package captcha

import (
	"magic/stock/service/conf"

	"github.com/gin-gonic/gin"

	"github.com/dchest/captcha"
)

func (ca *Captcha) NewCaptcha() string {
	return conf.Config.Host + "/common/captcha/" + captcha.New() + ".png"
}

func (ca *Captcha) Check(c *gin.Context) bool {
	captchaId := c.DefaultQuery("captchaId", "")
	captchaSolution := c.DefaultQuery("captchaSolution", "")
	if !captcha.VerifyString(captchaId, captchaSolution) {
		return false
	}
	return true
}
