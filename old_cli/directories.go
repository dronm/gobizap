package main

//This files contains constants with predefined directory names.

const (
	//project directories
	BUILD_DIR    = "src"
	SQL_DIR      = "sql"
	CONTR_DIR    = "controllers"
	MODEL_DIR    = "models"
	ENUM_DIR     = "enums"
	CONSTANT_DIR = "constants"
	CONTROLS_DIR = "custom_controls"
	CSS_DIR      = "css"
	JS_DIR       = "js"
	UPDATES_DIR  = "updates"

	SQL_MIG_DOWN_DIR = "mig_down"
	SQL_MIG_UP_DIR   = "mig_up"

	MD_FILE_NAME = "metadata.xml"

	//gobizap directories
	PROJ_DIR = "project"
	TMPL_DIR = "templates"

	//template names
	TMPL_CREATE_DB_UP   = "db_create_up.sql.tmpl"
	TMPL_CREATE_DB_DOWN = "db_create_down.sql.tmpl"
	TMPL_INIT_DB_UP     = "db_init_up.sql.tmpl"
	TMPL_INIT_DB_DOWN   = "db_init_down.sql.tmpl"
	TMPL_CREATE_TB_UP   = "table_create_up.sql.tmpl"
	TMPL_CREATE_TB_DOWN = "table_create_down.sql.tmpl"
	TMPL_ALT_TB_UP      = "table_alt_up.sql.tmpl"
	TMPL_ALT_TB_DOWN    = "table_alt_down.sql.tmpl"

	TMPL_CREATE_IND_UP   = "create_ind_up.sql.tmpl"
	TMPL_CREATE_IND_DOWN = "create_ind_down.sql.tmpl"
	TMPL_DROP_IND_UP     = "drop_ind_up.sql.tmpl"
	TMPL_DROP_IND_DOWN   = "drop_ind_down.sql.tmpl"

	MIG_POS_FILE = "current_mig.pos"

	NEW_FILE_PERMIS = 0600

	TMPL_MODEL    = "Model.go.tmpl"
	TMPL_MODEL_JS = "Model.js.tmpl"
	TMPL_CONST    = "Constant.go.tmpl"
	TMP_CONTR     = "Controller.go.tmpl"
	TMP_CONTR_IMP = "ControllerImp.go.tmpl"
	TMPL_ENUM     = "Enum.go.tmpl"
	TMPL_PM_COMPL = "PublicMethodCompleteModel.go.tmpl"
	TMPL_PM_MODEL = "PublicMethodModel.go.tmpl"
	TMPL_PERM     = "permissons.sql.tmpl"

	MODEL_NAME_TMPL     = "%s_Model.go"
	MODEL_JS_NAME_TEMPL = "%s_Model.js"

	ENUM_NAME_TMPL            = "Enum_%s.go"
	ENUM_JS_NAME_TEMPL        = "Enum_%s.js"
	ENUM_GR_COL_JS_NAME_TEMPL = "EnumGridColumn_%s.js"
	ENUM_LIST_JS_NAME         = "App.enums.js"
	ENUM_LIST_JS_TMPL_NAME    = "App.enums.js.tmpl"

	CONSTANT_NAME_TMPL = "Constant_%s.go"

	DT_ENUM  = "Enum"
	DT_JSON  = "JSON"
	DT_JSONB = "JSONB"
)
