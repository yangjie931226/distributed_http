package main

import (
	"context"
	"distributed/http/config"
	"distributed/http/log"
	"distributed/http/registry"
	"distributed/http/service"
	stlog "log"
	"fmt"
)

func main() {
	log.Run(config.GobalConfig.LogPath)
	addr := fmt.Sprintf("%s:%d", config.GobalConfig.IP, config.GobalConfig.Port)
	httpaddr := fmt.Sprintf("http://%s:%d", config.GobalConfig.IP, config.GobalConfig.Port)

	reg := registry.Registion{
		ServiceName: registry.ServiceName(config.GobalConfig.ServerName),
		ServiceUrl:httpaddr,
		ServiceUpdateUrl:config.GobalConfig.ServicesUpdateUrl,
		RequiresService:[]registry.ServiceName{
		},
		HeartbeatUrl:config.GobalConfig.HeartbeatUrl,
	}
	ctx,err := service.Start(context.Background(), reg, log.RegistyHandlers, addr)

	if err != nil {
		stlog.Println(err)
		return
	}

	<-ctx.Done()
}
