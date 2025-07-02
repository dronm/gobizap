package main

import "fmt"

func BuildScripts(md *Metadata, projDir string) error {
	for _, ent := range md.JSScripts {
		if ent.Cmd == nil {
			continue
		}
		switch *ent.Cmd {
		case MD_CMD_DEL:
			BuildScript_del()
		case MD_CMD_ALT:
			BuildScript_alt()
		case MD_CMD_ADD:
			BuildScript_add()
		default:
			fmt.Printf("unknown command '%s' for scrpt '%s'", *ent.Cmd, ent.File)
		}
	}
	return nil
}

func BuildScript_del() error {
	return nil
}

func BuildScript_alt() error {
	return nil
}

func BuildScript_add() error {
	return nil
}
