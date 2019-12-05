package enum

type ServiceType uint

const (
	All      ServiceType = 0
	Train    ServiceType = 1
	BaseLine ServiceType = 1
)

func (s ServiceType) Check() bool {
	return s >= 0 && s < 3
}
