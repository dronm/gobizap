package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dronm/gobizap/md"
	"github.com/dronm/sqlmigr"
	"github.com/hoisie/mustache"
)

const (
	INIT_DB_SUPERUSER    = "postgres"
	INIT_DB_SCHEMA       = "public"
	INIT_DB_HOST         = "localhost:5432"
	INIT_USER_PWD        = "123456"
	INIT_USER_ROLE_NAME  = "admin"
	INIT_USER_ROLE_DESCR = "Administrator"

	TMPL_CREATE_DB_UP   = "db_create_up.sql.tmpl"
	TMPL_CREATE_DB_DOWN = "db_create_down.sql.tmpl"
	TMPL_INIT_DB_UP     = "db_init_up.sql.tmpl"
	TMPL_INIT_DB_DOWN   = "db_init_down.sql.tmpl"

	TEMPL_EXT = ".tmpl"
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
	Name string `json:"name"`
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

func NewProjConfig(name string) *ProjConfig {
	return &ProjConfig{Name: name}
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

	return os.WriteFile(fileName, b, NEW_FILE_PERMISSIONS)
}

func (p *ProjConfig) RestoreFromFile(fileName string) error {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return p.Unmarshal(b)
}

func (p *ProjConfig) InteractiveFillig(projDir string) error {
	var err error
	confFileName := p.Name + PROJ_PARAM_FILE_EXT

	defer func() {
		if err != nil {
			if err := p.SaveToFile(confFileName); err != nil {
				TermColors.DangerColor.Printf("Error occured while filling new project parameters, error saving configuration to a file: %v", err)
			} else {
				TermColors.DangerColor.Printf("Error occured while filling new project parameters, configuration saved to file: %s", confFileName)
			}
		}
	}()

	TermColors.InfoColor.Printf("Project directory: %s\n", projDir)
	TermColors.InfoColor.Println("Fill in parameters. If empty string is entered, then default value will be used. On default value absence, empty string will be used if the parameter is not obligatory (not marked with asterisk *).")

	spacialDatabase := "No"

	//loop till correct
	done := false
	for !done {
		//GIT
		p.Git.Name, err = p.readText("GIT user name", p.Git.Name, false)
		if err != nil {
			return err
		}
		p.Git.Descr, err = p.readText("GIT description", p.Git.Descr, false)
		if err != nil {
			return err
		}

		//user role
		if p.UserRole.ID == "" {
			p.UserRole.ID = INIT_USER_ROLE_NAME
		}
		p.UserRole.ID, err = p.readText("Application user role name *", p.UserRole.ID, true)
		if err != nil {
			return err
		}
		if p.UserRole.Descr == "" {
			p.UserRole.Descr = INIT_USER_ROLE_DESCR
		}
		p.UserRole.Descr, err = p.readText("Application user role description *", p.UserRole.Descr, true)
		if err != nil {
			return err
		}

		//user name&&pwd
		if p.UserName == "" {
			p.UserName = p.Name
		}
		p.UserName, err = p.readText("Application user name *", p.UserName, true)
		if err != nil {
			return err
		}

		if p.UserPwd == "" {
			p.UserPwd = INIT_USER_PWD
		}
		p.UserPwd, err = p.readText("Application user password *", p.UserPwd, true)
		if err != nil {
			return err
		}

		//database
		if p.Database.Host == "" {
			p.Database.Host = INIT_DB_HOST
		}
		p.Database.Host, err = p.readText("Database host (SERVER:PORT) *", p.Database.Host, true)
		if err != nil {
			return err
		}
		if p.Database.Schema == "" {
			p.Database.Schema = INIT_DB_SCHEMA
		}
		p.Database.Schema, err = p.readText("Database schema *", p.Database.Schema, true)
		if err != nil {
			return err
		}

		if p.Database.Name == "" {
			p.Database.Name = p.Name
		}
		p.Database.Name, err = p.readText("Database name *", p.Database.Name, true)
		if err != nil {
			return err
		}
		if p.Database.UserName == "" {
			p.Database.UserName = p.Name
		}
		p.Database.UserName, err = p.readText("Database user name *", p.Database.UserName, true)
		if err != nil {
			return err
		}
		p.Database.UserPwd, err = p.readText("Database user password *", p.Database.UserPwd, true)
		if err != nil {
			return err
		}

		//superuser name&&password for create db migration
		p.Database.SuperuserName, err = p.readText("Database superuser name (to create database) *", INIT_DB_SUPERUSER, true)
		if err != nil {
			return err
		}

		//spacial
		if p.Database.Spacial {
			spacialDatabase = "Yes"
		}
		p.Database.Spacial, err = p.readBool("Is the database spacial? (postgist extention will be installed)", spacialDatabase, true)
		if err != nil {
			return err
		}

		//check if all's OK
		projStruct, err := json.MarshalIndent(p, "", "		")
		if err != nil {
			return err
		}
		TermColors.InfoColor.Println(string(projStruct))

		done, err = p.readBool("If everything is correct, input Yes, or any other value to start over", "", true)
		if err != nil {
			return err
		}
	}
	return nil
}

// Init initializes new project. Project config structure must be properly filled.
func (p *ProjConfig) Init(projDir string) error {
	//template parameters: all possible parameters from initialization templates
	tmplProjParams := map[string]interface{}{"APP_NAME": p.Name,
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
	templateDir := filepath.Join(path.Dir(exe), md.PROJ_DIR)
	//what if template directory does not exist
	if tmplEx, _ := FileExists(templateDir); !tmplEx {
		return fmt.Errorf("template directory %s does not exist", templateDir)
	}

	//clear all files in sql folder to prevent dubling of sql migrations

	if err := p.copyProjectFiles(templateDir, projDir, tmplProjParams); err != nil {
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
	_, err = md.RunCMD("go mod init "+p.Name, false)
	if err != nil {
		return err
	}
	TermColors.OkColor.Printf("Successfully created golang module: %s\n", p.Name)

	//apply all migrations
	var lastMigrFileName string
	mgrFiles, err := mgr.NewFileList(time.Time{}, sqlmigr.MG_UP)
	if err != nil {
		return err
	}
	for _, mgrForApply := range mgrFiles {
		TermColors.OkColor.Printf("Applying database migration: %s\n", mgrForApply.Name)

		scriptFile := mgr.GetMigrFullFileName(sqlmigr.MG_UP, mgrForApply.Name)
		if err := dbConn.ApplySQLScript(scriptFile); err != nil {
			return err
		}
		lastMigrFileName = mgrForApply.Name
	}
	if err := mgr.SetLastMigrFileName(lastMigrFileName); err != nil {
		return err
	}

	TermColors.OkColor.Printf("Successfully initialized new project in directory: %s\n", projDir)

	return nil
}

// helper functions
func (p *ProjConfig) readText(prompt, defValue string, obligatory bool) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var text string
	var err error
	for obligatory || text == "" {
		if defValue != "" {
			TermColors.OkColor.Printf("%s", prompt)
			TermColors.InfoColor.Printf(" (default: %s)-> ", defValue)
		} else {
			TermColors.OkColor.Printf("%s-> ", prompt)
		}
		text, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		text = strings.Replace(text, "\n", "", -1)
		if text == "" && defValue != "" {
			text = defValue
		}
		if obligatory && text == "" {
			TermColors.DangerColor.Println("This parameter can not be empty!")
		} else {
			//any answer will do
			break
		}
	}
	return text, nil
}

func (p *ProjConfig) readBool(prompt, defValue string, obligatory bool) (bool, error) {
	t, err := p.readText(prompt, defValue, obligatory)
	if err != nil {
		return false, err
	}
	if strings.ToUpper(t) == "Y" || strings.ToUpper(t) == "YES" {
		return true, nil
	}
	return false, nil
}

// copyProjectFiles recursively creates all project files (folders, files, symbolic links)
// from templates in sourceDir.
func (p *ProjConfig) copyProjectFiles(sourceDir, destDir string, params map[string]interface{}) error {
	return filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		rel_path := strings.TrimPrefix(path, sourceDir)
		rel_path = strings.ReplaceAll(rel_path, TEMPL_EXT, "")

		name := d.Name()
		is_template := (strings.HasSuffix(name, TEMPL_EXT) && rel_path != filepath.Join(md.BUILD_DIR, md.TMPL_DIR))
		if strings.Contains(name, "{{") && strings.Contains(name, "}}") {
			//may also have template in name
			rel_path = mustache.Render(rel_path, params)
		}
		new_file := filepath.Join(destDir, rel_path)
		//find out if it is a symbolic link folder and copy its contents if so
		if d.IsDir() {
			//directory
			if _, err := os.Stat(new_file); os.IsNotExist(err) {
				if err := os.Mkdir(new_file, NEW_FILE_PERMISSIONS); err != nil {
					return fmt.Errorf("os.Mkdir() failed: %v", err)
				}
			}

		} else if d.Type()&os.ModeSymlink != 0 {
			//symlink
			exists, err := FileExists(new_file)
			if err != nil {
				return fmt.Errorf("FileExists() for symlink %s failed: %v", new_file, err)
			}
			if !exists {
				target, err := os.Readlink(path)
				if err != nil {
					return fmt.Errorf("os.Readlink() failed: %v", err)
				}
				if err = os.Symlink(target, new_file); err != nil {
					return fmt.Errorf("os.Symlink() failed: %v", err)
				}
			}
		} else {
			//ordinary file
			exists, err := FileExists(new_file)
			if err != nil {
				return fmt.Errorf("FileExists() for normal file %s failed: %v", new_file, err)
			}
			if !exists {
				var new_file_data []byte
				if is_template {
					transformed_data := mustache.RenderFile(path, params)
					new_file_data = []byte(transformed_data)
				} else {
					var err error
					new_file_data, err = os.ReadFile(path)
					if err != nil {
						return fmt.Errorf("os.ReadFile() failed: %v", err)
					}
				}
				err := os.WriteFile(new_file, new_file_data, NEW_FILE_PERMISSIONS)
				if err != nil {
					return fmt.Errorf("os.WriteFile() failed: %v", err)
				}
			}
			//fmt.Println("Writing to file:", new_file)
		}
		return nil
	})
}
