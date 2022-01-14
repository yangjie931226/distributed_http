package main

import (
	"context"
	"distributed/http/config"
	"distributed/http/log"
	"distributed/http/portal"
	"distributed/http/registry"
	"distributed/http/service"
	"fmt"
	stlog "log"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	addr := fmt.Sprintf("%s:%d", config.GobalConfig.IP, config.GobalConfig.Port)
	httpaddr := fmt.Sprintf("http://%s:%d", config.GobalConfig.IP, config.GobalConfig.Port)
	reg := registry.Registion{
		ServiceName: registry.ServiceName(config.GobalConfig.ServerName),
		ServiceUrl:httpaddr,
		ServiceUpdateUrl:config.GobalConfig.ServicesUpdateUrl,
		RequiresService:[]registry.ServiceName{
			registry.LOG_SERVICE,
			registry.GRADES_SERVICE,
		},
		HeartbeatUrl:config.GobalConfig.HeartbeatUrl,

	}
	ctx,err := service.Start(context.Background(),reg, portal.RegisterHandlers,addr)

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
