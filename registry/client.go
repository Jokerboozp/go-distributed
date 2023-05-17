package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registration) error { // 定义RegisterService函数并接收Registration作为参数

	buf := new(bytes.Buffer)    // 创建一个新的Buffer类型变量buf
	enc := json.NewEncoder(buf) // 创建一个新的json编码器 enc 并将其设置为 buf 的输出
	err := enc.Encode(r)        // 将r编码为json格式并写入buf中
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

func ShutDownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesUrl, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deregister service. Registry "+"service responed with code %v", res.StatusCode)
	}
	return nil
}
