package main

import (
	"context"
	"fmt"
	"go-distributed/log"
	"go-distributed/registry"
	"go-distributed/service"
	stlog "log"
)

// main函数
func main() {
	// 运行log包，"./distributed.log"是日志文件路径
	log.Run("./distributed.log")
	// 将变量host和port分别设置为"localhost"和"4000"
	host, port := "localhost", "4000"
	// 使用host和port变量创建serviceAddress字符串
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	// 创建一个registry registration，使用"LogService"作为服务名称，使用serviceAddress作为服务URL
	r := registry.Registration{
		ServiceName: registry.LogService,
		ServiceUrl:  serviceAddress,
	}
	// 使用给定参数启动服务，并存储上下文和错误值
	ctx, err := service.Start(context.Background(), host, port, log.RegisterHandlers, r)
	// 如果启动服务时出现错误，则记录错误
	if err != nil {
		stlog.Fatalln(err)
	}
	// 等待上下文完成(即服务关闭)
	<-ctx.Done()

	// 打印指示日志服务正在关闭的消息
	fmt.Println("Shutting down log service")
}
