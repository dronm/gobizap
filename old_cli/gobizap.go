package main

import (
	"os"

	"github.com/fatih/color"
)

const (
	COMMAND_PARAM_IND = 1

	//all possible commands
	COMMAND_HELP   = "help"
	COMMAND_INIT   = "init"
	COMMAND_SQL    = "sql"
	COMMAND_MINIFY = "minify"
	COMMAND_BUILD  = "build"
	COMMAND_UPLOAD = "upload"
	COMMAND_PUSH   = "push"
)

type commandHandlerProto = func([]string) error

func main() {
	if len(os.Args) < 2 {
		ShowHelpFor_all()
		os.Exit(1)
	}

	var handler commandHandlerProto

	command := os.Args[COMMAND_PARAM_IND]
	switch command {
	case COMMAND_HELP:
		handler = ShowHelp

	case COMMAND_INIT:
		handler = InitNewApp

	case COMMAND_SQL:
		handler = RunSQL

	case COMMAND_MINIFY:
		handler = RunMinify

	case COMMAND_BUILD:
		handler = Build

	case COMMAND_UPLOAD:
		handler = Upload

	case COMMAND_PUSH:
		handler = Push

	default:
		color.Red("unknown command: %s", command)
		return
	}

	if err := handler(os.Args[COMMAND_PARAM_IND:]); err != nil {
		color.Red("Error executing command: %v\n", err)
		os.Exit(0)
	}
}
