package login

import (
	"reflect"	
	"encoding/json"
	
	"github.com/dronm/gobizap/repo/login/models"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/model"
)

//Controller
type LoginDevice_Controller struct {
	gobizap.Base_Controller
}

func NewController_LoginDevice() *LoginDevice_Controller{
	c := &LoginDevice_Controller{gobizap.Base_Controller{ID: "LoginDevice", PublicMethods: make(gobizap.PublicMethodCollection)}}	
	
	
	
	
	
	
	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &LoginDevice_Controller_get_list{
		gobizap.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},
	}
	
	//************************** method switch_banned *************************************
	c.PublicMethods["switch_banned"] = &LoginDevice_Controller_switch_banned{
		gobizap.Base_PublicMethod{
			ID: "switch_banned",
			Fields: fields.GenModelMD(reflect.ValueOf(models.LoginDevice_switch_banned{})),			
		},
	}
			
	
	return c
}





//************************* GET LIST **********************************************
//Public method: get_list
type LoginDevice_Controller_get_list struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *LoginDevice_Controller_get_list) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &model.Controller_get_list_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}


//************************* switch_banned **********************************************
//Public method: switch_banned
type LoginDevice_Controller_switch_banned struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *LoginDevice_Controller_switch_banned) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &models.LoginDevice_switch_banned_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

