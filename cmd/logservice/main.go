package main

import (
	"context"
	"fmt"
	"go-distributed/log"
	"go-distributed/registry"
	"go-distributed/service"
	stlog "log"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName: registry.LogService,
		ServiceUrl:  serviceAddress,
	}
	ctx, err := service.Start(context.Background(), host, port, log.RegisterHandlers, r)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()

	fmt.Println("Shutting down log service")
}
