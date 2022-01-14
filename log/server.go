package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

type fileLog string
func (fl fileLog)Write(data []byte) (int,error){
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0766)
	if err != nil {
		return 0, nil
	}
	defer f.Close()
	return f.Write(data)

}

func Run(name string)  {
	log = stlog.New(fileLog(name), "[go]", stlog.LstdFlags)
}

func write(data string) {
	log.Printf("%v\n", data)
}

func RegistyHandlers() {
	http.Handle("/log", &logHandler{})
}

type logHandler struct{}

func (logHandler *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil || len(data) == 0{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		write(string(data))
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
