// @Time:       2019/11/28 下午3:35
package control

import (
	"magic/stock/model"
	cap "magic/stock/service/captcha"

	"github.com/gin-gonic/gin"
)

type CommonIF interface {
	Response(c *gin.Context, data interface{}, err error)
	ReloadCaptcha(c *gin.Context)
}

var CommonControlGlobal CommonIF

func init() {
	tmp := new(CommonControl)
	tmp.response = new(model.HttpResponse)
	tmp.service = cap.CaptchaGlobal
	CommonControlGlobal = tmp
}

type CommonControl struct {
	response *model.HttpResponse
	service  cap.CaptchaIF
}

func (d *CommonControl) Response(c *gin.Context, data interface{}, err error) {
	c.AbortWithStatusJSON(200, d.response.Response(data, err))
}

func (d *CommonControl) ReloadCaptcha(c *gin.Context) {
	str := d.service.NewCaptcha()
	d.Response(c, str, nil)
}
