package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

//管理总服务和集
type Registry struct {
	registions []Registion
	mutex      sync.RWMutex
}

func (r *Registry) add(reg Registion) error {
	r.mutex.Lock()
	r.registions = append(r.registions, reg)
	r.mutex.Unlock()
	//服务本身的依赖发给自己
	err := r.sendRequiresRegistion(reg)
	//通知其他依赖自己的服务更新依赖集和
	r.notifyAdd(reg)
	fmt.Println("Registry add", r.registions)
	return err
}

func (r *Registry) sendRequiresRegistion(reg Registion) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	p := patch{
		Add:    []patchEntry{},
		Remove: []patchEntry{},
	}
	for _, regServiceName := range reg.RequiresService {
		for _, regionstion := range r.registions {
			if regionstion.ServiceName == regServiceName {
				pe := patchEntry{
					ServiceName: regionstion.ServiceName,
					ServiceUrl:  regionstion.ServiceUrl,
				}
				p.Add = append(p.Add, pe)
			}
		}
	}

	//调用远程服务通知服发现
	err := r.sendPatch(p, reg.ServiceUpdateUrl)
	return err
}

func (r *Registry) notifyAdd(reg Registion) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, regionstion := range r.registions {
		//并发通知所有在线服务
		go func(regionstion Registion) {
			for _, rs := range regionstion.RequiresService {
				if rs == reg.ServiceName {
					p := patch{
						Add:    []patchEntry{},
						Remove: []patchEntry{},
					}
					pe := patchEntry{
						ServiceName: reg.ServiceName,
						ServiceUrl:  reg.ServiceUrl,
					}
					p.Add = append(p.Add, pe)
					err := r.sendPatch(p, regionstion.ServiceUpdateUrl)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}

		}(regionstion)

	}

}

func (r *Registry) notifyRemove(reg Registion) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, regionstion := range r.registions {
		//并发通知所有在线服务
		go func(regionstion Registion) {
			for _, rs := range regionstion.RequiresService {
				if rs == reg.ServiceName {
					p := patch{
						Add:    []patchEntry{},
						Remove: []patchEntry{},
					}
					pe := patchEntry{
						ServiceName: reg.ServiceName,
						ServiceUrl:  reg.ServiceUrl,
					}
					p.Remove = append(p.Remove, pe)
					err := r.sendPatch(p, regionstion.ServiceUpdateUrl)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}

		}(regionstion)

	}

}

func (r *Registry) sendPatch(p patch, url string) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(p)
	if err != nil {
		return fmt.Errorf("Encode patch error: %v", err)
	}

	_, err = http.Post(url, "application/json", buf)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *Registry) remove(reg Registion) error {

	for key, value := range r.registions {
		if value.ServiceUrl == reg.ServiceUrl {
			r.mutex.Lock()
			r.registions = append(r.registions[:key], r.registions[key+1:]...)
			r.mutex.Unlock()
			r.notifyRemove(reg)

			fmt.Println("Registry remove", r.registions)
			return nil
		}
	}
	//通知其他依赖自己的服务更新依赖集和
	return fmt.Errorf("service name: %v, service url: %v not found", reg.ServiceName, reg.ServiceUrl)
}
func (r *Registry) heartbeat(duration time.Duration)  {
	for  {

		var wg sync.WaitGroup

		for  _,v := range r.registions{
			wg.Add(1)
			go func() {
				defer wg.Done()
			}()
			passFlag := true
			for i:=0;i<3;i++ {
				resp, err := http.Get(v.HeartbeatUrl)
				if err != nil {
					log.Println(err)
				} else if resp.StatusCode == http.StatusOK {
					log.Printf("Heartbeat check passed for %v", v.ServiceName)

					if !passFlag {
						r.add(v)
					}
					break
				}

				if passFlag {
					passFlag = false
					r.remove(v)
				}
				time.Sleep(time.Second)
			}

		}
		wg.Wait()
		time.Sleep(duration)

	}

}

var reg = Registry{
	registions: make([]Registion, 0),
	mutex:      sync.RWMutex{},
}

func RegistyHandlers() {
	http.Handle("/services", &registyHandler{})
}

type registyHandler struct{}

func (rh *registyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var r Registion
		buf := bytes.NewBuffer(data)
		dec := json.NewDecoder(buf)
		err = dec.Decode(&r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = reg.add(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var r Registion
		buf := bytes.NewBuffer(data)
		dec := json.NewDecoder(buf)
		err = dec.Decode(&r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = reg.remove(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

var once = sync.Once{}
func DoHeartbeat (duration time.Duration) {
	once.Do(func() {
		go reg.heartbeat(duration)
	})
}