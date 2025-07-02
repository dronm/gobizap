package main

import (
	"fmt"
)

const CMD_MINIFY = "minify"

type MinifyHandler struct {
}

func (h MinifyHandler) OnCommand(args []string) error {
	return nil
}

func (c MinifyHandler) Help() {
	formatter := TermColors.OkColor.SprintfFunc()
	fmt.Printf("Command:\t\t\t%s\n", formatter(CMD_MINIFY))
	fmt.Println("")
}
