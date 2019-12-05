// @Contact:    huaxinrui
// @Time:       2019/10/8 下午6:32

package alarm

type TYPE int

const (
	LARK TYPE = iota + 1
	SENTRY
)

type AlarmMessage struct {
	Type     TYPE
	Messages map[string]interface{}
	Group    bool
}

func NewAlarmObject(t TYPE) AlarmIF {
	switch t {
	case LARK:
		return new(Lark)
	case SENTRY:
		return new(Sentry)
	default:
		panic("NO TYPE FOUND")
	}
}
