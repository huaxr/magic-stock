// @Time:       2019/12/17 下午5:01

package engine

import (
	"encoding/json"
	"log"

	"code.byted.org/data/databus_client"
	"code.byted.org/gopkg/logs"
)

var DatabusProduce *Producer

//func init() {
//	DatabusProduce, _ = newProducer()
//	DatabusProduce.ProducerForever()
//}

type Producer struct {
	Producer *databus_client.DatabusCollector
	Channel  string
	Messages chan map[string]interface{}
}

func newProducer() (*Producer, error) {
	kl := new(Producer)
	kl.Channel = "bytesecurity_event_detail"
	kl.Producer = databus_client.NewDefaultCollector()
	kl.Messages = make(chan map[string]interface{}, 100000)
	return kl, nil
}

func (p *Producer) Send(e map[string]interface{}) {
	result, err := json.Marshal(e)
	if err == nil {
		err = p.Producer.Collect(p.Channel, result, nil, 0)
		if err != nil {
			log.Println("Error send to databus", err)
		}
	} else {
		logs.Errorf("Json marshal err: %v", err)
	}
}

func (p *Producer) ProducerForever() {
	log.Println("[*] Starting producer forever...")
	for i := 1; i <= 1; i++ {
		go func() {
			for {
				select {
				case m := <-p.Messages:
					//p.OnCallHandler(m)
					p.Send(m)
					log.Println("[*] Sending to es and oncall:", m)
				}
			}
		}()
	}
}
