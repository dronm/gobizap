package main

import (
	"fmt"
	"path/filepath"

	"github.com/dronm/gobizap/config"
	"github.com/dronm/gobizap/md"
	"github.com/dronm/sqlmigr"
)

//TODO:add command to show all migrations after current

const CMD_SQL = "sql"
const (
	SQL_CMD_MG_UP    = "up"
	SQL_CMD_MG_UPALL = "upall"
	SQL_CMD_MG_DOWN  = "down"
	SQL_CMD_MG_POS   = "pos"
	SQL_CMD_MG_ADD   = "add"
)

type SQLHandler struct {
}

func (h SQLHandler) OnCommand(args []string) error {
	if len(args) < 1 {
		h.Help()
		return nil
	}
	switch args[0] {
	case SQL_CMD_MG_UP:
		return h.migrateUp()
	case SQL_CMD_MG_UPALL:
		return h.migrateUpAll()
	case SQL_CMD_MG_DOWN:
		return h.migrateDown()
	case SQL_CMD_MG_POS:
		return h.migratePos()
	case SQL_CMD_MG_ADD:
		return h.migrateAdd()
	}
	return nil
}

func (c SQLHandler) Help() {
	formatter := TermColors.OkColor.SprintfFunc()
	fmt.Printf("Command:\t\t\t%s\n", formatter(CMD_SQL))
	fmt.Println("")
}

func (h SQLHandler) migrateUpAll() error {
	return nil
}

func (h SQLHandler) migrateUp() error {
	projDir, err := GetProjectDir()
	if err != nil {
		return fmt.Errorf("GetProjectDir() failed: %v", err)
	}
	mgr := sqlmigr.NewMigrator(filepath.Join(projDir, md.BUILD_DIR, md.SQL_DIR))
	mgrFile, err := mgr.Up()
	if err != nil {
		return err
	}
	configFile, err := GetProjectConfigFile()
	if err != nil {
		return err
	}
	appConf := config.AppConfig{}
	if err := config.ReadConf(configFile, &appConf); err != nil {
		return err
	}
	appConn, err := md.NewDbConn(appConf.Db.Primary)
	if err != nil {
		return err
	}

	mgrScript := mgr.GetMigrFullFileName(sqlmigr.MG_UP, mgrFile.Name)
	if err := ApplySQLScript(appConn, mgrScript); err != nil {
		return err
	}
	return nil
}

func (h SQLHandler) migrateDown() error {
	projDir, err := GetProjectDir()
	if err != nil {
		return fmt.Errorf("GetProjectDir() failed: %v", err)
	}
	mgr := sqlmigr.NewMigrator(filepath.Join(projDir, md.BUILD_DIR, md.SQL_DIR))
	mgrFile, err := mgr.Down()
	if err != nil {
		return err
	}
	configFile, err := GetProjectConfigFile()
	if err != nil {
		return err
	}
	appConf := config.AppConfig{}
	if err := config.ReadConf(configFile, &appConf); err != nil {
		return err
	}
	appConn, err := md.NewDbConn(appConf.Db.Primary)
	if err != nil {
		return err
	}

	mgrScript := mgr.GetMigrFullFileName(sqlmigr.MG_DOWN, mgrFile.Name)
	if err := ApplySQLScript(appConn, mgrScript); err != nil {
		return err
	}
	return nil
}

func (h SQLHandler) migratePos() error {
	projDir, err := GetProjectDir()
	if err != nil {
		return fmt.Errorf("GetProjectDir() failed: %v", err)
	}
	mgr := sqlmigr.NewMigrator(filepath.Join(projDir, md.BUILD_DIR, md.SQL_DIR))
	posFileName, err := mgr.GetLastMigrFileName()
	if err != nil {
		return err
	}
	fmt.Printf("Current postion file: %s\n", TermColors.OkColor.SprintFunc()(posFileName))
	return nil
}

func (h SQLHandler) migrateAdd() error {
	return nil
}
