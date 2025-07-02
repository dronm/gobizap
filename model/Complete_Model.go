package model

import (
	"github.com/dronm/gobizap/fields"
)

//Complete model
type Complete_Model struct {	
	//Pattern fields.ValText `json:"pattern" length:500`
	Count fields.ValInt `json:"count" default:10`	
	Ic fields.ValInt `json:"ic" default:1 minValue:0 maxValue:1`
	Mid fields.ValInt `json:"mid" default:1 minValue:0 maxValue:1`
	Ord_directs fields.ValText `json:"ord_directs" length:500`
	Field_sep fields.ValText `json:"field_sep" length:2`
	Cond_fields fields.ValText `json:"cond_fields" length:"1000"`
	Cond_sgns   fields.ValText `json:"cond_sgns" length:"1000"`
	Cond_vals   fields.ValText `json:"cond_vals" length:"1000"`
}
