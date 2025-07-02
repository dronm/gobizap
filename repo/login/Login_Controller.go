package login

import (
	"reflect"	
	"encoding/json"
	
	"github.com/dronm/gobizap/repo/login/models"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/model"
	"github.com/dronm/gobizap/evnt"
)

//Controller
type Login_Controller struct {
	gobizap.Base_Controller
}

func NewController_Login() *Login_Controller{
	c := &Login_Controller{gobizap.Base_Controller{ID: "Login", PublicMethods: make(gobizap.PublicMethodCollection)}}	
	keys_fields := fields.GenModelMD(reflect.ValueOf(models.Login_keys{}))
	
	
	
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &Login_Controller_get_object{
		gobizap.Base_PublicMethod{
			ID: "get_object",
			Fields: keys_fields,
		},
	}
	
	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &Login_Controller_get_list{
		gobizap.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},
	}
	
	c.PublicMethods["destroy_session"] = &Login_Controller_destroy_session{
		gobizap.Base_PublicMethod{
			ID: "destroy_session",
			Fields: fields.GenModelMD(reflect.ValueOf(evnt.Event{})),
		},
	}
			
	
	return c
}

type Login_Controller_keys_argv struct {
	Argv models.Login_keys `json:"argv"`	
}



//************************* GET OBJECT **********************************************
type Login_Controller_get_object struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Login_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &models.Login_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET LIST **********************************************
//Public method: get_list
type Login_Controller_get_list struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Login_Controller_get_list) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &model.Controller_get_list_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* destroy_session **********************************************
//Public method: destroy_session
type Login_Controller_destroy_session struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Login_Controller_destroy_session) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &evnt.Event_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}


