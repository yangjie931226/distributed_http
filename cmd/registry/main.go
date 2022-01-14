package main

import (
	"context"
	"distributed/http/config"
	"distributed/http/registry"
	"fmt"
	stlog "log"
	"net/http"
	"time"
)

func main() {
	//注册handler
	registry.DoHeartbeat(3*time.Second)
	registry.RegistyHandlers()
	addr := fmt.Sprintf("%s:%d", config.GobalConfig.IP, config.GobalConfig.Port)
	srv := http.Server{
		Addr: addr,
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stlog.Println(srv.ListenAndServe())
		cancel()

	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop \n", config.GobalConfig.ServerName)
		var stop string
		fmt.Scanln(&stop)
		cancel()
	}()


	<-ctx.Done()
}
