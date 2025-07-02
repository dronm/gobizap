package gobizap

import(		
	"github.com/dronm/gobizap/response"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/logger"	
	
	"github.com/dronm/session"
)

// Applicationer interface is used in all standart controller methods as a parameter.
type Applicationer interface {
	GetConfig() AppConfiger
	GetLogger() logger.Logger
	GetMD() *Metadata
	GetServer(string) srv.Server
	GetServers() ServerList
	SendToClient(srv.Server, socket.ClientSocketer, *response.Response, string) error
	HandleRequest(srv.Server, socket.ClientSocketer, string, string, string, []byte, string)
	HandleJSONRequest(srv.Server, socket.ClientSocketer, []byte, string)
	HandleSession(socket.ClientSocketer) error
	DestroySession(sessID string)
	HandlePermission(socket.ClientSocketer, string, string) error
	HandleServerError(srv.Server, socket.ClientSocketer, string, string)
	GetSessManager() *session.Manager
	//XSLTransform([]byte, string, string, string) ([]byte, error)	
	//XSLToPDFTransform(string, string, []string, []byte, string, string, string) ([]byte, error)
	GetFrameworkVersion() string
	GetPermisManager() Permissioner
	GetOnPublishEvent() OnPublishEventProto
	GetDataStorage() interface{}
	PublishPublicMethodEvents(PublicMethod, map[string]interface{})
	GetEncryptKey() string
	GetBaseDir() string
	LoadAppVersion() error
	ReloadAppConfig() error
}


