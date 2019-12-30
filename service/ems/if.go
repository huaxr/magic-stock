// @Time:       2019/12/2 上午11:56

package ems

import "net/http"

type Ems interface {
	SendEms(phone string)
	Callback(err error, resp *http.Response, resData string)
}
