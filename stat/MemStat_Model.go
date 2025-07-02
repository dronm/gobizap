package stat

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 *
 */

import (
	"reflect"
	
		
	"github.com/dronm/gobizap/fields"
	"github.com/dronm/gobizap/model"
)

//
type MemStat struct {
	Sys fields.ValUint `json:"sys"`
	Lookups fields.ValUint `json:"lookups"`
	
	Heap_sys fields.ValUint `json:"heap_sys"`
	Heap_inuse fields.ValUint `json:"heap_inuse"`
	Heap_objects fields.ValUint `json:"heap_objects"`
	
	Stack_sys fields.ValUint `json:"stack_sys"`
	Stack_inuse fields.ValUint `json:"stack_inuse"`
	
}

func NewModelMD_MemStat() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(MemStat{})),
		ID: "MemStat_Model",
	}
}

