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
type SrvStatistics struct {
	Name fields.ValText `json:"name"`
	Max_client_count fields.ValInt `json:"max_client_count"`
	Client_count fields.ValInt `json:"client_count"`
	Downloaded_bytes fields.ValUint `json:"downloaded_bytes"`
	Uploaded_bytes fields.ValUint `json:"uploaded_bytes"`
	Handshakes fields.ValUint `json:"handshakes"`
	Run_seconds fields.ValUint `json:"run_seconds"`
	Mem_stat_sys fields.ValUint `json:"mem_stat_sys"`
	Mem_stat_lookups fields.ValUint `json:"mem_stat_lookups"`
	
	Mem_stat_heap_sys fields.ValUint `json:"mem_stat_heap_sys"`
	Mem_stat_heap_inuse fields.ValUint `json:"mem_stat_heap_inuse"`
	Mem_stat_heap_objects fields.ValUint `json:"mem_stat_heap_objects"`
	
	Mem_stat_stack_sys fields.ValUint `json:"mem_stat_stack_sys"`
	Mem_stat_stack_inuse fields.ValUint `json:"mem_stat_stack_inuse"`
	
}

func NewModelMD_SrvStatistics() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(SrvStatistics{})),
		ID: "SrvStatistics_Model",
	}
}

