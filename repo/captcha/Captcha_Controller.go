package captcha

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 *
 * THIS FILE IS GENERATED FROM TEMPLATE build/templates/controllers/Controller.go.tmpl
 * ALL DIRECT MODIFICATIONS WILL BE LOST WITH THE NEXT BUILD PROCESS!!!
 *
 * This file contains method descriptions only.
 * Controller implimentation is in Captcha_ControllerImp.go file
 *
 */

import (
	"reflect"	
	"encoding/json"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/fields"		
)

//Controller
type Captcha_Controller struct {
	gobizap.Base_Controller
}

func NewController_Captcha() *Captcha_Controller{
	c := &Captcha_Controller{gobizap.Base_Controller{ID: "Captcha", PublicMethods: make(gobizap.PublicMethodCollection)}}	

	//************************** method get *************************************
	c.PublicMethods["get"] = &Captcha_Controller_get{
		gobizap.Base_PublicMethod{
			ID: "get",
			Fields: fields.GenModelMD(reflect.ValueOf(Captcha_get{})),
		},
	}
	
	return c
}

//************************* get **********************************************
//Public method: get
type Captcha_Controller_get struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Captcha_Controller_get) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Captcha_get_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return

}
