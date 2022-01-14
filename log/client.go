package log

import (
	"bytes"
	"distributed/http/config"
	"fmt"
	"net/http"
	stlog "log"
)

func SetLogger(serviceName string, serviceUrl string) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", serviceName))
	stlog.SetFlags(0)
	stlog.SetOutput(&logWriter{url:serviceUrl})
}

type logWriter struct{
	url string
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	resp, err := http.Post(config.GobalConfig.LogServerUrl, "text/plain", bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to send log message. Service responded with %d - %s", resp.StatusCode, resp.Status)
	}

	return len(p), nil
}
