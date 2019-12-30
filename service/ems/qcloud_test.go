// @Time:       2019/12/30 下午5:30

package ems

import "testing"

func TestEmsObj_Callback(t *testing.T) {
	SmsGlobal.SendEms("15210360661")
}
