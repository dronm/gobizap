package menu

//Controller method model
import (
	"github.com/dronm/gobizap/fields"
)

//
type VariantStorage_get_filter_data struct {
	Storage_name fields.ValText `json:"storage_name" required:"true"`
	Variant_name fields.ValText `json:"variant_name"`
}

type VariantStorage_get_filter_data_argv struct {
	Argv *VariantStorage_get_filter_data `json:"argv"`	
}

