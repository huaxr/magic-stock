package models

type BasicResponse struct {
	ErrorCode int         `json:"error_code"`
	ErrMsg    string      `json:"err_msg"`
	Data      interface{} `json:"data"`
}

func (b *BasicResponse) Error(errMsg string) *BasicResponse {
	b.ErrorCode = 1
	b.ErrMsg = errMsg
	b.Data = nil
	return b
}

func (b *BasicResponse) Success(data interface{}) *BasicResponse {
	b.ErrorCode = 0
	b.ErrMsg = ""
	b.Data = data
	return b
}
