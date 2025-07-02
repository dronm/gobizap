package permission

import (
	"fmt"
	"reflect"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/srv"
	"github.com/dronm/gobizap/socket"
	"github.com/dronm/gobizap/response"
	
)

//Controller
type Permission_Controller struct {
	gobizap.Base_Controller
}

func NewController_Permission() *Permission_Controller{
	c := &Permission_Controller{gobizap.Base_Controller{ID: "Permission", PublicMethods: make(gobizap.PublicMethodCollection)}}
	
	//************************** method change **********************************
	c.PublicMethods["change"] = &Permission_Controller_change{
		gobizap.Base_PublicMethod{
			ID: "change",
		},
	}
	return c
}

//**************************************************************************************
//Public method: change
type Permission_Controller_change struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Permission_Controller_change) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	return res, nil
}

//custom method
func (pm *Permission_Controller_change) Run(app gobizap.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	app.GetLogger().Warn("Permission_Controller_change")
	if err := app.GetPermisManager().Reload(); err != nil {
		return gobizap.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Permission_Controller_change: %v",err))	
	}
	return nil
}

