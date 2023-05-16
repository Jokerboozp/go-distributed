package registry

// 定义一个 Registration 结构体，用于保存服务的名称和 URL 信息
type Registration struct {
	ServiceName ServiceName
	ServiceUrl  string
}

// 定义一个 ServiceName 类型为 string
type ServiceName string

// 定义一个常量 LogService 的值为 ServiceName("LogService")
const (
	LogService = ServiceName("LogService")
)
