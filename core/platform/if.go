// @Contact:    huaxinrui
// @Time:       2019/10/18 上午11:57

package platform

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type Route interface {
	InitGin() Route
	EnableHTML()
}

type Router struct {
	Router *gin.Engine
	once   *sync.Once
	env    string
}

var Gin Route

func InitGin() {
	Gin = &Router{}
	Gin.InitGin()
}
