package docAttachment

import (
	"reflect"	
	"encoding/json"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/model"
)

//Controller
type Attachment_Controller struct {
	gobizap.Base_Controller
}

func NewController_Attachment() *Attachment_Controller{
	c := &Attachment_Controller{gobizap.Base_Controller{ID: "Attachment", PublicMethods: make(gobizap.PublicMethodCollection)}}	
	keys_fields := fields.GenModelMD(reflect.ValueOf(Attachment_keys{}))
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &Attachment_Controller_get_object{
		gobizap.Base_PublicMethod{
			ID: "get_object",
			Fields: keys_fields,
		},
	}
	
	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &Attachment_Controller_get_list{
		gobizap.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},
	}

	//************************** method clear_cache *************************************
	c.PublicMethods["clear_cache"] = &Attachment_Controller_clear_cache{
		gobizap.Base_PublicMethod{
			ID: "clear_cache",
			Fields: fields.GenModelMD(reflect.ValueOf(Attachment_clear_cache{})),
		},
	}
	
	//************************** method delete_file *************************************
	c.PublicMethods["delete_file"] = &Attachment_Controller_delete_file{
		gobizap.Base_PublicMethod{
			ID: "delete_file",
			Fields: fields.GenModelMD(reflect.ValueOf(Attachment_delete_file{})),
		},
	}
	//************************** method add_file *************************************
	c.PublicMethods["add_file"] = &Attachment_Controller_add_file{
		gobizap.Base_PublicMethod{
			ID: "add_file",
			Fields: fields.GenModelMD(reflect.ValueOf(Attachment_add_file{})),
		},
	}
	//************************** method get_file *************************************
	c.PublicMethods["get_file"] = &Attachment_Controller_get_file{
		gobizap.Base_PublicMethod{
			ID: "get_file",
			Fields: fields.GenModelMD(reflect.ValueOf(Attachment_get_file{})),
		},
	}			
	
	return c
}

type Attachment_Controller_keys_argv struct {
	Argv Attachment_keys `json:"argv"`	
}



//************************* GET OBJECT **********************************************
type Attachment_Controller_get_object struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Attachment_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &Attachment_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET LIST **********************************************
//Public method: get_list
type Attachment_Controller_get_list struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Attachment_Controller_get_list) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &model.Controller_get_list_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}


//************************* delete_file **********************************************
//Public method: delete_file
type Attachment_Controller_delete_file struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Attachment_Controller_delete_file) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &Attachment_delete_file_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* get_file **********************************************
//Public method: get_file
type Attachment_Controller_get_file struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Attachment_Controller_get_file) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &Attachment_get_file_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* add_file **********************************************
//Public method: add_file
type Attachment_Controller_add_file struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Attachment_Controller_add_file) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &Attachment_add_file_argv{}
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* clear_cache **********************************************
//Public method: clear_cache
type Attachment_Controller_clear_cache struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Attachment_Controller_clear_cache) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &Attachment_clear_cache_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}


