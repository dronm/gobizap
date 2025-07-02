package service

import (
	"reflect"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"
)

//Controller
type Service_Controller struct {
	gobizap.Base_Controller
}

func NewController_Service() *Service_Controller{
	c := &Service_Controller{gobizap.Base_Controller{ID: "Service", PublicMethods: make(gobizap.PublicMethodCollection)}}
	
	//************************** method reload_config **********************************
	c.PublicMethods["reload_config"] = &Service_Controller_reload_config{
		gobizap.Base_PublicMethod{
			ID: "reload_config",
			Fields: nil,
		},
	}

	//************************** method reload_version **********************************
	c.PublicMethods["reload_version"] = &Service_Controller_reload_version{
		gobizap.Base_PublicMethod{
			ID: "reload_version",
			Fields: nil,
		},
	}
	
	return c
}

//**************************************************************************************
//Public method: reload_config
type Service_Controller_reload_config struct {	
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Service_Controller_reload_config) Unmarshal(payload []byte) (res reflect.Value, err error) {
	return res, nil
}

//custom method
func (pm *Service_Controller_reload_config) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	if err := app.ReloadAppConfig(); err != nil {
		return gobizap.NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	return nil
}

//**************************************************************************************
//Public method: reload_version
type Service_Controller_reload_version struct {	
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Service_Controller_reload_version) Unmarshal(payload []byte) (res reflect.Value, err error) {
	return res, nil
}

//custom method
func (pm *Service_Controller_reload_version) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	if err := app.LoadAppVersion(); err != nil {
		return gobizap.NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	
	return nil
}

