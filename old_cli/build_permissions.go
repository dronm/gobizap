package main

import "fmt"

func BuildPermissions(md *Metadata, projDir string) error {
	for _, ent := range md.Permissions {
		if ent.Cmd == nil {
			continue
		}
		switch *ent.Cmd {
		case MD_CMD_DEL, MD_CMD_ALT, MD_CMD_ADD:
			BuildPermission_update()
			return nil
		default:
			fmt.Printf("unknown command '%s' for permission '%s:%s.%s.%s'", *ent.Cmd,
				ent.Type, ent.RoleID, ent.ControllerID, ent.MethodID)
		}
	}
	return nil
}

func BuildPermission_update() error {
	return nil
}
