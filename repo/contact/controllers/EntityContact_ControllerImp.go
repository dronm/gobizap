package contact

import (
	"reflect"	
	
	models "github.com/dronm/gobizap/repo/contact/models"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"	
	
	"github.com/dronm/session"	
)

type ClientSockSessioner interface {
	GetSession() session.Session
}

//Method implemenation insert
func (pm *EntityContact_Controller_insert) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["EntityContact"], &models.EntityContact_keys{}, sock.GetPresetFilter("EntityContact"))	
}

//Method implemenation delete
func (pm *EntityContact_Controller_delete) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["EntityContact"], sock.GetPresetFilter("EntityContact"))	
}

//Method implemenation get_object
func (pm *EntityContact_Controller_get_object) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["EntityContactList"], &models.EntityContactList{}, sock.GetPresetFilter("EntityContactList"))	
}

//Method implemenation get_list
func (pm *EntityContact_Controller_get_list) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["EntityContactList"], &models.EntityContactList{}, sock.GetPresetFilter("EntityContactList"))	
}

//Method implemenation update
func (pm *EntityContact_Controller_update) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["EntityContact"], sock.GetPresetFilter("EntityContact"))	
}

//Method implemenation complete

