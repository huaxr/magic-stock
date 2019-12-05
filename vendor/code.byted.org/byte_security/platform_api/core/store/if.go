// @Contact:    huaxinrui
// @Time:       2019/9/23 上午11:26

package store

import (
	"code.byted.org/byte_security/platform_api/models"
	"sync"
)

type StoreIF interface {
	MI
	Query(*models.NewQuery) (interface{}, error)
	Count(*models.NewQuery) (int, error)
	GetType() string
	NewQuery() *models.NewQuery
}

type StorageClient struct {
	typ   string // 存储类型
	mutex *sync.Mutex
	debug bool
}
