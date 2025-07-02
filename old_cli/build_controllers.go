package main

import "fmt"

func BuildControllers(md *Metadata, projDir string) error {
	for _, ent := range md.Controllers.Controller {
		if ent.Cmd == nil {
			continue
		}
		switch *ent.Cmd {
		case MD_CMD_DEL:
			BuildController_del()
		case MD_CMD_ALT:
			BuildController_alt()
		case MD_CMD_ADD:
			BuildController_add()
		default:
			fmt.Printf("unknown command '%s' for controller '%s'", *ent.Cmd, ent.ID)
		}
	}
	return nil
}

func BuildController_del() error {
	return nil
}

func BuildController_alt() error {
	return nil
}

func BuildController_add() error {
	return nil
}
