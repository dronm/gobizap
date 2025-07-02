package main

import (
	"fmt"
	"path/filepath"
)

func BuildConstants(md *Metadata, projDir string) error {
	for _, ent := range md.Constants.Constant {
		if ent.Cmd == nil {
			continue
		}
		switch *ent.Cmd {
		case MD_CMD_DEL:
			if err := delConstant(md, projDir, ent); err != nil {
				return err
			}
		case MD_CMD_ALT:
			if err := buildConstant(md, projDir, ent); err != nil {
				return err
			}
		case MD_CMD_ADD:
			if err := buildConstant(md, projDir, ent); err != nil {
				return err
			}
		default:
			fmt.Errorf("unknown command '%s' for constant '%s'", *ent.Cmd, ent.ID)
		}
	}
	return nil
}

func delConstant(md *Metadata, projDir string, ent Constant) error {
	//database migration

	//golang file
	go_file_name := filepath.Join(projDir, CONSTANT_DIR, fmt.Sprintf(CONSTANT_NAME_TMPL, ent.ID))
	if err := DeleteFileIfExists(go_file_name); err != nil {
		return err
	}
	return nil
}

func buildConstant(md *Metadata, projDir string, ent Constant) error {
	// params := map[string]interface{}{"ID": ent.ID,
	// 	DATA_TYPE:  $ct_tp,
	// 	FIELDS_MODULE: TRUE,
	// 	AUTOLOAD: ((string)$constant->attributes()->autoload=='TRUE')? 'true':'false'
	// }
	return nil
}
