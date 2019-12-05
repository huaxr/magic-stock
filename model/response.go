// @Time:       2019/11/28 下午3:39

package model

type HttpResponse struct {
	ErrorCode int         `json:"error_code"`
	ErrMsg    interface{} `json:"err_msg"`
	Data      interface{} `json:"data"`
}

func (b *HttpResponse) Response(data interface{}, err error) *HttpResponse {
	if err != nil {
		b.ErrorCode = 1
		b.ErrMsg = err.Error()
		b.Data = nil
	} else {
		b.ErrorCode = 0
		b.ErrMsg = ""
		b.Data = data
	}
	return b
}
