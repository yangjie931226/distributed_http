package registry

type ServiceName string

//服务
type Registion struct {
	ServiceName      ServiceName
	ServiceUrl       string
	RequiresService  []ServiceName
	ServiceUpdateUrl string
	HeartbeatUrl     string
}

const (
	LOG_SERVICE    = ServiceName("log_service")
	GRADES_SERVICE = ServiceName("grades_service")
)

//有变化的服务需要通知服务依赖
type patchEntry struct {
	ServiceName ServiceName
	ServiceUrl  string
}

//变化服务集和
type patch struct {
	Add    []patchEntry
	Remove []patchEntry
}
