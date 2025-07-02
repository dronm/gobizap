package main

import (
	"fmt"
)

const CMD_HELP = "help"

type HelpHandler struct {
}

func (h HelpHandler) OnCommand(args []string) error {
	if len(args) > 0 {
		return h.ShowForCmd(args[0])
	}
	//common usage
	h.ShowForAll()

	return nil
}

func (h HelpHandler) ShowForCmd(cmdID string) error {
	cmd, ok := Commands[cmdID]
	if !ok {
		return fmt.Errorf("Command '%s' not found", cmdID)
	}
	cmd.Help()

	return nil
}

func (h HelpHandler) ShowForAll() {
	fmt.Printf("Usage: <%s> <%s>\n",
		TermColors.OkColor.SprintFunc()("COMMAND"),
		TermColors.OkColor.SprintFunc()("ARGUMENTS"),
	)
	for id, cmd := range Commands {
		if id == CMD_HELP {
			continue
		}
		cmd.Help()
	}
}

func (c HelpHandler) Help() {
	//stub
}
