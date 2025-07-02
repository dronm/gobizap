package contact

import (
	"reflect"	
	
	models "github.com/dronm/gobizap/repo/contact/models"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"	
	
	//"github.com/jackc/pgx/v5"
)

//********************
//Method implemenation insert
func (pm *Contact_Controller_insert) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["Contact"], &models.Contact_keys{}, sock.GetPresetFilter("Contact"))	
}

//Method implemenation delete
func (pm *Contact_Controller_delete) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["Contact"], sock.GetPresetFilter("Contact"))	
}

//Method implemenation get_object
func (pm *Contact_Controller_get_object) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["ContactDialog"], &models.ContactDialog{}, sock.GetPresetFilter("ContactDialog"))	
}

//Method implemenation get_list
func (pm *Contact_Controller_get_list) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["ContactList"], &models.ContactList{}, sock.GetPresetFilter("ContactList"))	
}

//Method implemenation update
func (pm *Contact_Controller_update) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["Contact"], sock.GetPresetFilter("Contact"))	
}

func (pm *Contact_Controller_complete) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.CompleteOnArgs(app, resp, rfltArgs, app.GetMD().Models["ContactList"], &models.ContactList{}, sock.GetPresetFilter("ContactList"))	
}


