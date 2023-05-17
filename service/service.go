package service

import (
	"context"
	"fmt"
	"go-distributed/registry"
	"log"
	"net/http"
)

// Start函数启动服务，用于注册服务并启动HTTP服务。
func Start(ctx context.Context, host, port string, registerHandlersFunc func(), reg registry.Registration) (context.Context, error) {
	registerHandlersFunc() // 注册处理函数

	ctx = startService(ctx, reg.ServiceName, host, port) // 启动HTTP服务

	err := registry.RegisterService(reg) // 向注册中心注册服务
	if err != nil {
		return ctx, err // 如果发生错误，返回上下文和错误
	}
	return ctx, nil // 返回上下文
}

// startService函数启动HTTP服务，并在这个函数创建的goroutine启动HTTP服务器。
func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx) // 创建一个新的上下文和一个可以取消的方法
	var srv http.Server
	srv.Addr = host + ":" + port // 拼接host和端口号

	go func() {
		log.Println(srv.ListenAndServe()) // 启动HTTP服务器
		err := registry.ShutDownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Fatal(err)
		}
		cancel() // 如果启动失败，将取消上下文传递给cancel变量
	}()

	go func() {
		fmt.Printf("%v started.Press any key to stop\n", serviceName) // 输出服务名称并提示用户按任意键停止服务
		var s string
		fmt.Scanln(&s) // 用户按下任意键
		err := registry.ShutDownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Fatal(err)
		}
		srv.Shutdown(ctx) // 关闭HTTP服务器
		cancel()          // 取消上下文传递给cancel变量
	}()
	return ctx // 返回上下文
}
