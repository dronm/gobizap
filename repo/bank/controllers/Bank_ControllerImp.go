package bank

import (
	"reflect"	
	
	models "github.com/dronm/gobizap/repo/bank/models"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"	
)



//Method implemenation get_object
func (pm *Bank_Controller_get_object) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["BankList"], &models.BankList{}, sock.GetPresetFilter("BankList"))	
}

//Method implemenation get_list
func (pm *Bank_Controller_get_list) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["BankList"], &models.BankList{}, sock.GetPresetFilter("BankList"))	
}


//Method implemenation complete
func (pm *Bank_Controller_complete) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return gobizap.CompleteOnArgs(app, resp, rfltArgs, app.GetMD().Models["BankList"], &models.BankList{}, sock.GetPresetFilter("BankList"))	
}

