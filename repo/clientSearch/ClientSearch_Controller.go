package clientSearch

import (
	"reflect"	
	"encoding/json"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/fields"
	
)

//Controller
type ClientSearch_Controller struct {
	gobizap.Base_Controller
}

func NewController_ClientSearch() *ClientSearch_Controller{
	c := &ClientSearch_Controller{gobizap.Base_Controller{ID: "ClientSearch", PublicMethods: make(gobizap.PublicMethodCollection)}}	
			
	//************************** method search *************************************
	c.PublicMethods["search"] = &ClientSearch_Controller_search{
		gobizap.Base_PublicMethod{
			ID: "search",
			Fields: fields.GenModelMD(reflect.ValueOf(ClientSearch_search{})),
		},
	}
	
	return c
}

type ClientSearch_Controller_search struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *ClientSearch_Controller_search) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &ClientSearch_search_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

