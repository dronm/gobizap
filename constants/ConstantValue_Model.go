package constants

/**
 * Andrey Mikhalevich 16/12/21
 * This file is part of the OSBE framework
 */

import (
	"reflect"
	
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/model"
)

//
type ConstantValue struct {
	Id fields.ValText `json:"id" primaryKey:"true"`
	Val fields.ValText `json:"val"`
	Val_type fields.ValText `json:"val_type"`
}

func NewModelMD_ConstantValue() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(ConstantValue{})),
		ID: "ConstantValueList_Model",
	}
}

