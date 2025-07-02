package login

import (
	"reflect"	
	
	"github.com/dronm/gobizap/repo/login/models"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"	
	
	//"github.com/jackc/pgx/v5"
)

func (pm *LoginDeviceBan_Controller_insert) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["LoginDeviceBan"], &models.LoginDeviceBan_keys{}, sock.GetPresetFilter("LoginDeviceBan"))	
}

//Method implemenation
func (pm *LoginDeviceBan_Controller_delete) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["LoginDeviceBan"], sock.GetPresetFilter("LoginDeviceBan"))	
}

//Method implemenation
func (pm *LoginDeviceBan_Controller_get_object) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["LoginDeviceBan"], &models.LoginDeviceBan{}, sock.GetPresetFilter("LoginDeviceBan"))	
}

//Method implemenation
func (pm *LoginDeviceBan_Controller_get_list) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["LoginDeviceBan"], &models.LoginDeviceBan{}, sock.GetPresetFilter("LoginDeviceBan"))	
}

