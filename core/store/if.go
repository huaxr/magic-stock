// @Contact:    huaxinrui
// @Time:       2019/9/23 上午11:26

package store

import (
	"magic/stock/model"
	"sync"
)

type StoreIF interface {
	MI
	Query(*model.NewQuery) (interface{}, error)
	Count(*model.NewQuery) (int, error)
	GetType() string
	NewQuery() *model.NewQuery
}

type StorageClient struct {
	typ   string // 存储类型
	mutex *sync.Mutex
	debug bool
}
