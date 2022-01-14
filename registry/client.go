package registry

import (
	"bytes"
	"distributed/http/config"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegistryAdd(reg Registion) error {

	//心跳
	heartBeatUrl, err := url.Parse(reg.HeartbeatUrl)
	if err != nil {
		return err
	}
	http.HandleFunc(heartBeatUrl.Path, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	//启动更自身服务依赖
	serviceUpdateUrl, err := url.Parse(reg.ServiceUpdateUrl)
	if err != nil {
		return err
	}
	fmt.Println(serviceUpdateUrl.Path)
	http.Handle(serviceUpdateUrl.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(reg)
	if err != nil {
		return fmt.Errorf("RegistryAdd Encode err: %v\n", err)
	}

	//远程调用注册中心的增加服务方法
	resp, err := http.Post(config.GobalConfig.RegistryServer, "application/json", buf)
	if err != nil {
		return fmt.Errorf("RegistryAdd Post err: %v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RegistryAdd Status err\n")
	}
	return nil
}

func RegistryRemove(reg Registion) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(reg)
	if err != nil {
		return fmt.Errorf("RegistryRemove Encode err: %v\n", err)
	}

	//远程调用注册中心的删除服务方法
	req, err := http.NewRequest(http.MethodDelete, config.GobalConfig.RegistryServer, buf)
	if err != nil {
		return fmt.Errorf("RegistryRemove NewRequest err: %v\n", err)
	}
	req.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("RegistryRemove DefaultClient Do err: %v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RegistryRemove Status err\n")
	}
	return nil
}

type serviceUpdateHandler struct{}

func (suh *serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var p patch
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	prov.update(p)
	fmt.Printf("Updated received %v\n", p)
}

type provider struct {
	serviceProviders map[ServiceName][]string
	mutex            *sync.RWMutex
}

func (pvd *provider) update(p patch) {
	pvd.mutex.Lock()
	defer pvd.mutex.Unlock()

	for _, add := range p.Add {
		if _, ok := pvd.serviceProviders[add.ServiceName]; ok {
			pvd.serviceProviders[add.ServiceName] = append(pvd.serviceProviders[add.ServiceName], add.ServiceUrl)
		} else {
			pvd.serviceProviders[add.ServiceName] = []string{add.ServiceUrl}
		}
	}
	for _, remove := range p.Remove {
		if services, ok := pvd.serviceProviders[remove.ServiceName]; ok {
			for index, serviceUrl := range services {
				if serviceUrl == remove.ServiceUrl {
					pvd.serviceProviders[remove.ServiceName] = append(pvd.serviceProviders[remove.ServiceName][:index],pvd.serviceProviders[remove.ServiceName][index+1:]...)
				}
			}
		}
	}
}

func (pvd *provider) get(serviceName ServiceName) (string,error) {
	pvd.mutex.RLock()
	defer pvd.mutex.RUnlock()
	if services,ok := pvd.serviceProviders[serviceName] ; ok {
		return services[int(rand.Float32() * float32(len(services)))],nil
	}
	return "", fmt.Errorf("%v service not found", serviceName)
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

var prov = provider{
	serviceProviders: map[ServiceName][]string{},
	mutex:            &sync.RWMutex{},
}
