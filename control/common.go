// @Time:       2019/11/28 下午3:35
package control

import (
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	cap "magic/stock/service/captcha"

	"github.com/gin-gonic/gin"
)

type CommonIF interface {
	Response(c *gin.Context, data interface{}, err error)
	ReloadCaptcha(c *gin.Context)
	PaymentList(c *gin.Context)
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

func (d *CommonControl) PaymentList(c *gin.Context) {
	typ := c.DefaultQuery("type", "")
	if typ == "data" {
		var res []dal.Price
		store.MysqlClient.GetDB().Model(&dal.Price{}).Where("type = ?", "data").Find(&res)
		d.Response(c, res, nil)
		return
	} else if typ == "mask" {
		var res []dal.Price
		store.MysqlClient.GetDB().Model(&dal.Price{}).Where("type = ?", "mask").Find(&res)
		d.Response(c, res, nil)
		return
	} else {
		var res []dal.Price
		store.MysqlClient.GetDB().Model(&dal.Price{}).Where("type = ? OR type = ?", "member", "query").Find(&res)
		d.Response(c, res, nil)
		return
	}
}
