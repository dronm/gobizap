package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dronm/gobizap/md"
	"github.com/dronm/sqlmigr"
	"github.com/fatih/color"
)

var TermColors = TerminalColors{CommentColor: color.New(color.FgYellow),
	PromptColor: color.New(color.FgGreen),
	ErrorColor:  color.New(color.FgRed),
}

const (
	INIT_DB_SUPERUSER    = "postgres"
	INIT_DB_SCHEMA       = "public"
	INIT_DB_HOST         = "localhost:5432"
	INIT_USER_PWD        = "123456"
	INIT_USER_ROLE_NAME  = "admin"
	INIT_USER_ROLE_DESCR = "Administrator"

	PROJ_DIR = "project" //project templates/files directory name

	TMPL_CREATE_DB_UP   = "db_create_up.sql.tmpl"
	TMPL_CREATE_DB_DOWN = "db_create_down.sql.tmpl"
	TMPL_INIT_DB_UP     = "db_init_up.sql.tmpl"
	TMPL_INIT_DB_DOWN   = "db_init_down.sql.tmpl"
)

type ProjGit struct {
	Name  string `json:"name"`
	Descr string `json:"descr"`
}

type ProjDb struct {
	Host          string `json:"host"`
	Schema        string `json:"schema"`
	Name          string `json:"name"`
	UserName      string `json:"userName"`
	UserPwd       string `json:"UserPwd"`
	Spacial       bool   `json:"spacial"`
	SuperuserName string `json:"superuserName"`
}

type InitTimeZone struct {
	Name   string `json:"name"`
	Descr  string `json:"descr"`
	Offset string `json:"offset"`
}

// Application configuration structure
type ProjConfig struct {
	ID string `json:"id"`
	//AuthorName    string
	//TechEmail     string
	Git        ProjGit      `json:"git"`
	UserLocale md.Enum      `json:"userLocale"`
	UserRole   md.Enum      `json:"userRole"`
	UserName   string       `json:"userName"`
	UserPwd    string       `json:"UserPwd"`
	TimeZone   InitTimeZone `json:"timeZone"`
	Database   ProjDb       `json:"database"`
}

func NewProjConfig(id string) *ProjConfig {
	return &ProjConfig{ID: id}
}

func (p *ProjConfig) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ProjConfig) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *ProjConfig) SaveToFile(fileName string) error {
	b, err := p.Marshal()
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, b, FILE_PERMISSION)
}

func (p *ProjConfig) RestoreFromFile(fileName string) error {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return p.Unmarshal(b)
}

type TerminalColors struct {
	CommentColor *color.Color
	PromptColor  *color.Color
	ErrorColor   *color.Color
}

func (p *ProjConfig) InteractiveFillig(projDir string) error {
	var err error
	confFileName := p.ID + PROJ_PARAM_FILE_EXT

	defer func() {
		if err != nil {
			if err := p.SaveToFile(confFileName); err != nil {
				TermColors.ErrorColor.Printf("Error occured while filling new project parameters, error saving configuration to a file: %v", err)
			} else {
				TermColors.ErrorColor.Printf("Error occured while filling new project parameters, configuration saved to file: %s", confFileName)
			}
		}
	}()

	TermColors.CommentColor.Printf("Project directory: %s\n", projDir)
	TermColors.CommentColor.Println("Fill in parameters. If empty string is entered, then default value will be used. On default value absence, empty string will be used if the parameter is not obligatory (not marked with asterisk *).")

	spacialDatabase := "No"

	//loop till correct
	done := false
	for !done {
		//GIT
		p.Git.Name, err = readText("GIT user name", p.Git.Name, false, TermColors)
		if err != nil {
			return err
		}
		p.Git.Descr, err = readText("GIT description", p.Git.Descr, false, TermColors)
		if err != nil {
			return err
		}

		//user role
		if p.UserRole.ID == "" {
			p.UserRole.ID = INIT_USER_ROLE_NAME
		}
		p.UserRole.ID, err = readText("Application user role name *", p.UserRole.ID, true, TermColors)
		if err != nil {
			return err
		}
		if p.UserRole.Descr == "" {
			p.UserRole.Descr = INIT_USER_ROLE_DESCR
		}
		p.UserRole.Descr, err = readText("Application user role description *", p.UserRole.Descr, true, TermColors)
		if err != nil {
			return err
		}

		//user name&&pwd
		if p.UserName == "" {
			p.UserName = p.ID
		}
		p.UserName, err = readText("Application user name *", p.UserName, true, TermColors)
		if err != nil {
			return err
		}

		if p.UserPwd == "" {
			p.UserPwd = INIT_USER_PWD
		}
		p.UserPwd, err = readText("Application user password *", p.UserPwd, true, TermColors)
		if err != nil {
			return err
		}

		//database
		if p.Database.Host == "" {
			p.Database.Host = INIT_DB_HOST
		}
		p.Database.Host, err = readText("Database host (SERVER:PORT) *", p.Database.Host, true, TermColors)
		if err != nil {
			return err
		}
		if p.Database.Schema == "" {
			p.Database.Schema = INIT_DB_SCHEMA
		}
		p.Database.Schema, err = readText("Database schema *", p.Database.Schema, true, TermColors)
		if err != nil {
			return err
		}

		if p.Database.Name == "" {
			p.Database.Name = p.ID
		}
		p.Database.Name, err = readText("Database name *", p.Database.Name, true, TermColors)
		if err != nil {
			return err
		}
		if p.Database.UserName == "" {
			p.Database.UserName = p.ID
		}
		p.Database.UserName, err = readText("Database user name *", p.Database.UserName, true, TermColors)
		if err != nil {
			return err
		}
		p.Database.UserPwd, err = readText("Database user password *", p.Database.UserPwd, true, TermColors)
		if err != nil {
			return err
		}

		//superuser name&&password for create db migration
		p.Database.SuperuserName, err = readText("Database superuser name (to create database) *", INIT_DB_SUPERUSER, true, TermColors)
		if err != nil {
			return err
		}

		//spacial
		if p.Database.Spacial {
			spacialDatabase = "Yes"
		}
		p.Database.Spacial, err = readBool("Is the database spacial? (postgist extention will be installed)", spacialDatabase, true, TermColors)
		if err != nil {
			return err
		}

		//check if all's OK
		projStruct, err := json.MarshalIndent(p, "", "		")
		if err != nil {
			return err
		}
		TermColors.CommentColor.Println(string(projStruct))

		done, err = readBool("If everything is correct, input Yes, or any other value to start over", "", true, TermColors)
		if err != nil {
			return err
		}
	}
	return nil
}

// Init initializes new project. Project config structure must be properly filled.
func (p *ProjConfig) Init(projDir string) error {
	//template parameters: all possible parameters from initialization templates
	tmplProjParams := map[string]interface{}{"APP_NAME": p.ID,
		"DB_USER":           p.Database.UserName,
		"DB_PASSWORD":       p.Database.UserPwd,
		"DB_NAME":           p.Database.Name,
		"DB_SCHEMA":         p.Database.Schema,
		"YEAR":              time.Now().Format("2006"), // used in html templates
		"TIME_ZONE_NAME":    p.TimeZone.Name,
		"TIME_ZONE_DESCR":   p.TimeZone.Descr,
		"TIME_ZONE_OFFSET":  p.TimeZone.Offset,
		"USER":              p.UserName,
		"USER_PWD":          p.UserPwd,
		"USER_ROLE_ID":      p.UserRole.ID,
		"USER_ROLE_DESCR":   p.UserRole.Descr,
		"USER_LOCALE_ID":    p.UserLocale.ID,
		"USER_LOCALE_DESCR": p.UserLocale.Descr,
		"CREATE_DB_MIG_UP":  "", //gitignore, late initialization
		"INIT_DB_MIG_UP":    "", //gitignore, late initialization
	}
	if p.Database.Spacial {
		tmplProjParams["DB_SPACIAL"] = "TRUE"
	}

	//gobizap cli directory with templates
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	templateDir := filepath.Join(path.Dir(exe), PROJ_DIR)

	//clear all files in sql folder to prevent dubling of sql migrations

	if err := copyProjectFiles(templateDir, projDir, tmplProjParams); err != nil {
		return fmt.Errorf("copyProjectFiles() failed: %v", err)
	}

	//split host:port
	dbHost, dbPort, err := md.ParseDbHost(p.Database.Host)
	if err != nil {
		return err
	}
	dbConn := &md.DbConn{User: p.Database.SuperuserName,
		Pwd:    "",
		Server: dbHost,
		Port:   dbPort,
		DbName: "postgres",
	}

	mgr := sqlmigr.NewMigrator(filepath.Join(projDir, md.BUILD_DIR, md.SQL_DIR))

	//create db migration
	crMigTime := time.Now()
	projMigr := md.ProjMigration{Migrator: mgr,
		MigrationTime:   crMigTime,
		MigrationAction: "dbCreate",
		TemplateDir:     templateDir,
		TemplateParams:  tmplProjParams,
		TemplateUp:      TMPL_CREATE_DB_UP,
		TemplateDown:    TMPL_CREATE_DB_UP,
	}
	if err := projMigr.ApplyMigration(); err != nil {
		return err
	}

	//init db migration
	projMigr.MigrationTime = projMigr.MigrationTime.Add(1)
	projMigr.MigrationAction = "dbInit"
	projMigr.TemplateUp = TMPL_INIT_DB_UP
	projMigr.TemplateDown = TMPL_INIT_DB_DOWN
	if err := projMigr.ApplyMigration(); err != nil {
		return err
	}

	//init new go module
	_, err = md.RunCMD("go mod init "+p.ID, false)
	if err != nil {
		return err
	}
	TermColors.PromptColor.Printf("Successfully created golang module: %s\n", p.ID)

	//apply all migrations
	var lastMigrFileName string
	mgrFiles, err := mgr.NewFileList(time.Time{}, sqlmigr.MG_UP)
	if err != nil {
		return err
	}
	for _, mgrForApply := range mgrFiles {
		TermColors.PromptColor.Printf("Applying database migration: %s\n", mgrForApply.Name)

		scriptFile := mgr.GetMigrFullFileName(sqlmigr.MG_UP, mgrForApply.Name)
		if err := dbConn.ApplySQLScript(scriptFile); err != nil {
			return err
		}
		lastMigrFileName = mgrForApply.Name
	}
	if err := mgr.SetLastMigrFileName(lastMigrFileName); err != nil {
		return err
	}

	TermColors.PromptColor.Printf("Successfully initialized new project in directory: %s\n", projDir)

	return nil
}
