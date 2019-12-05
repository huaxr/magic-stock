// @Contact:    huaxinrui
// @Time:       2019/9/19 下午5:44

package alarm

import (
	"errors"
	"sync"

	"github.com/getsentry/raven-go"
)

type Sentry struct {
	ReportClient
	sentry  *raven.Client
	message chan map[string]interface{}
}

func (r *Sentry) InitAlarm() {
	r.typ = "Sentry"
	r.path = ""
	r.mutex = new(sync.Mutex)
	client, err := raven.New("")
	r.message = make(chan map[string]interface{}, 100)
	if err != nil {
		panic(err)
	}
	r.sentry = client
}

func (r *Sentry) GetType() string {
	return r.typ
}

func (r *Sentry) SetType(t string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typ = t
}

func (r *Sentry) Report(users []string, msg string, group bool, token string) error {
	res := r.sentry.CaptureError(errors.New(msg), nil)
	if res == "" {
		return errors.New("Error checked")
	}
	return nil
}

func (r *Sentry) Put(m map[string]interface{}) {
	r.message <- m
}

func (r *Sentry) GetQueue() chan map[string]interface{} {
	return r.message
}

func InitSentry() {
	AlarmClient = NewAlarmObject(2)
	AlarmClient.InitAlarm()
}
