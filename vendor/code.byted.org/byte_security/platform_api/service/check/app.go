// @Time:       2019/11/15 下午4:36

package check


const (
	NAME = "byte_security"
	ID = "840"
	SECRET = "BD3E775670334DBD8D2481840E0B65B0"
)

type App struct {
	name string
	id string
	secret string
	resources []*Resource
}

var KaNiApp AppIF

func init() {
	tmp := new(App)
	tmp.name = NAME
	tmp.id = ID
	tmp.secret = SECRET
	KaNiApp = tmp
}

type AppIF interface {
	GetName() string
	GetBasicId() string
	GetBasicSecret() string
}

func (a *App) GetName() string {
	return a.name
}

func (a *App) GetBasicId() string {
	return a.id
}

func (a *App) GetBasicSecret() string {
	return a.secret
}