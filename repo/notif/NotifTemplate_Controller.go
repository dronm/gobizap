package notif

import(
	"reflect"
	"encoding/json"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/model"
)

//Controller
type NotifTemplate_Controller struct {
	gobizap.Base_Controller
}

func NewController_NotifTemplate() *NotifTemplate_Controller{
	c := &NotifTemplate_Controller{gobizap.Base_Controller{ID: "NotifTemplate", PublicMethods: make(gobizap.PublicMethodCollection)}}	
	keys_fields := fields.GenModelMD(reflect.ValueOf(NotifTemplate_keys{}))
	
	//************************** method insert **********************************
	c.PublicMethods["insert"] = &NotifTemplate_Controller_insert{
		gobizap.Base_PublicMethod{
			ID: "insert",
			Fields: fields.GenModelMD(reflect.ValueOf(NotifTemplate{})),
			EventList: gobizap.PublicMethodEventList{"NotifTemplate.insert"},
		},
	}
	
	//************************** method delete *************************************
	c.PublicMethods["delete"] = &NotifTemplate_Controller_delete{
		gobizap.Base_PublicMethod{
			ID: "delete",
			Fields: keys_fields,
			EventList: gobizap.PublicMethodEventList{"NotifTemplate.delete"},
		},
	}
	
	//************************** method update *************************************
	c.PublicMethods["update"] = &NotifTemplate_Controller_update{
		gobizap.Base_PublicMethod{
			ID: "update",
			Fields: fields.GenModelMD(reflect.ValueOf(NotifTemplate_old_keys{})),
			EventList: gobizap.PublicMethodEventList{"NotifTemplate.update"},
		},
	}
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &NotifTemplate_Controller_get_object{
		gobizap.Base_PublicMethod{
			ID: "get_object",
			Fields: keys_fields,
		},
	}
	
	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &NotifTemplate_Controller_get_list{
		gobizap.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},
	}
	
			
	
	return c
}

type NotifTemplate_Controller_keys_argv struct {
	Argv NotifTemplate_keys `json:"argv"`	
}

//************************* INSERT **********************************************
//Public method: insert
type NotifTemplate_Controller_insert struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_insert) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &NotifTemplate_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* DELETE **********************************************
type NotifTemplate_Controller_delete struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_delete) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &NotifTemplate_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET OBJECT **********************************************
type NotifTemplate_Controller_get_object struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &NotifTemplate_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET LIST **********************************************
//Public method: get_list
type NotifTemplate_Controller_get_list struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_get_list) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &model.Controller_get_list_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* UPDATE **********************************************
//Public method: update
type NotifTemplate_Controller_update struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_update) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &NotifTemplate_old_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

