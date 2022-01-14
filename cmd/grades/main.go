package main

import (
	"context"
	"distributed/http/config"
	"distributed/http/grades"
	"distributed/http/log"
	"distributed/http/registry"
	"distributed/http/service"
	"fmt"
	stlog "log"
)

func main() {
	addr := fmt.Sprintf("%s:%d", config.GobalConfig.IP, config.GobalConfig.Port)
	httpaddr := fmt.Sprintf("http://%s:%d", config.GobalConfig.IP, config.GobalConfig.Port)
	reg := registry.Registion{
		ServiceName: registry.ServiceName(config.GobalConfig.ServerName),
		ServiceUrl:httpaddr,
		ServiceUpdateUrl:config.GobalConfig.ServicesUpdateUrl,
		RequiresService:[]registry.ServiceName{
			registry.LOG_SERVICE,
		},
		HeartbeatUrl:config.GobalConfig.HeartbeatUrl,

	}
	ctx,err := service.Start(context.Background(),reg, grades.RegistryHandlers,addr)

	if err != nil {
		stlog.Println(err)
		return
	}
	if logProvider, err := registry.GetProvider(registry.LOG_SERVICE); err == nil {
		log.SetLogger(config.GobalConfig.ServerName,logProvider)

	}



	stlog.Println("测试日志服务")
	<-ctx.Done()
}
