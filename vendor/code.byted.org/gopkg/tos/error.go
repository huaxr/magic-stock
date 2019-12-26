package tos

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrChecksum        = errors.New("missmatch md5")
	ErrContentTooSmall = errors.New("content too small")
)

type ErrRes struct {
	Success int `json:"success"`
	Err     struct {
		HttpCode int
		Code     int    `json:"code"`
		Message  string `json:"message"`
	} `json:"error"`
	RemoteAddr string `json:"-"`
	RequestID  string `json:"-"`
}

func (e ErrRes) Error() string {
	return fmt.Sprintf("code=%d message=%s remoteAddr=%s reqID=%s", e.Err.Code, e.Err.Message, e.RemoteAddr, e.RequestID)
}

func DecodeErr(res *http.Response) error {
	decoder := json.NewDecoder(res.Body)
	errRes := new(ErrRes)
	if err := decoder.Decode(errRes); err != nil {
		errRes.Err.Code = res.StatusCode
		errRes.Err.Message = http.StatusText(res.StatusCode)
	}
	errRes.Err.HttpCode = res.StatusCode
	errRes.RequestID = res.Header.Get("X-Tos-Request-Id")
	errRes.RemoteAddr = res.Request.Host
	return errRes
}

func IsObjectNotRestored(err error) bool {
	if resErr, ok := err.(*ErrRes); ok && resErr.Err.Code == 4032 {
		return true
	}
	return false
}

func IsNotArchiveObject(err error) bool {
	if resErr, ok := err.(*ErrRes); ok && resErr.Err.Code == 4034 {
		return true
	}
	return false
}

func IsRestoreInProgress(err error) bool {
	if resErr, ok := err.(*ErrRes); ok && resErr.Err.Code == 4033 {
		return true
	}
	return false
}

func IsNotFound(err error) bool {
	er, ok := err.(*ErrRes)
	if ok && (er.Err.Code == 4008 || er.Err.Code == 404) {
		return true
	}
	return false
}
