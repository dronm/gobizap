package main

import (
	"fmt"

	"github.com/fatih/color"
)

const CMD_PARAM_IND = 1

func getCmdTextFormatter() func(a ...interface{}) string {
	return color.New(color.FgGreen, color.Bold).SprintFunc()
}

func ShowHelp(args []string) error {
	if len(args) < CMD_PARAM_IND+1 {
		ShowHelpFor_help()
		return nil
	}

	cmd := args[CMD_PARAM_IND]
	switch cmd {
	case COMMAND_INIT:
		ShowHelpFor_init()

	case COMMAND_SQL:
		ShowHelpFor_sql()

	case COMMAND_MINIFY:
		ShowHelpFor_minify()

	case COMMAND_BUILD:
		ShowHelpFor_build()

	case COMMAND_UPLOAD:
		ShowHelpFor_upload()

	case COMMAND_PUSH:
		ShowHelpFor_push()

	default:
		ShowHelpFor_all()
	}

	return nil
}

func ShowHelpFor_all() {
	cmd_col := color.New(color.FgRed).SprintFunc()
	arg_col := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("Usage: <%s> <%s>\n", cmd_col("COMMAND"), arg_col("ARGUMENTS"))
	fmt.Println("Commands:")
	fmt.Printf("		%s\n", cmd_col("help"))
	fmt.Printf("		%s\n", cmd_col("init"))
	fmt.Printf("		%s\n", cmd_col("sql"))
	fmt.Printf("		%s\n", cmd_col("minify"))
	fmt.Printf("		%s\n", cmd_col("build"))
	fmt.Printf("		%s\n", cmd_col("upload"))
	fmt.Printf("		%s\n", cmd_col("push"))
	fmt.Println("")
}

func ShowHelpFor_help() {
	formatter := getCmdTextFormatter()
	fmt.Printf("Command:\t\t\t%s\n", formatter("help"))
	fmt.Printf("Arguments:\t\t\t%s\n", formatter("<COMMAND_NAME>"))
	fmt.Printf("%s command shows help information for a giving command.\n", formatter("help"))
	fmt.Printf("%s is the command you want to get help for.\n", formatter("<COMMAND_NAME>"))
	fmt.Println("")
}

func ShowHelpFor_init() {
	formatter := getCmdTextFormatter()
	fmt.Printf("Command:\t\t\t%s\n", formatter("init"))
	fmt.Printf("Arguments:\t\t\t%s\n", formatter("<APP_NAME>"))
	fmt.Printf("%s command initializes new project in the working directory. All necessary files, folders, symlinks are created.\n", formatter("init"))
	fmt.Printf("%s is the name for a new project. The argument is not obligatory. If not defind, project folder will be used for a name.\n", formatter("<APP_NAME>"))
	fmt.Println("")
}

func ShowHelpFor_sql() {
	formatter := getCmdTextFormatter()
	fmt.Printf("Command:\t\t\t%s\n", formatter("sql"))
	fmt.Printf("Arguments:\t\t\t%s\n", formatter("<up> <down> <add> <pos>"))
	fmt.Printf("%s command manages db migrations.\n", formatter("sql"))
	fmt.Printf("%s runs one up migration.\n", formatter("sql up"))
	fmt.Printf("%s runs all up migrations.\n", formatter("sql up all"))
	fmt.Printf("%s runs one down migration.\n", formatter("sql down"))
	fmt.Printf("%s runs all down migrations.\n", formatter("sql down all"))
	fmt.Printf("%s creates new db migration.\n", formatter("sql add <ACTION>"))
	fmt.Printf("If run with no arguments, sql up migration is assumed.\n")
	fmt.Println("")
}

func ShowHelpFor_build() {
	formatter := getCmdTextFormatter()
	fmt.Printf("Command:\t\t\t%s\n", formatter("build"))
	fmt.Println("")
}

func ShowHelpFor_upload() {
	formatter := getCmdTextFormatter()
	fmt.Printf("Command:\t\t\t%s\n", formatter("upload"))
	fmt.Printf("Arguments:\t\t\t%s\n", formatter("<CONFIG_FILE>"))
	fmt.Printf("%s command requires one argument %s which is the project configuration file. By default project folder name will be used plus file extension .json. If not found, the first file with .json extension will be used. Otherwise error exeption will be thrown.\n",
		formatter("upload"),
		formatter("<CONFIG_FILE>"))
}

func ShowHelpFor_minify() {
	formatter := getCmdTextFormatter()
	fmt.Printf("Command:\t\t\t%s\n", formatter("minify"))
	fmt.Println("")
}

func ShowHelpFor_push() {
	formatter := getCmdTextFormatter()
	fmt.Printf("Command:\t\t\t%s\n", formatter("push"))
	fmt.Println("")
}
