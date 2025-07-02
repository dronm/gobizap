package menu

import (
	"reflect"	
	"encoding/json"
	
	"github.com/dronm/gobizap"
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/model"
)

//Controller
type VariantStorage_Controller struct {
	gobizap.Base_Controller
}

func NewController_VariantStorage() *VariantStorage_Controller{
	c := &VariantStorage_Controller{gobizap.Base_Controller{ID: "VariantStorage", PublicMethods: make(gobizap.PublicMethodCollection)}}
	
	keys_fields := fields.GenModelMD(reflect.ValueOf(VariantStorage_keys{}))
	
	//************************** method insert **********************************
	c.PublicMethods["insert"] = &VariantStorage_Controller_insert{
		gobizap.Base_PublicMethod{
			ID: "insert",
			Fields: fields.GenModelMD(reflect.ValueOf(VariantStorage{})),
			EventList: gobizap.PublicMethodEventList{"VariantStorage.insert"},
		},
	}
	
	//************************** method delete *************************************
	c.PublicMethods["delete"] = &VariantStorage_Controller_delete{
		gobizap.Base_PublicMethod{
			ID: "delete",
			Fields: keys_fields,
			EventList: gobizap.PublicMethodEventList{"VariantStorage.delete"},
		},				
	}

	//************************** method update *************************************
	c.PublicMethods["update"] = &VariantStorage_Controller_update{
		gobizap.Base_PublicMethod{
			ID: "update",
			Fields: fields.GenModelMD(reflect.ValueOf(VariantStorage_old_keys{})),
			EventList: gobizap.PublicMethodEventList{"VariantStorage.update"},
		},				
	}
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &VariantStorage_Controller_get_object{
		gobizap.Base_PublicMethod{
			ID: "get_object",
			Fields: keys_fields,
		},	
	}

	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &VariantStorage_Controller_get_list{
		gobizap.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},		
	}

	
	//************************** method upsert_filter_data **********************************
	c.PublicMethods["upsert_filter_data"] = &VariantStorage_Controller_upsert_filter_data{
		gobizap.Base_PublicMethod{
			ID: "upsert_filter_data",
			Fields: fields.GenModelMD(reflect.ValueOf(VariantStorage_upsert_filter_data{})),
		},				
	}	
	//************************** method upsert_col_visib_data **********************************
	c.PublicMethods["upsert_col_visib_data"] = &VariantStorage_Controller_upsert_col_visib_data{
		gobizap.Base_PublicMethod{
			ID: "upsert_col_visib_data",
			Fields: fields.GenModelMD(reflect.ValueOf(VariantStorage_upsert_col_visib_data{})),
		},				
	}	
	//************************** method upsert_col_order_data **********************************
	c.PublicMethods["upsert_col_order_data"] = &VariantStorage_Controller_upsert_col_order_data{
		gobizap.Base_PublicMethod{
			ID: "upsert_col_order_data",
			Fields: fields.GenModelMD(reflect.ValueOf(VariantStorage_upsert_col_order_data{})),
		},				
	}	
	//************************** method get_filter_data **********************************
	c.PublicMethods["get_filter_data"] = &VariantStorage_Controller_get_filter_data{
		gobizap.Base_PublicMethod{
			ID: "get_filter_data",
			Fields: fields.GenModelMD(reflect.ValueOf(VariantStorage_get_filter_data{})),
		},				
	}	
	//************************** method get_col_visib_data **********************************
	c.PublicMethods["get_col_visib_data"] = &VariantStorage_Controller_get_col_visib_data{
		gobizap.Base_PublicMethod{
			ID: "get_col_visib_data",
			Fields: fields.GenModelMD(reflect.ValueOf(VariantStorage_get_col_visib_data{})),
		},				
	}	
	//************************** method get_col_order_data **********************************
	c.PublicMethods["get_col_order_data"] = &VariantStorage_Controller_get_col_order_data{
		gobizap.Base_PublicMethod{
			ID: "get_col_order_data",
			Fields: fields.GenModelMD(reflect.ValueOf(VariantStorage_get_col_order_data{})),
		},				
	}
	return c	
}

type VariantStorage_Controller_keys_argv struct {
	Argv VariantStorage_keys `json:"argv"`	
}

//************************* INSERT **********************************************
//Public method: insert
type VariantStorage_Controller_insert struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_insert) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}

	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//************************* DELETE **********************************************
type VariantStorage_Controller_delete struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_delete) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_keys_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//************************* GET OBJECT **********************************************
type VariantStorage_Controller_get_object struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_get_object) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_keys_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//************************* GET LIST **********************************************
//Public method: get_list
type VariantStorage_Controller_get_list struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_get_list) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &model.Controller_get_list_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//************************* UPDATE **********************************************
//Public method: update
type VariantStorage_Controller_update struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_update) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_old_keys_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Custom method
type VariantStorage_Controller_upsert_filter_data struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_upsert_filter_data) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_upsert_filter_data_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Custom method
type VariantStorage_Controller_upsert_col_visib_data struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_upsert_col_visib_data) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_upsert_col_visib_data_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Custom method
type VariantStorage_Controller_upsert_col_order_data struct {
	gobizap.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_upsert_col_order_data) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_upsert_col_order_data_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Custom method
type VariantStorage_Controller_get_filter_data struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_get_filter_data) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_get_filter_data_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Custom method
type VariantStorage_Controller_get_col_visib_data struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_get_col_visib_data) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_get_filter_data_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Custom method
type VariantStorage_Controller_get_col_order_data struct {
	gobizap.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *VariantStorage_Controller_get_col_order_data) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &VariantStorage_get_col_order_data_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

