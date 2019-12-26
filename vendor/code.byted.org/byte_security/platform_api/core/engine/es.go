// @Contact:    huaxinrui
// @Time:       2019/9/16 上午10:41

package engine

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"code.byted.org/byte_security/platform_api/utils"
	"gopkg.in/olivere/elastic.v5"
)

const (
	KEY   = "uuid"
	INDEX = "bytesec_event-*"
)

type ES struct {
	EngineClient
	host    string
	channel chan uint
	client  *elastic.Client
}

func InitEngine(host string, debug bool, sniffer bool) (*ES, error) {
	es := new(ES)
	es.host = host
	es.typ = "es"
	es.mutex = new(sync.Mutex)
	es.channel = make(chan uint, 20)

	var err error
	opts := []elastic.ClientOptionFunc{elastic.SetURL(host)}

	if debug {
		opts = append(opts, elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	}

	if sniffer {
		opts = append(opts, elastic.SetSniff(sniffer))
	}

	es.client, err = elastic.NewClient(opts...)
	if err != nil {
		log.Println("Error when connecting ES'host")
		if utils.TellEnv() == "tce" { // boe could't connect es
			panic(err)
		}
	}
	return es, nil
}

func (e *ES) TermQuery(query map[string]string, callback func(err error)) []string { //[]models.EventLog{
	var err error
	e.channel <- 1
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer func(err error) {
		cancel()
		if callback != nil {
			callback(err)
		}
	}(err)

	select {
	case <-e.channel:
		q := elastic.NewBoolQuery()
		for k, v := range query {
			q = q.Must(elastic.NewTermQuery(k, v))
		}
		search := e.client.Search().Index(INDEX).Pretty(true)
		search = search.Query(q)
		sr, err := search.Do(ctx)
		if err != nil {
			fmt.Println("err querying es", err)
		}
		if sr == nil || sr.TotalHits() == 0 {
			log.Println("the es result is null")
			return nil
		}
		return e.extract(sr)
	}
}

func (e *ES) extract(res *elastic.SearchResult) []string {
	var objects []string
	for _, hit := range res.Hits.Hits {
		objects = append(objects, string(*hit.Source))
	}
	return objects
}

func (e *ES) loopConnectUntilSuccess() {

}
