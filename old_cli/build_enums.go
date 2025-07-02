package main

import (
	"fmt"
	"path/filepath"
)

func BuildEnums(md *Metadata, projDir string) error {
	js_scripts_to_add := []string{}
	md_modif := false
	enums_mod := false
	for _, ent := range md.Enums {
		if ent.Cmd == nil || *ent.Cmd == "" {
			continue
		}
		switch *ent.Cmd {
		case MD_CMD_DEL:
			if err := delEnum(md, projDir, ent, &md_modif); err != nil {
				return err
			}
			if !enums_mod {
				enums_mod = true
			}
		case MD_CMD_ALT:
			if err := buildEnum(md, projDir, ent, &js_scripts_to_add); err != nil {
				return err
			}
			if !enums_mod {
				enums_mod = true
			}
		case MD_CMD_ADD:
			if err := buildEnum(md, projDir, ent, &js_scripts_to_add); err != nil {
				return err
			}
			if !enums_mod {
				enums_mod = true
			}
		default:
			fmt.Printf("unknown command '%s' for enum '%s'", *ent.Cmd, ent.ID)
		}
	}
	//App.enums
	if enums_mod {
		values := make([]map[string]interface{}, len(md.Enums))
		for i, ent := range md.Enums {
			descriptions := make([]map[string]interface{}, 0)
			not_first := true
			if i == 0 {
				not_first = false
			}
			for _, en_val := range ent.Values {
				for lang_id, lang_descr := range en_val.LangDescriptions {
					not_first_lang := true
					if len(descriptions) == 0 {
						not_first_lang = false
					}
					descriptions = append(descriptions, map[string]interface{}{"VALUE": en_val.ID,
						"NOT_FIRST":   not_first_lang,
						"LANG":        lang_id,
						"DESCRIPTION": lang_descr,
					})
				}
			}
			values[i] = map[string]interface{}{"ID": ent.ID,
				"NOT_FIRST":    not_first,
				"DESCRIPTIONS": descriptions,
			}
		}
		params := map[string]interface{}{"VALUES": values}
		//add script to metadata
		rel_scr := filepath.Join(CONTROLS_DIR, ENUM_LIST_JS_NAME)
		source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, ENUM_LIST_JS_TMPL_NAME)
		dest := filepath.Join(projDir, BUILD_DIR, JS_DIR, rel_scr)
		if err := RenderTemplate(source, dest, params); err != nil {
			return err
		}
		js_scripts_to_add = append(js_scripts_to_add, rel_scr)
	}

	if len(js_scripts_to_add) > 0 {
		md_modif = md.AddNewJsScripts("build", js_scripts_to_add)
	}
	if md_modif {
		// if err := md.WriteMd(projDir); err != nil {
		// 	return err
		// }
		// LogInfo("build", "metadata file was modified", nil)
	}
	return nil
}

func delEnum(md *Metadata, projDir string, ent Enum, mdModif *bool) error {
	//remove enum file
	//javascript enum file
	js_rel_file_name := filepath.Join(BUILD_DIR, JS_DIR, fmt.Sprintf(ENUM_JS_NAME_TEMPL, ent.ID))
	js_file_name := filepath.Join(projDir, js_rel_file_name)
	if err := DeleteFileIfExists(js_file_name); err != nil {
		return err
	}

	//javascript grid column file
	js_file_name = filepath.Join(projDir, BUILD_DIR, JS_DIR, fmt.Sprintf(ENUM_GR_COL_JS_NAME_TEMPL, ent.ID))
	if err := DeleteFileIfExists(js_file_name); err != nil {
		return err
	}

	//golang enum file
	go_file_name := filepath.Join(projDir, ENUM_DIR, fmt.Sprintf(ENUM_NAME_TMPL, ent.ID))
	if err := DeleteFileIfExists(go_file_name); err != nil {
		return err
	}

	//remome from md
	DeleteJavascript(md, js_rel_file_name)

	*mdModif = true

	return nil
}

func buildEnum(md *Metadata, projDir string, ent Enum, jsScriptsToAdd *[]string) error {

	enum_values := make([]interface{}, len(ent.Values))
	enum_lang_values := make([]interface{}, 0)
	for i, en_val := range ent.Values {
		descriptions := make([]interface{}, len(en_val.LangDescriptions))
		not_first_val := true
		if i == 0 {
			not_first_val = false
		}
		n := 0
		for lang_id, lang_descr := range en_val.LangDescriptions {
			not_first := true
			if n == 0 {
				not_first = false
			}
			descriptions[n] = map[string]interface{}{"VALUE": lang_descr,
				"NOT_FIRST": not_first,
				"LANG":      lang_id,
			}

			enum_lang_values = append(enum_lang_values, map[string]interface{}{"ID": en_val.ID,
				"NOT_FIRST":   not_first_val,
				"FIRST":       !not_first_val,
				"DESCRIPTION": lang_descr,
				"LANG":        lang_id,
			})
		}
		enum_values[i] = map[string]interface{}{"ID": en_val.ID,
			"NOT_FIRST":    not_first_val,
			"FIRST":        !not_first_val,
			"DESCRIPTIONS": descriptions,
		}

	}
	params := map[string]interface{}{"ID": ent.ID,
		"VALUES":      enum_values,
		"LANG_VALUES": enum_lang_values,
	}

	if *ent.Cmd == MD_CMD_ADD || *ent.Cmd == MD_CMD_ALT {
		//golang file
		source := filepath.Join(projDir, BUILD_DIR, TMPL_DIR, TMPL_ENUM)
		dest := filepath.Join(projDir, ENUM_DIR, fmt.Sprintf(ENUM_NAME_TMPL, ent.ID))
		if err := RenderTemplate(source, dest, params); err != nil {
			return err
		}
		//format new file with formatter
		if err := FormatGoFile(dest); err != nil {
			return err
		}

		//javascript files: enum && grid column
		//add script to metadata
		rel_scr := filepath.Join(ENUM_DIR, fmt.Sprintf(ENUM_JS_NAME_TEMPL, ent.ID))
		*jsScriptsToAdd = append(*jsScriptsToAdd, rel_scr)

		rel_scr_gr_col := filepath.Join(ENUM_DIR, fmt.Sprintf(ENUM_GR_COL_JS_NAME_TEMPL, ent.ID))
		*jsScriptsToAdd = append(*jsScriptsToAdd, rel_scr_gr_col)

	}
	return nil
}
