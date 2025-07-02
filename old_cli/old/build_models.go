package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func BuildModels(md *Metadata, projDir string) error {
	for _, ent := range md.Models.Model {
		if ent.Cmd == nil {
			continue
		}
		switch *ent.Cmd {
		case MD_CMD_DEL:
			BuildModel_del()
		case MD_CMD_ALT:
			BuildModel_alt(md, projDir, ent)
		case MD_CMD_ADD:
			BuildModel_add()
		default:
			fmt.Printf("unknown command '%s' for enum '%s'", *ent.Cmd, ent.ID)
		}
	}
	return nil
}

func BuildModel_del() error {
	return nil
}

func BuildModel_alt(md *Metadata, projDir string, ent Model) error {
	schema := md.DataSchema
	if ent.DataSchema != "" {
		schema = ent.DataSchema
	}
	var base_model Model
	if ent.BaseModelID != "" {
		base_model_id := ""
		if strings.HasSuffix(ent.BaseModelID, "List") {
			base_model_id = strings.CutSuffix(ent.BaseModelID, "List")

		} else if strings.HasSuffix(ent.BaseModelID, "Dialog") {
			base_model_id = strings.CutSuffix(ent.BaseModelID, "Dialog")
		}
		if base_model_id != "" {
			for _, bm := range md.Models.Model {
				if bm.ID == base_model_id {
					base_model = bm
					break
				}
			}
		}
	}
	params := map[string]interface{}{"ID": ent.ID,
		"OBJECT_DATA_TABLE": ent.DataTable,
		"DATA_SCHEMA":       schema,
		"FIELDS":            map[string]string{},
		"VIRTUAL":           ent.Virtual,
		"NOT_VIRTUAL":       !ent.Virtual,
		"KEYS":              map[string]string{},
		"ENUMS_EXIST":       false,
		"APP_NAME":          md.AppName,
		"BASE_DATA_TABLE":   base_model.DataTable,
	}

	params["AGG_FUNCTIONS_EXIST"] = true // + total
	total_count_exists := true
	agg_count := len(ent.AggFunctions)
	for i, agg := range ent.AggFunctions {
		if agg.Alias == "totalCount" {
			total_count_exists = true
			break
		}
	}
	if !total_count_exists {
		agg_count++
	}
	agg_funcs := make([]map[string]interface{}, agg_count)
	if agg_count > 0 {
		for i, agg := range ent.AggFunctions {
			not_first := false
			if i > 0 {
				not_first = true
			}
			tmpl_agg := map[string]interface{}{"NOT_FIRST": not_first,
				"ALIAS": agg.Alias,
				"EXPR":  agg.Expr,
			}
			agg_funcs[i] = tmpl_agg
		}
	}

	//add total
	if !total_count_exists {
		tmpl_agg := map[string]interface{}{"NOT_FIRST": not_first,
			"ALIAS": "totalCount",
			"EXPR":  "count(*)",
		}
		agg_funcs[agg_count-1] = tmpl_agg
	}
	params["AGG_FUNCTIONS"] = agg_funcs

	if ent.LimitCount != "" && ent.DocPerPageCount != "" {
		params["LIMIT_COUNT"] = ent.LimitCount
		params["DOC_PER_PAGE_COUNT"] = ent.DocPerPageCount
	}

	if len(ent.DefaultOrder) > 0 {
		params["ORDERS"] = ""
		var order_v strings.Builder
		order_v.Grow(len(ent.DefaultOrder))
		for i, order_f := range ent.DefaultOrder {
			if i > 0 {
				order_v.WriteString(",")
			}
			s_dir := order_f.Field.SortDirect
			if s_dir == "" {
				s_dir = "ASC"
			}

			order_v.WriteString(order_f.Field.ID + " " + s_dir)
		}
	}

	//field
	// add_field_exists := false
	for i, f := range ent.Field {
		enum_exists := false
		if f.DataType == DT_ENUM {
			enum_exists = true
			params["ENUMS_EXIST"] = true
			go_type := golangDataType(f.DataType)
			not_first := false
			if i > 0 {
				not_first = false
			}
			f_params := map[string]interface{}{
				"NAME":          strings.ToUpper(f.ID[0:1]) + f.ID[1:],
				"ID":            f.ID,
				"TYPE":          f.DataType,
				"NOT_FIRST":     not_first,
				"NOT_REF_TABLE": true,
				"SYS_COL":       f.SysCol,
				"NAME_ALIAS":    f.Alias,
			}
		}
	}

	source := filepath.Join(BUILD_DIR, TMPL_DIR, TMPL_MODEL)
	dest := filepath.Join(projDir, MODEL_DIR, fmt.Sprintf(MODEL_NAME_TMPL))
	if err := RenderTemplate(source, dest, params); err != nil {
		return err
	}
	return nil
}

func golangDataType(mdDataType string) string {
	return ""
}

func BuildModel_add() error {
	return nil
}
