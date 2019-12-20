// @Contact:    huaxinrui
// @Time:       2019/10/18 下午12:00

package platform

import (
	"magic/stock/core/store"
	"magic/stock/service/conf"
	"magic/stock/service/middleware/normal"
	"magic/stock/service/middleware/session"
	"magic/stock/service/middleware/session/mysql"
	"magic/stock/utils"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func (r *Router) InitGin() Route {
	e := new(Router)
	e.Router = gin.New()
	e.once = new(sync.Once)
	e.env = utils.TellEnv()

	e.registerMiddleware(
		gin.Logger(),
		gin.Recovery(),
		session.Sessions("session", mysql.NewGormStore(store.MysqlClient.GetDB(), []byte(conf.Config.SessionSecret))),
	)
	switch e.env {
	case "loc":
		e.bindLoc()
	case "dev":
		e.bindDev()
	case "online":
		e.bindOnline()
	}
	return e
}

func (r *Router) registerMiddleware(middleware ...gin.HandlerFunc) {
	// By default gin.DefaultWriter = os.Stdout
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.once.Do(func() {
		for _, m := range middleware {
			r.Router.Use(m)
		}
	})
}

func (r *Router) EnableHTML() {
	r.Router.Static("/static", "./static")
	r.Router.LoadHTMLFiles("./templates/index.html")
	r.Router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
}

func (r *Router) bindLoc() {
	// register your middleware here by order
	r.Router.Use(normal.DebugCORS(), normal.LoginRequired())
	r.bindRouters()
	//r.EnableHTML()
	r.Router.Run("0.0.0.0:8881")
}

func (r *Router) bindDev() {
	r.Router.Use(normal.DebugCORS())
	r.Router.Use(normal.LoginRequired())
	r.bindRouters()
	r.Router.Run()
}

func (r *Router) bindOnline() {
	//r.Router.Use(normal.Recover())
	r.Router.Use(normal.DebugCORS())
	r.Router.Use(normal.LoginRequired())
	r.bindRouters()
	r.Router.Run("0.0.0.0:8081")
}
