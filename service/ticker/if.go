// @Contact:    huaxinrui
// @Time:       2019/11/3 下午6:32

package ticker

import (
	"log"
	"time"

	"magic/stock/utils"
)

type TickerIF interface {
	GetTicker(f f)
}

type f func()

var GlobalTicker TickerIF

func init() {
	tmp := new(Ticker)
	tmp.env = utils.TellEnv()
	GlobalTicker = tmp
}

type Ticker struct {
	env string
}

func (ticker *Ticker) GetTicker(f f) {
	t := time.NewTicker(3 * time.Minute)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			log.Println("starting sync src reports")
			f()
		}
	}
}

func InitTicker() {
	// 后者可以在此处注册定时任务函数
	//go GlobalTicker.GetTicker(GlobalTicker.SyncReports)
}
