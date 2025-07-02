package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const (
	COMMAND_PARAM_IND = 1
)

type commandHandlerProto = func([]string) error

type TerminalColors struct {
	InfoColor   *color.Color
	OkColor     *color.Color
	DangerColor *color.Color
}

var TermColors = TerminalColors{InfoColor: color.New(color.FgYellow),
	OkColor:     color.New(color.FgGreen),
	DangerColor: color.New(color.FgRed),
}

type Command interface {
	OnCommand(args []string) error
	Help()
}

var Commands = map[string]Command{CMD_HELP: HelpHandler{},
	CMD_INIT:   InitHandler{},
	CMD_PUSH:   PushHandler{},
	CMD_MINIFY: MinifyHandler{},
	CMD_SQL:    SQLHandler{},
	CMD_BUILD:  BuildHandler{},
	CMD_UPLOAD: UploadHandler{},
}

func main() {
	if len(os.Args) < COMMAND_PARAM_IND+1 {
		Commands[CMD_HELP].(HelpHandler).ShowForAll()
		os.Exit(1)
	}

	cmdID := os.Args[COMMAND_PARAM_IND]
	cmd, ok := Commands[cmdID]
	if !ok {
		fmt.Printf("Command '%s' not found\n", TermColors.DangerColor.SprintFunc()(cmdID))
		os.Exit(1)
	}
	if err := cmd.OnCommand(os.Args[COMMAND_PARAM_IND+1:]); err != nil {
		fmt.Printf("Error executing command %s: %s\n", cmdID, TermColors.DangerColor.SprintfFunc()("%v", err))
	}
}
