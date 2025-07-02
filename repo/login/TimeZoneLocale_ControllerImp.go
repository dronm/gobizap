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

func (pm *TimeZoneLocale_Controller_insert) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["TimeZoneLocale"], &models.TimeZoneLocale_keys{}, sock.GetPresetFilter("TimeZoneLocale"))
}

//Method implemenation
func (pm *TimeZoneLocale_Controller_delete) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["TimeZoneLocale"], sock.GetPresetFilter("TimeZoneLocale"))
}

//Method implemenation
func (pm *TimeZoneLocale_Controller_get_object) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["TimeZoneLocale"], &models.TimeZoneLocale{}, sock.GetPresetFilter("TimeZoneLocale"))
}

//Method implemenation
func (pm *TimeZoneLocale_Controller_get_list) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["TimeZoneLocale"], &models.TimeZoneLocale{}, sock.GetPresetFilter("TimeZoneLocale"))
}

//Method implemenation
func (pm *TimeZoneLocale_Controller_update) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["TimeZoneLocale"], sock.GetPresetFilter("TimeZoneLocale"))
}

