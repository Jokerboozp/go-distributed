package log

import (
	"bytes"
	"fmt"
	"go-distributed/registry"
	stlog "log"
	"net/http"
)

// SetClientLogger 是一个用于设置客户端日志记录器的函数。
// 它将客户端服务的名称作为前缀设置到日志记录器中，并禁用所有标记设置。
// 它还将日志记录器的输出设置为 clientLogger 结构体的实例。
func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] -", clientService)) // 将日志前缀设置为客户端服务名称，用方括号括起来。
	stlog.SetFlags(0)                                     // 禁用所有标准标记设置。
	stlog.SetOutput(&clientLogger{url: serviceURL})       // 将日志记录器的输出设置为 clientLogger 结构体的实例，其中 url 字段为 serviceURL。
}

// clientLogger 是一个包含服务 URL 的结构体。
type clientLogger struct {
	url string
}

// Write 是 clientLogger 上的一个方法，它将日志数据写入服务 URL。
// 它将字节数组作为输入，返回 (int, error)。
func (cl clientLogger) Write(data []byte) (int, error) {
	//TODO implement me
	b := bytes.NewBuffer([]byte(data))                    // 从字节数组中创建一个新的缓冲区。
	res, err := http.Post(cl.url+"/log", "text/plain", b) // 将缓冲区数据作为 HTTP POST 请求发送到 cl.url + "/log"。
	if err != nil {
		return 0, err // 如果有错误，返回 0 和错误。
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log message. Service responded with %v", res.StatusCode) // 如果响应状态码不是 OK，则返回 0 和错误信息。
	}
	return len(data), nil // 如果成功，则返回数据长度和无错误信息。
}
