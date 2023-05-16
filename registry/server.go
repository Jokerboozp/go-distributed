package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// 声明常量 ServerPort，并赋值为 ":3000"
const ServerPort = ":3000"

// 声明常量 ServicesUrl，并赋值为 "http://localhost" + ServerPort + "/services"
const ServicesUrl = "http://localhost" + ServerPort + "/services"

// 声明 struct registry，包含 registrations 和 mutex 两个字段
type registry struct {
	registrations []Registration // 存储服务注册信息的切片
	mutex         *sync.Mutex    // 互斥锁，防止多个 goroutine 同时修改 registrations 切片
}

// 定义 registry 的 add 方法，向 registrations 切片中添加 Registration，加锁确保并发安全
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

// 创建 registry 类型变量 reg，初始化其中的 registrations 和 mutex 字段
var reg = registry{
	registrations: make([]Registration, 0), // 初始化 registrations 为空切片
	mutex:         new(sync.Mutex),         // 初始化 mutex 为空互斥锁
}

// 声明 RegistryService 类型
type RegistryService struct {
}

// 实现 ServeHTTP 方法，当接收到 POST 请求时将请求体解码成 Registration 类型，然后调用 add 方法将其加入 registrations 切片中
func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("request receive") // 打印日志记录请求到达
	switch r.Method {              // 根据请求方法选择不同的处理方式
	case http.MethodPost: // 如果是 POST 请求
		dec := json.NewDecoder(r.Body) // 创建解码器 dec 来解码请求体
		var r Registration             // 声明变量 r 来存储解码后的数据
		err := dec.Decode(&r)          // 将请求体解码成 Registration 类型，并保存到 r 变量中
		if err != nil {                // 如果解码发生错误
			log.Println(err) // 输出错误日志
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with url : %s\n", r.ServiceName, r.ServiceUrl) // 输出日志记录服务注册信息
		err = reg.add(r)                                                              // 调用 registry 的 add 方法将新注册的服务信息加入 registrations 中
		if err != nil {                                                               // 如果 add 方法返回了错误
			log.Println(err) // 输出错误日志
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default: // 如果请求方法不是 POST
		w.WriteHeader(http.StatusMethodNotAllowed) // 返回状态码 405
	}
}
