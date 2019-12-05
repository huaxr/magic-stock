// @Contact:    huaxinrui
// @Time:       2019/9/26 上午11:28

package alarm

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/parnurzeal/gorequest"
)

const ()

type Lark struct {
	ReportClient
	message chan map[string]interface{}
}

func (r *Lark) InitAlarm() {
	r.typ = "lark"
	r.path = ""
	r.mutex = new(sync.Mutex)
	r.message = make(chan map[string]interface{}, 100)
	go r.consume()
}

func (r *Lark) Report(users []string, msg string, group bool, token string) error {
	if group {
		larkGroup(users, msg, token)
		return nil
	}
	larkGuy(users, msg, token)
	return nil
}

func (r *Lark) consume() {
	for {
		select {
		case ms := <-r.message:
			users := ms["users"].([]string)
			msg := ms["msg"].(string)
			group, ok := ms["is_group"]
			if !ok {
				group = false
			}
			token, ok2 := ms["token"]
			if !ok2 {
				token = ""
			}
			if len(users) > 0 && len(msg) > 0 {
				r.Report(users, msg, group.(bool), token.(string))
			}
		}
	}
}

//func (r *Lark) ReportError(err error) error{
//	return errors.New("Not implement error")
//}

func (r *Lark) GetType() string {
	return r.typ
}

func (r *Lark) SetType(t string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.typ = t
}

func (r *Lark) Put(m map[string]interface{}) {
	r.message <- m
}

func (r *Lark) GetQueue() chan map[string]interface{} {
	return r.message
}

func getUserID(name string) string {
	_, body, _ := gorequest.New().Post("").
		Send(map[string]interface{}{
			"token": "", "email": name + "@bytedance.com",
		}).End()
	type rs struct {
		Code    int    `json:"code"`
		User_id string `json:"user_id"`
	}

	var r rs
	json.Unmarshal([]byte(body), &r)
	if r.Code != 0 {
		log.Println("error user not found", name)
		return ""
	}
	return r.User_id
}

type rs struct {
	Code    int               `json:"code"`
	Channel map[string]string `json:"channel"`
}

func larkGuy(names []string, msg, token string) {
	fmt.Println(names, msg)
	for _, name := range names {
		uid := getUserID(name)
		if uid == "" {
			log.Println("error uid is not exist")
			return
		}
		_, body, _ := gorequest.New().Post("").
			Send(map[string]interface{}{
				"token": token, "user": uid,
			}).End()
		var r rs
		json.Unmarshal([]byte(body), &r)
		gorequest.New().Post("").
			Send(map[string]interface{}{
				"token": token, "channel": r.Channel["id"], "text": msg,
			}).End()
	}
}

// 群组消息
func larkGroup(channel []string, msg, token string) {
	for _, i := range channel {
		gorequest.New().Post("").
			Send(map[string]interface{}{
				"token": token, "channel": i, "text": msg,
			}).End()
	}
}

func InitLark() {
	LarkClient = NewAlarmObject(1)
	LarkClient.InitAlarm()
}
