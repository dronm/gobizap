package login

import (
	"reflect"	
	
	"github.com/dronm/gobizap/repo/login/models"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/evnt"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"	
)



//Method implemenation get_object
func (pm *Login_Controller_get_object) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["LoginList"], &models.LoginList{}, sock.GetPresetFilter("LoginList"))	
}

//Method implemenation get_list
func (pm *Login_Controller_get_list) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["LoginList"], &models.LoginList{}, sock.GetPresetFilter("LoginList"))	
}

func (pm *Login_Controller_destroy_session) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	args := rfltArgs.Interface().(*evnt.Event)
	session_id_i, ok := args.Params["session_id"]
	if !ok {
		return gobizap.NewPublicMethodError(response.RESP_ER_INTERNAL, "Login_Controller_destroy_session session_id parameter is missing")
	}
	session_id, ok := session_id_i.(string)
	if !ok {
		return gobizap.NewPublicMethodError(response.RESP_ER_INTERNAL, "Login_Controller_destroy_session session_id parameter is not a string")
	}
	app.GetSessManager().SessionDestroy(session_id)	
	return nil
}

