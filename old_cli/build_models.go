package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dronm/sqlmigr"
)

const (
	SORT_UP   = "ASC"
	SORT_DOWN = "DESC"
)

func BuildModels(md *Metadata, projDir string) error {
	js_scripts_to_add := []string{}
	md_modif := false
	for _, ent := range md.Models.Model {
		if ent.Cmd == nil {
			continue
		}
		switch *ent.Cmd {
		case MD_CMD_DEL:
			if err := delModel(md, projDir, ent, &md_modif); err != nil {
				return err
			}
			LogInfo("build", "deleted model '%s'", ent.ID)
		case MD_CMD_ALT:
			if err := buildModel(md, projDir, ent, &js_scripts_to_add); err != nil {
				return err
			}
			LogInfo("build", "altered model '%s'", ent.ID)
		case MD_CMD_ADD:
			if err := buildModel(md, projDir, ent, &js_scripts_to_add); err != nil {
				return err
			}
			LogInfo("build", "added new model '%s' ", ent.ID)
		default:
			LogWarning("build", "unknown command '%s' for model '%s'", *ent.Cmd, ent.ID)
		}
	}
	if len(js_scripts_to_add) > 0 {
		md_modif = md.AddNewJsScripts("build", js_scripts_to_add)
	}
	if md_modif {
		if err := md.WriteMd(projDir); err != nil {
			return err
		}
		// LogInfo("build", "metadata file was modified", nil)
	}
	return nil
}

func delModel(md *Metadata, projDir string, ent Model, mdModif *bool) error {
	//get template parameters
	params, err := ModelTemplateParams(md, projDir, ent, nil)
	if err != nil {
		return err
	}

	//remove model file
	js_rel_file_name := filepath.Join(BUILD_DIR, JS_DIR, fmt.Sprintf(MODEL_JS_NAME_TEMPL, ent.ID))
	js_file_name := filepath.Join(projDir, js_rel_file_name)
	if err := DeleteFileIfExists(js_file_name); err != nil {
		return err
	}

	//database migration
	mgr := sqlmigr.NewMigrator(filepath.Join(projDir, BUILD_DIR, SQL_DIR))
	create_db_fname := mgr.GetMigrFileName(create_mig_t, "createDB", sqlmigr.MG_UP)
	create_db_up_mgr := mgr.GetMigrFullFileName(sqlmigr.MG_UP, create_db_fname)
	mgr.GetMigrFileName()
	err := ApplyTemplateFromFile(mig_up_source, mig_up_dest, *params)
	if err != nil {
		return err
	}

	mig_name := dbMigrationName(time.Now(), ent.DataTable+"drop")
	mig_up_source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_CREATE_TB_DOWN)
	mig_up_dest := filepath.Join(projDir, BUILD_DIR, SQL_DIR, SQL_MIG_UP_DIR, mig_name)
	if err := RenderTemplate(mig_up_source, mig_up_dest, *params); err != nil {
		return err
	}
	mig_down_source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_CREATE_TB_UP)
	mig_down_dest := filepath.Join(projDir, BUILD_DIR, SQL_DIR, SQL_MIG_DOWN_DIR, mig_name)
	if err := RenderTemplate(mig_down_source, mig_down_dest, *params); err != nil {
		return err
	}
	LogInfo("build", "added database migration %s", mig_name)

	//remome from md
	DeleteJavascript(md, js_rel_file_name)

	*mdModif = true

	return nil
}

// ModelTemplateParams() returns params to be used in a template to create/drop/alter a model
func ModelTemplateParams(md *Metadata, projDir string, ent Model, jsScriptsToAdd *[]string) (*map[string]interface{}, error) {
	schema := md.DataSchema
	if ent.DataSchema != "" {
		schema = ent.DataSchema
	}

	type ref_field struct {
		fieldId    string
		refFieldId string
	}
	var base_model Model
	var base_model_ref_fields map[string]ref_field //key is a ref table for search
	if ent.BaseModelID != "" {
		base_model_id := ""
		base_model_found := false // base model found
		base_model_id, base_model_found = strings.CutSuffix(ent.BaseModelID, "List")
		if !base_model_found {
			base_model_id, base_model_found = strings.CutSuffix(ent.BaseModelID, "Dialog")
		}
		if base_model_found {
			//find base model
			for _, bm := range md.Models.Model {
				if bm.ID != base_model_id {
					continue
				}
				base_model = bm
				//iterate over fields, get all ref fields
				for _, base_f := range base_model.Field {
					if base_f.RefTable != "" && base_f.RefField != "" {
						base_model_ref_fields[base_f.RefTable] = ref_field{fieldId: base_f.ID,
							refFieldId: base_f.RefField,
						}
					}
				}
			}
		}
	}
	mod_name, err := GetProjectModuleName(projDir)
	if err != nil {
		return nil, err
	}
	//template parameters
	params := map[string]interface{}{"ID": ent.ID,
		"OBJECT_DATA_TABLE": ent.DataTable,
		"DATA_SCHEMA":       schema,
		"DATA_OWNER":        md.Owner,
		"FIELDS":            map[string]string{},
		"VIRTUAL":           ent.Virtual,
		"NOT_VIRTUAL":       !bool(*ent.Virtual),
		"KEYS":              map[string]string{},
		"ENUMS_EXIST":       false,
		"APP_MODULE":        mod_name,
		"BASE_DATA_TABLE":   base_model.DataTable,
	}
	params["AGG_FUNCTIONS_EXIST"] = true // + total
	total_count_exists := true
	agg_count := len(ent.AggFunctions)
	for _, agg := range ent.AggFunctions {
		if agg.Alias == "totalCount" {
			total_count_exists = true
			break
		}
	}
	if !total_count_exists {
		agg_count++
	}
	agg_funcs := make([]map[string]interface{}, agg_count)
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

	//add total
	if !total_count_exists {
		not_first := false
		if len(ent.AggFunctions) > 0 {
			not_first = true
		}
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
		var order_v strings.Builder
		order_v.Grow(len(ent.DefaultOrder))
		for i, order_f := range ent.DefaultOrder {
			if i > 0 {
				order_v.WriteString(",")
			}
			s_dir := order_f.Field.SortDirect
			if s_dir == "" || (s_dir != SORT_UP && s_dir != SORT_DOWN) {
				LogWarning("build", "unknown sort order for model '%s', default sort order is used", ent.ID)
				s_dir = SORT_UP
			}
			order_v.WriteString(order_f.Field.ID + " " + s_dir)
		}
		params["ORDERS"] = order_v.String()
	}

	fields := make([]map[string]interface{}, len(ent.Field))
	//field
	enum_id := ""
	for i, f := range ent.Field {
		if f.DataType == DT_ENUM {
			enum_id = f.ID
		}
		go_type, err := GolangDataType(ent, f)
		if err != nil {
			return nil, err
		}
		sql_type, err := SQLDataType(ent, f, f.DataType)
		if err != nil {
			return nil, err
		}
		not_first := false
		if i > 0 {
			not_first = true
		}
		f_params := map[string]interface{}{
			"NAME":          strings.ToUpper(f.ID[0:1]) + f.ID[1:],
			"ID":            f.ID,
			"TYPE":          go_type,
			"NOT_FIRST":     not_first,
			"NOT_REF_TABLE": true,
			"SYS_COL":       f.SysCol,
			"SQL_DATA_TYPE": sql_type,
		}
		if f.Alias != "" {
			f_params["ALIAS"] = f.Alias
		}

		if *ent.Cmd == MD_CMD_ADD || f.Cmd != nil && *f.Cmd == MD_CMD_ADD {
			f_params["ADD"] = true

		} else if f.Cmd != nil && *f.Cmd == MD_CMD_DEL {
			f_params["DROP"] = true

		} else if f.Cmd != nil && *f.Cmd == MD_CMD_ALT {
			if f.OldDataType != "" {
				//new data type
				f_params["ALTER_TYPE"] = true

			} else if f.OldID != "" {
				f_params["RENAME"] = true

			}
		}

		if f.OldDataType != "" {
			sql_type, err := SQLDataType(ent, f, f.OldDataType)
			if err != nil {
				return nil, err
			}
			f_params["OLD_SQL_DATA_TYPE"] = sql_type
		}
		if f.OldID != "" {
			f_params["OLD_ID"] = f.OldID
		}
		//ref table if json type
		if f.DataType == DT_JSON || f.DataType == DT_JSONB {
			if ref_table_id, found := strings.CutSuffix(f.ID, "_ref"); found {
				if ref_table, ok := base_model_ref_fields[ref_table_id]; ok {
					f_params["REF_TABLE"] = ref_table_id
					f_params["REF_FIELD"] = ref_table.refFieldId
					f_params["NOT_REF_TABLE"] = false
					f_params["BASE_ID"] = ref_table.fieldId

					params["REFS_EXIST"] = true
				}
			}
		}
		SetTemplateFieldParams(f_params, f, ent.DefaultOrder)
		if f.PrimaryKey {
			params["KEYS"] = map[string]interface{}{"ID": f_params["ID"],
				"NAME": f_params["NAME"],
				"TYPE": f_params["TYPE"],
			}
		}
		fields[i] = f_params
	}
	params["FIELDS"] = fields
	if enum_id != "" {
		params["ENUMS_EXIST"] = true
	}

	//add indexes
	if *ent.Cmd != MD_CMD_DEL && len(ent.Index) > 0 {
		indexes_len := 0
		if *ent.Cmd == MD_CMD_ADD {
			//all indexes
			indexes_len = len(ent.Index)
		} else {
			for _, ind := range ent.Index {
				if ind.Cmd != nil && (*ind.Cmd == MD_CMD_ADD || *ind.Cmd == MD_CMD_ALT || *ind.Cmd == MD_CMD_DEL) {
					indexes_len++
				}
			}
		}
		if indexes_len > 0 {
			indexes := make([]map[string]interface{}, indexes_len)
			for i, ind := range ent.Index {
				if *ent.Cmd == MD_CMD_ALT && ind.Cmd == nil {
					//alter model, no cmd on index - skeep
					continue
				}

				ind_type := ""
				if ind.Type != "" {
					ind_type = ind.Type
				} else {
					ind_type = "bree" //default index type
				}

				var index_fields []map[string]interface{}
				var index_expr string
				if ind.Expr != "" {
					index_expr = ind.Expr

				} else {
					index_fields = make([]map[string]interface{}, len(ind.Field))
					for n, ind_f := range ind.Field {
						not_first := false
						if n > 0 {
							not_first = true
						}
						order := ind_f.Order
						if order == "" {
							order = "ASC"
						}
						nulls := ind_f.Nulls
						if nulls == "" {
							nulls = "LAST"
						}
						index_fields[n] = map[string]interface{}{"NOT_FIRST": not_first,
							"ID":    ind_f.ID,
							"ORDER": order,
							"NULLS": nulls,
						}
					}
				}
				indexes[i] = map[string]interface{}{"UNIQUE": ind.Unique,
					"ID":           ind.ID,
					"TYPE":         ind_type,
					"ADD":          *ent.Cmd == MD_CMD_ADD || (ind.Cmd != nil && (*ind.Cmd == MD_CMD_ADD || *ind.Cmd == MD_CMD_ALT)),
					"DROP":         *ent.Cmd == MD_CMD_ALT || (ind.Cmd != nil && (*ind.Cmd == MD_CMD_ALT || *ind.Cmd == MD_CMD_DEL)),
					"INDEX_FIELDS": index_fields,
					"EXPR":         index_expr,
				}
			}
			params["INDEXES"] = indexes
		}
	}
	return &params, nil
}

// BuildGolangModel() builds: model file, migrations
func buildModel(md *Metadata, projDir string, ent Model, jsScriptsToAdd *[]string) error { //get template parameters
	params, err := ModelTemplateParams(md, projDir, ent, jsScriptsToAdd)
	if err != nil {
		return err
	}

	//migration based on command
	var mig_name string
	if *ent.Cmd == MD_CMD_ADD {
		mig_name = dbMigrationName(time.Now(), ent.DataTable+"_add")
		mig_up_source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_CREATE_TB_UP)
		mig_up_dest := filepath.Join(projDir, BUILD_DIR, SQL_DIR, SQL_MIG_UP_DIR, mig_name)
		if err := RenderTemplate(mig_up_source, mig_up_dest, *params); err != nil {
			return err
		}
		mig_down_source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_CREATE_TB_DOWN)
		mig_down_dest := filepath.Join(projDir, BUILD_DIR, SQL_DIR, SQL_MIG_DOWN_DIR, mig_name)
		if err := RenderTemplate(mig_down_source, mig_down_dest, *params); err != nil {
			return err
		}

	} else if *ent.Cmd == MD_CMD_ALT {
		mig_name = dbMigrationName(time.Now(), ent.DataTable+"_alt")
		mig_up_source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_ALT_TB_UP)
		mig_up_dest := filepath.Join(projDir, BUILD_DIR, SQL_DIR, SQL_MIG_UP_DIR, mig_name)
		if err := RenderTemplate(mig_up_source, mig_up_dest, *params); err != nil {
			return err
		}
		mig_down_source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_ALT_TB_DOWN)
		mig_down_dest := filepath.Join(projDir, BUILD_DIR, SQL_DIR, SQL_MIG_DOWN_DIR, mig_name)
		if err := RenderTemplate(mig_down_source, mig_down_dest, *params); err != nil {
			return err
		}

	} else if *ent.Cmd == MD_CMD_DEL {
		//reverse migration
		mig_name = dbMigrationName(time.Now(), ent.DataTable+"drop")
		mig_up_source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_CREATE_TB_UP)
		mig_up_dest := filepath.Join(projDir, BUILD_DIR, SQL_DIR, SQL_MIG_UP_DIR, mig_name)
		if err := RenderTemplate(mig_up_source, mig_up_dest, *params); err != nil {
			return err
		}
		mig_down_source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_CREATE_TB_DOWN)
		mig_down_dest := filepath.Join(projDir, BUILD_DIR, SQL_DIR, SQL_MIG_DOWN_DIR, mig_name)
		if err := RenderTemplate(mig_down_source, mig_down_dest, *params); err != nil {
			return err
		}
	}
	if mig_name != "" {
		LogInfo("build", "added database migration %s", mig_name)
	}

	//rewrite model template in case of add/alter
	if *ent.Cmd == MD_CMD_ADD || *ent.Cmd == MD_CMD_ALT {
		source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_MODEL)
		dest := filepath.Join(projDir, MODEL_DIR, fmt.Sprintf(MODEL_NAME_TMPL, ent.ID))
		if err := RenderTemplate(source, dest, *params); err != nil {
			return err
		}
		//add script to metadata
		rel_scr := filepath.Join(MODEL_DIR, fmt.Sprintf(MODEL_JS_NAME_TEMPL, ent.ID))
		*jsScriptsToAdd = append(*jsScriptsToAdd, rel_scr)

		//format new file with formatter
		if err := FormatGoFile(dest); err != nil {
			return err
		}
	}
	return nil
}

// setFieldsParams set field template parameters
func SetTemplateFieldParams(fParams map[string]interface{}, field Field, defOrder []DefaultOrder) {
	if field.Alias != "" {
		fParams["ALIAS"] = field.Alias
	}
	if field.Length != "" {
		fParams["LEN"] = field.Length
	}
	if field.Precision != "" {
		fParams["PREC"] = field.Precision
	}
	fParams["REC"] = field.Required
	fParams["PK"] = field.PrimaryKey
	fParams["AI"] = field.AutoInc
	for _, ord := range defOrder {
		fParams["ORD"] = ord.Field.SortDirect
		break
	}
}
