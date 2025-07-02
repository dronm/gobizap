package main

import (
	"fmt"
)

const CMD_UPLOAD = "upload"

type UploadHandler struct {
}

func (h UploadHandler) OnCommand(args []string) error {
	return nil
}

func (c UploadHandler) Help() {
	formatter := TermColors.OkColor.SprintfFunc()
	fmt.Printf("Command:\t\t\t%s\n", formatter(CMD_UPLOAD))
	fmt.Println("")
}
