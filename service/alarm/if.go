// @Contact:    huaxinrui
// @Time:       2019/9/19 下午5:31

package alarm

import "sync"

type AlarmIF interface {
	InitAlarm()
	Report(users []string, msg string, group bool, token string) error
	GetType() string
	SetType(string)
	Put(map[string]interface{})
	GetQueue() chan map[string]interface{}
}

type ReportClient struct {
	typ   string // 报警类型
	path  string // 报警地址
	mutex *sync.Mutex
}

var AlarmClient, LarkClient AlarmIF
