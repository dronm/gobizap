package main

import (
	"fmt"
)

const CMD_PUSH = "push"

type PushHandler struct {
}

func (h PushHandler) OnCommand(args []string) error {
	return nil
}

func (c PushHandler) Help() {
	formatter := TermColors.OkColor.SprintfFunc()
	fmt.Printf("Command:\t\t\t%s\n", formatter(CMD_PUSH))
	fmt.Println("")
}
