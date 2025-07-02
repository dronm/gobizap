package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dronm/gobizap/md"
	"github.com/fatih/color"

	"github.com/dronm/sqlmigr"
)

type TerminalColors struct {
	CommentColor *color.Color
	OkColor      *color.Color
	ErrorColor   *color.Color
}

var TermColors = TerminalColors{CommentColor: color.New(color.FgYellow),
	OkColor:    color.New(color.FgGreen),
	ErrorColor: color.New(color.FgRed),
}

type SQL struct {
}

func (s SQL) GetProjectDbConn() (*md.DbConn, error) {

	confFile, err := md.GetProjectConfigFile()
	if err != nil {
		return nil, err
	}
	md.
}

func (s SQL) MigrateUp() error {
	projDir, err := md.GetProjectDir()
	if err != nil {
		return fmt.Errorf("md.GetProjectDir() failed: %v", err)
	}

	mgr := sqlmigr.NewMigrator(filepath.Join(projDir, md.BUILD_DIR, md.SQL_DIR))
	mgrFile, err := mgr.Up()
	if err != nil {
		return err
	}
	mgrScript := mgr.GetMigrFullFileName(sqlmigr.MG_UP, mgrFile.Name)
	//create db migration
	md.NewDbConn()
	return nil
}
func (s SQL) MigrateDown() error {
	return nil
}
func (s SQL) ShowPosition() error {
	return nil
}

func main() {
	if len(os.Args) < 2 {
		ShowHelp()
		os.Exit(1)
	}

	sql := SQL{}
	cmd := os.Args[1]
	switch cmd {
	case "up":
		if err := sql.MigrateUp(); err != nil {
			TermColors.ErrorColor.Println(err)
			os.Exit(-1)
		}
	case "down":
		if err := sql.MigrateDown(); err != nil {
			TermColors.ErrorColor.Println(err)
			os.Exit(-1)
		}
	case "pos":
		if err := sql.ShowPosition(); err != nil {
			TermColors.ErrorColor.Println(err)
			os.Exit(-1)
		}
	}
}

func ShowHelp() {

}
