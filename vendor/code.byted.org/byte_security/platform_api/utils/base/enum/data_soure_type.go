package enum

type DataSourceType uint

const (
	Kafka  DataSourceType = 0
	UpLoad DataSourceType = 1
)

func (d DataSourceType) Check() bool {
	return d >= 0 && d < 2
}
