// About package holds About application information.
// It contains a controller and a model. Controller has a get_object method.
// Information is fetched from application configuration.
//
// OSBE project by Andrey Mikhalevich
package about

import (
	"reflect"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"	
)

//Controller
type About_Controller struct {
	gobizap.Base_Controller
}

func NewController_About() *About_Controller{
	c := &About_Controller{gobizap.Base_Controller{ID: "About", PublicMethods: make(gobizap.PublicMethodCollection)}}

	c.PublicMethods = make(gobizap.PublicMethodCollection)
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &About_Controller_get_object{
		gobizap.Base_PublicMethod{
			ID: "get_object",
		},
	}
	
	return c
}

//************************* GET OBJECT **********************************************
type About_Controller_get_object struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *About_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {	
	var res reflect.Value
	return res, nil
}

//Method implemenation
func (pm *About_Controller_get_object) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	conf := app.GetConfig()	
	m_row := &About{Author: fields.ValText{TypedValue: conf.GetAuthor()},
		Tech_mail: fields.ValText{TypedValue: conf.GetTechMail()},
		App_name: fields.ValText{TypedValue: conf.GetAppID()},
		Fw_version: fields.ValText{TypedValue: app.GetFrameworkVersion()},
		App_version: fields.ValText{TypedValue: app.GetMD().Version.Value},
		Db_name: fields.ValText{},
	}
	resp.AddModelFromStruct("About_Model", m_row)
	return nil
}



