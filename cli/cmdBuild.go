package main

import (
	"fmt"
)

const CMD_BUILD = "build"

type BuildHandler struct {
}

func (h BuildHandler) OnCommand(args []string) error {
	return nil
}

func (c BuildHandler) Help() {
	formatter := TermColors.OkColor.SprintfFunc()
	fmt.Printf("Command:\t\t\t%s\n", formatter(CMD_BUILD))
	fmt.Println("")
}
