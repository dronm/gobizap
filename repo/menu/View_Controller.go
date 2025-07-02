package menu

import (
	"encoding/json"
	"reflect"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/model"
)

//Controller
type View_Controller struct {
	gobizap.Base_Controller
}

func NewController_View() *View_Controller {
	c := &View_Controller{gobizap.Base_Controller{ID: "View", PublicMethods: make(gobizap.PublicMethodCollection)}}

	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &View_Controller_get_list{
		gobizap.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},				
	}

	//************************** method complete *************************************
	c.PublicMethods["complete"] = &View_Controller_complete{
		gobizap.Base_PublicMethod{
			ID: "complete",
			Fields: fields.GenModelMD(reflect.ValueOf(View_complete{})),
		},		
	}
	
	//************************** method get_section_list **********************************
	c.PublicMethods["get_section_list"] = &View_Controller_get_section_list{
		gobizap.Base_PublicMethod{
			ID: "get_section_list",
			Fields: model.Cond_Model_fields,
		},		
	}
	
	return c	
}

//************************* GET LIST **********************************************
//Public method: get_list
type View_Controller_get_list struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *View_Controller_get_list) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &model.Controller_get_list_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//********************************************************************************************
//Public method: complete
type View_Controller_complete struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *View_Controller_complete) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &View_complete_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Custom method
type View_Controller_get_section_list struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *View_Controller_get_section_list) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &View_get_section_list_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

