// @Time:       2019/12/30 下午5:22

package ems

import (
	"fmt"
	"log"
	"net/http"

	"github.com/qinxin0720/QcloudSms-go/QcloudSms"
)

const (
	AppID       = 1400300638
	AppKey      = "108f5cb8eb0b65d9f5475dadc2a45420"
	TemplateId2 = 511596 // 注册会员通知
	TemplateId  = 511591 // 验证码
)

var params = []string{"12312", "2"}

var SmsGlobal Ems

type EmsObj struct {
	AppID      int
	AppKey     string
	TemplateId int
}

func init() {
	tmp := new(EmsObj)
	tmp.AppID = AppID
	tmp.AppKey = AppKey
	tmp.TemplateId = TemplateId
	SmsGlobal = tmp
}

func (s *EmsObj) SendEms(phone string) {
	qcloudsms, err := QcloudSms.NewQcloudSms(s.AppID, s.AppKey)
	if err != nil {
		log.Println(err)
		return
	}
	qcloudsms.SmsSingleSender.SendWithParam(86, phone, s.TemplateId, params, "", "", "", s.Callback)
}

func (s *EmsObj) Callback(err error, resp *http.Response, resData string) {
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Println("response data: ", resData)
	}
}
