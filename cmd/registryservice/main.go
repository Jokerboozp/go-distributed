package main

import (
	"context"
	"fmt"
	"go-distributed/registry"
	"log"
	"net/http"
)

// 定义 main 函数
func main() {
	registry.SetupRegistryService()
	// 注册服务，注册地址为 "/services"，使用 RegistryService 结构体作为处理器
	http.Handle("/services", &registry.RegistryService{})

	// 创建一个上下文和取消函数
	ctx, cancel := context.WithCancel(context.Background())

	// 确保当主函数返回时调用取消函数
	defer cancel()

	// 创建一个 HTTP 服务器实例
	var srv http.Server

	// 将服务器地址设置为 SERVER_PORT 的值
	srv.Addr = registry.ServerPort

	// 开始一个 Goroutine，监听已注册的端点上的请求并提供服务
	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	// 开始一个 Goroutine，打印信息并等待输入以关闭服务器
	go func() {
		fmt.Println("Registry service started. Press any key to stop")
		var s string
		// 等待 console 输入
		fmt.Scanln(&s)

		// 优雅地关闭 HTTP 服务器
		srv.Shutdown(ctx)

		// 取消上下文
		cancel()
	}()

	// 等待上下文被取消
	<-ctx.Done()

	// 打印关于关闭服务的消息
	fmt.Println("shutting down registry service")
}
