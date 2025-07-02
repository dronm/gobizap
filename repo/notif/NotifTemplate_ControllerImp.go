package notif

import (
	"reflect"	
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"	
	
)

//Method implemenation insert
func (pm *NotifTemplate_Controller_insert) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["NotifTemplate"], &NotifTemplate_keys{}, sock.GetPresetFilter("NotifTemplate"))	
}

//Method implemenation delete
func (pm *NotifTemplate_Controller_delete) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["NotifTemplate"], sock.GetPresetFilter("NotifTemplate"))	
}

//Method implemenation get_object
func (pm *NotifTemplate_Controller_get_object) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["NotifTemplate"], &NotifTemplate{}, sock.GetPresetFilter("NotifTemplate"))	
}

//Method implemenation get_list
func (pm *NotifTemplate_Controller_get_list) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["NotifTemplateList"], &NotifTemplateList{}, sock.GetPresetFilter("NotifTemplateList"))	
}

//Method implemenation update
func (pm *NotifTemplate_Controller_update) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["NotifTemplate"], sock.GetPresetFilter("NotifTemplate"))	
}


