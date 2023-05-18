package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error { // 定义RegisterService函数并接收Registration作为参数

	heartbeatURL, err := url.Parse(r.HeartBeatURL)
	if err != nil {
		return err
	}
	http.HandleFunc(heartbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)    // 创建一个新的Buffer类型变量buf
	enc := json.NewEncoder(buf) // 创建一个新的json编码器 enc 并将其设置为 buf 的输出
	err = enc.Encode(r)         // 将r编码为json格式并写入buf中
	if err != nil {             // 如果出错，返回err
		return err
	}
	res, err := http.Post(ServicesUrl, "application/json", buf) // 向ServicesUrl发起POST请求并向其发布buf内容，返回响应和错误
	if err != nil {                                             // 如果出错，返回err
		return err
	}
	if res.StatusCode != http.StatusOK { // 如果响应状态码不为200
		return fmt.Errorf("failed to register service. Registry service "+"responded with code %v", res.StatusCode) // 抛出一个新的错误
	}
	return nil // 返回空值
}

type serviceUpdateHandler struct {
}

// 定义Struct：serviceUpdateHandler
func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 如果请求方法不是POST，则返回"Method Not Allowed"状态码
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 创建JSON解析器
	dec := json.NewDecoder(r.Body)
	// 定义patch类型变量p
	var p patch
	// 解析请求体r，并将其存储到变量p中
	err := dec.Decode(&p)
	// 如果发生错误，则打印错误并返回状态码"Bad Request"
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 打印更新的内容，以及更新内容（变量p）的值
	fmt.Printf("updated received %v\n", p)
	// 调用prov.Update()，传递变量p作为参数
	prov.Update(p)
}

// ShutDownService 是一个函数，将以text/plain内容类型为参数发送DELETE请求来注销服务。
func ShutDownService(url string) error {
	// 创建一个新的DELETE http请求，url为参数，并带有text/plain消息体。
	req, err := http.NewRequest(http.MethodDelete, ServicesUrl, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	// 设置请求的Content-Type为"text/plain"。
	req.Header.Add("Content-Type", "text/plain")
	// 通过DefaultClient发送请求并等待返回响应。
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// 检查响应状态码是否不等于200 OK。
	if res.StatusCode != http.StatusOK {
		// 返回一个错误，其中包含格式化后的消息指示失败。
		return fmt.Errorf("deregister service失败。注册服务的响应代码为：%v", res.StatusCode)
	}
	// 如果没有错误，则返回nil。
	return nil
}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

func (p *providers) Update(pat patch) { // 定义一个方法 Update，并传入一个 pat 的 patch 类型参数，p 为 providers 结构体指针类型参数
	p.mutex.Lock()         // 加锁
	defer p.mutex.Unlock() // 解锁，保证函数执行结束后一定会执行此行代码

	// 遍历 pat.Added 切片中的每一个元素
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok { // 如果 services 中没有名称为 patchEntry.Name 的服务，则将其初始化为一个空的切片
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL) // 将 patchEntry.URL 添加到对应服务的切片中
	}

	// 遍历 pat.Removed 切片中的每一个元素
	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; ok { // 如果 services 中有名称为 patchEntry.Name 的服务，则遍历其对应的切片
			for i := range providerURLs {
				if providerURLs[i] == patchEntry.URL { // 如果找到了对应的 URL，则从切片中删除它
					p.services[patchEntry.Name] = append(providerURLs[:i], providerURLs[i+1:]...)
				}
			}
		}
	}
}

// 定义方法get，其中p是一个类型为providers的接收器，name是ServiceName类型的参数
func (p providers) get(name ServiceName) (string, error) {
	// 从p中的services字段中获取name对应的slice（值和是否找到）
	providers, ok := p.services[name]
	// 若没找到则返回一个包含错误信息的error
	if !ok {
		return "", fmt.Errorf("no providers available for service %v", name)
	}
	// 随机获取providers中的一个元素的索引
	idx := int(rand.Float32() * float32(len(providers)))
	// 返回providers[idx]和nil
	return providers[idx], nil
}

// 定义GetProvider函数，其中name是ServiceName类型的参数，返回一个string和error
func GetProvider(name ServiceName) (string, error) {
	// 返回prov调用get方法后的结果
	return prov.get(name)
}

// 定义变量prov，值为一个providers类型的struct
var prov = providers{
	// 初始化services字段为一个空的map
	services: make(map[ServiceName][]string),
	// 初始化mutex字段为一个新的RWMutex结构体的指针
	mutex: new(sync.RWMutex),
}
