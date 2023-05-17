package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	mutex         *sync.RWMutex  // 互斥锁，防止多个 goroutine 同时修改 registrations 切片
}

// 定义 registry 的 add 方法，向 registrations 切片中添加 Registration，加锁确保并发安全
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	err := r.sendRequiredServices(reg)
	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceUrl,
			},
		},
	})
	return err
}

// 定义一个registry类型的方法notify，参数为fullPatch
func (r registry) notify(fullPatch patch) {
	// 读写锁加读锁
	r.mutex.RLock()
	defer r.mutex.RUnlock() // 延迟执行解锁操作

	// 遍历registrations数组中的每一个元素，将其赋值给reg
	for _, reg := range r.registrations {
		// 使用go启动一个协程
		go func(reg Registration) {
			// 遍历Registration结构体中的RequiredServices字段
			for _, reqService := range reg.RequiredServices {
				// 定义一个patch类型的变量p
				p := patch{
					Added:   []patchEntry{},
					Removed: []patchEntry{},
				}
				// 定义一个布尔类型的变量sendUpdate，用于判断是否向外发送更新信息
				sendUpdate := false
				// 遍历fullPatch中的Added字段
				for _, added := range fullPatch.Added {
					// 如果added的Name字段等于reqService
					if added.Name == reqService {
						// 将added添加到p的Added字段中
						p.Added = append(p.Added, added)
						// 将sendUpdate设置为true
						sendUpdate = true
					}
				}
				// 遍历fullPatch中的Removed字段
				for _, removed := range fullPatch.Removed {
					// 如果removed的Name字段等于reqService
					if removed.Name == reqService {
						// 将removed添加到p的Removed字段中
						p.Removed = append(p.Removed, removed)
						// 将sendUpdate设置为true
						sendUpdate = true
					}
				}
				// 如果sendUpdate为true
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL) // 调用registry的方法sendPatch，向reg.ServiceUpdateURL发送p的内容
					if err != nil {                             // 如果出现错误
						log.Println(err) // 打印错误信息
						return
					}
				}
			}
		}(reg) // 将reg作为参数传递给go函数
	}
}

// 定义了一个名为registry的结构体类型，代表了服务注册中心
func (r registry) sendRequiredServices(reg Registration) error {
	// 上读锁，防止其他 goroutine 修改注册中心数据，推迟解锁操作
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var p patch // 定义一个 patch 类型的变量 p
	// 遍历注册中心已经注册的服务
	for _, serviceReg := range r.registrations {
		// 遍历当前给定的 Registration 实例所需的服务
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				// 将已存在的服务添加到 patch 类型实例 p 中
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceUrl,
				})
			}
		}
	}
	// 发送更新到服务所定义的 URL，如果有错误则返回错误信息
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil // 成功，则返回 nil
}

// 定义了一个方法 sendPatch，参数为 patch 和 url，返回值为 error 类型
func (r registry) sendPatch(p patch, url string) error {
	// 使用 json.Marshal 方法将 patch 对象序列化为 JSON 字节数组
	d, err := json.Marshal(p)
	if err != nil {
		return err // 若出现错误，则返回该错误
	}
	// 使用 http.Post 方法发送 POST 请求，请求体为 JSON 数据
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err // 若出现错误，则返回该错误
	}
	return nil // 返回 nil 表示无错误
}

// 定义了一个方法 remove，参数为 url，返回值为 error 类型
func (r *registry) remove(url string) error {
	// 遍历 registrations 数组中的所有元素
	for i := range reg.registrations {
		// 如果当前元素的 ServiceUrl 字段等于指定 url
		if reg.registrations[i].ServiceUrl == url {
			// 调用 notify 方法，将包含要删除的服务信息的 patch 对象作为参数传入
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: r.registrations[i].ServiceName,
						URL:  r.registrations[i].ServiceUrl,
					},
				},
			})
			// 获取互斥锁，修改 registrations 数组，然后释放互斥锁
			r.mutex.Lock()
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil // 返回 nil 表示删除成功
		}
	}
	// 如果未找到要删除的服务，则返回一个错误
	return fmt.Errorf("service at URL %s not found", url)
}

// 创建 registry 类型变量 reg，初始化其中的 registrations 和 mutex 字段
var reg = registry{
	registrations: make([]Registration, 0), // 初始化 registrations 为空切片
	mutex:         new(sync.RWMutex),       // 初始化 mutex 为空互斥锁
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
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("removing service at URL : %s", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default: // 如果请求方法不是 POST
		w.WriteHeader(http.StatusMethodNotAllowed) // 返回状态码 405
	}
}
