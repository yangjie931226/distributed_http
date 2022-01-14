package service

import (
	"context"
	"distributed/http/registry"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, reg registry.Registion, registyHandlerFunc func(),addr string) (context.Context, error) {
	registyHandlerFunc()
	ctx, cancel := context.WithCancel(ctx)


	srv := http.Server{
		Addr: addr,
	}
	go func() {
		log.Println(srv.ListenAndServe())
		err := registry.RegistryRemove(reg)
		if err!= nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		var stop string
		fmt.Printf("%v started. Press any key to stop \n", reg.ServiceName)
		fmt.Scanln(&stop)
		srv.Shutdown(ctx)
		err := registry.RegistryRemove(reg)
		if err!= nil {
			log.Println(err)
		}
	}()
	err := registry.RegistryAdd(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}
