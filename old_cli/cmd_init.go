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

	"github.com/dronm/sqlmigr"
	"github.com/fatih/color"
	"github.com/hoisie/mustache"
)

const (
	TEMPL_EXT string = ".tmpl" // template extension folder or file

	APP_NAME_PARAM_IND = 1

	FILE_PERMISSION = 0775

	INIT_TMPL_FILE_BACKUP = "init.tmpl"

	//default values if no config
	DEF_APP_USER_ROLE_ID    = "admin"
	DEF_APP_USER_ROLE_DESCR = "Administrator"
)

type AppGit struct {
	Name  string `json:"name"`
	Descr string `json:"descr"`
}

type AppDb struct {
	Host          string `json:"host"`
	Schema        string `json:"schema"`
	Name          string `json:"name"`
	UserName      string `json:"userName"`
	UserPwd       string `json:"UserPwd"`
	Spacial       bool   `json:"spacial"`
	SuperuserName string `json:"superuserName"`
}

// Application configuration structure
type AppConfig struct {
	ID string `json:"id"`
	//AuthorName    string
	//TechEmail     string
	Git        AppGit       `json:"git"`
	UserLocale Enum         `json:"userLocale"`
	UserRole   Enum         `json:"userRole"`
	UserName   string       `json:"userName"`
	UserPwd    string       `json:"UserPwd"`
	TimeZone   InitTimeZone `json:"timeZone"`
	Database   AppDb        `json:"database"`
}

func (app *AppConfig) Backup() ([]byte, error) {
	return json.Marshal(app)
}

func AppConfigFromBackup(backupData []byte) (*AppConfig, error) {
	app := AppConfig{}
	if err := json.Unmarshal(backupData, &app); err != nil {
		return nil, err
	}
	return &app, nil
}

func NewAppConfig(appName string, conf *Config) *AppConfig {
	app := AppConfig{ID: appName,
		UserLocale: conf.InitDefValues.Enums.UserLocale,
		UserRole:   conf.InitDefValues.Enums.UserRole,
		UserName:   conf.InitDefValues.UserName,
		TimeZone:   conf.InitDefValues.TimeZone,
	}
	if app.UserRole.ID == "" {
		app.UserRole.ID = DEF_APP_USER_ROLE_ID
		app.UserRole.Descr = DEF_APP_USER_ROLE_DESCR
	}
	return &app
}

// InitNewApp initializes new application in a working directory.
// Directory
// The function requires parameters:
//   - APP_NAME
func InitNewApp(args []string) error {
	var err error

	//new project directory
	proj_dir, err := GetProjectDir()
	if err != nil {
		return fmt.Errorf("GetProjectDir() failed: %v", err)
	}

	color.Yellow("Project directory: %s", proj_dir)

	//check if directory is empty. It not, warn.
	proj_dir_empty := true
	err = filepath.Walk(proj_dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if proj_dir != path {
			proj_dir_empty = false
		}
		return nil
	})
	if err != nil {
		return err
	}
	if !proj_dir_empty {
		res, err := readBool("Project directory is not empty, continue?", "No", true, color.FgRed)
		if err != nil {
			return err
		}
		if !res {
			return nil
		}
	}

	app_def_name := ""
	if len(args) >= APP_NAME_PARAM_IND+1 {
		app_def_name = args[APP_NAME_PARAM_IND]
	} else {
		_, app_def_name = filepath.Split(proj_dir)
	}

	conf, err := NewConfig()
	if err != nil {
		return err
	}

	var app *AppConfig
	backup_file := filepath.Join(proj_dir, INIT_TMPL_FILE_BACKUP)
	spacial_database := "No"
	prompt_color := color.FgGreen

	//if backup template exists - restore
	if bk_exists, _ := FileExists(backup_file); bk_exists {
		res, err := readBool("Found backup application configuration, restore?", "Yes", true, color.FgRed)
		if err != nil {
			return err
		}
		if res {
			backup_data, err := os.ReadFile(backup_file)
			if err != nil {
				return err
			}
			app, err = AppConfigFromBackup(backup_data)
			if err != nil {
				return err
			}
			goto appCheck
		}
	}

	app = NewAppConfig(app_def_name, conf)
	color.Yellow("Fill in parameters. If empty string is entered, then default value will be used. On default value absence, empty string will be used if the parameter is not obligatory (not marked with asterisk *).")

	// go back to this spot in case of user rejection
appInitialization:
	//GIT
	app.Git.Name, err = readText("GIT user name", app.Git.Name, false, prompt_color)
	if err != nil {
		return err
	}
	app.Git.Descr, err = readText("GIT description", app.Git.Descr, false, prompt_color)
	if err != nil {
		return err
	}

	//user role
	app.UserRole.ID, err = readText("Application user role name *", app.UserRole.ID, true, prompt_color)
	if err != nil {
		return err
	}

	app.UserRole.Descr, err = readText("Application user role description *", app.UserRole.Descr, true, prompt_color)
	if err != nil {
		return err
	}

	//user name&&pwd
	if app.UserName == "" {
		app.UserName = app.ID
	}
	app.UserName, err = readText("Application user name *", app.UserName, true, prompt_color)
	if err != nil {
		return err
	}

	if app.UserPwd == "" {
		app.UserPwd = conf.InitDefValues.UserPwd
	}
	app.UserPwd, err = readText("Application user password *", app.UserPwd, true, prompt_color)
	if err != nil {
		return err
	}

	//database
	if app.Database.Host == "" {
		app.Database.Host = conf.InitDefValues.Db.Host
	}
	app.Database.Host, err = readText("Database host (SERVER:PORT) *", app.Database.Host, true, prompt_color)
	if err != nil {
		return err
	}
	if app.Database.Schema == "" {
		app.Database.Schema = conf.InitDefValues.Db.Schema
	}
	app.Database.Schema, err = readText("Database schema *", app.Database.Schema, true, prompt_color)
	if err != nil {
		return err
	}

	if app.Database.Name == "" {
		app.Database.Name = app.ID
	}
	app.Database.Name, err = readText("Database name *", app.Database.Name, true, prompt_color)
	if err != nil {
		return err
	}
	if app.Database.UserName == "" {
		app.Database.UserName = app.ID
	}
	app.Database.UserName, err = readText("Database user name *", app.Database.UserName, true, prompt_color)
	if err != nil {
		return err
	}
	app.Database.UserPwd, err = readText("Database user password *", app.Database.UserPwd, true, prompt_color)
	if err != nil {
		return err
	}

	//superuser name&&password for create db migration
	app.Database.SuperuserName, err = readText("Database superuser name (to create database) *", conf.InitDefValues.Db.Superuser, true, prompt_color)
	if err != nil {
		return err
	}

	//spacial
	if app.Database.Spacial {
		spacial_database = "Yes"
	}
	app.Database.Spacial, err = readBool("Is the database spacial? (postgist extention will be installed)", spacial_database, true, prompt_color)
	if err != nil {
		return err
	}

	//check if all's OK
appCheck:
	app_struct, err := json.MarshalIndent(app, "", "		")
	if err != nil {
		return err
	}
	color.Yellow("%s", string(app_struct))

	res, err := readBool("If everything is correct, input Yes, or any other value to start over", "", true, prompt_color)
	if err != nil {
		return err
	}
	if !res {
		goto appInitialization
	}

	//backup configuration
	bk_data, err := app.Backup()
	if err != nil {
		color.Red("Error backingup application configuration: %v", err)
	}
	if err := os.WriteFile(backup_file, bk_data, FILE_PERMISSION); err != nil {
		color.Red("Error backingup application configuration: %v", err)
	}

	//template parameters: all possible parameters from initialization templates
	start_proj_params := map[string]interface{}{"APP_NAME": app.ID,
		"DB_USER":           app.Database.UserName,
		"DB_PASSWORD":       app.Database.UserPwd,
		"DB_NAME":           app.Database.Name,
		"DB_SCHEMA":         app.Database.Schema,
		"YEAR":              time.Now().Format("2006"), // used in html templates
		"TIME_ZONE_NAME":    app.TimeZone.Name,
		"TIME_ZONE_DESCR":   app.TimeZone.Descr,
		"TIME_ZONE_OFFSET":  app.TimeZone.Offset,
		"USER":              app.UserName,
		"USER_PWD":          app.UserPwd,
		"USER_ROLE_ID":      app.UserRole.ID,
		"USER_ROLE_DESCR":   app.UserRole.Descr,
		"USER_LOCALE_ID":    app.UserLocale.ID,
		"USER_LOCALE_DESCR": app.UserLocale.Descr,
		"CREATE_DB_MIG_UP":  "", //gitignore, late initialization
		"INIT_DB_MIG_UP":    "", //gitignore, late initialization
	}
	if app.Database.Spacial {
		start_proj_params["DB_SPACIAL"] = "TRUE"
	}

	//gobizap cli directory with templates
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	templates := filepath.Join(path.Dir(exe), PROJ_DIR)

	if err := copyProjectFiles(templates, proj_dir, start_proj_params); err != nil {
		return fmt.Errorf("copyProjectFiles() failed: %v", err)
	}

	//split host:port
	db_host, db_port, err := ParseDbHost(app.Database.Host)
	if err != nil {
		return err
	}

	//db creation - migration script is executed immediately with superuser priviledges
	create_db_up_tmpl, err := os.ReadFile(filepath.Join(templates, BUILD_DIR, TMPL_DIR, TMPL_CREATE_DB_UP))
	if err != nil {
		return err
	}
	create_db_up := mustache.Render(string(create_db_up_tmpl), start_proj_params)

	//temp file with superuser migration script
	create_db_f, err := os.CreateTemp("", "gobiazp")
	if err != nil {
		return err
	}
	defer os.Remove(create_db_f.Name())
	if _, err := create_db_f.Write([]byte(create_db_up)); err != nil {
		return err
	}
	if err := create_db_f.Close(); err != nil {
		return err
	}
	//apply create migration
	if err := ApplySQLScript(&DbConn{User: app.Database.SuperuserName,
		Pwd:    "",
		Server: db_host,
		Port:   db_port,
		DbName: "postgres"},
		create_db_f.Name()); err != nil {

		return err
	}
	color.Blue("- New database is created: %s", app.Database.Name)

	create_mig_t := time.Now()
	mgr := sqlmigr.NewMigrator(filepath.Join(proj_dir, BUILD_DIR, SQL_DIR))

	//create db migration up&&down
	create_db_fname := mgr.GetMigrFileName(create_mig_t, "createDB", sqlmigr.MG_UP)
	create_db_up_mgr := mgr.GetMigrFullFileName(sqlmigr.MG_UP, create_db_fname)
	if err := RenderTemplate(filepath.Join(templates, BUILD_DIR, TMPL_DIR, TMPL_CREATE_DB_UP),
		create_db_up_mgr,
		start_proj_params); err != nil {
		return err
	}
	if err := RenderTemplate(filepath.Join(templates, BUILD_DIR, TMPL_DIR, TMPL_CREATE_DB_DOWN),
		mgr.GetMigrFullFileName(sqlmigr.MG_DOWN, create_db_fname),
		start_proj_params); err != nil {
		return err
	}
	if err := ApplySQLScript(&DbConn{User: app.Database.SuperuserName,
		Pwd:    "",
		Server: db_host,
		Port:   db_port,
		DbName: "postgres"},
		create_db_up_mgr); err != nil {

		return err
	}
	color.Blue("- New database is created: %s", app.Database.Name)

	//db initialization
	init_db_fname := mgr.GetMigrFileName(create_mig_t, "initDB", sqlmigr.MG_UP)
	init_db_up_mgr := mgr.GetMigrFullFileName(sqlmigr.MG_UP, init_db_fname)
	if err := RenderTemplate(filepath.Join(templates, BUILD_DIR, TMPL_DIR, TMPL_INIT_DB_UP),
		init_db_up_mgr,
		start_proj_params); err != nil {
		return err
	}
	if err := RenderTemplate(filepath.Join(templates, BUILD_DIR, TMPL_DIR, TMPL_INIT_DB_DOWN),
		mgr.GetMigrFullFileName(sqlmigr.MG_DOWN, init_db_fname),
		start_proj_params); err != nil {
		return err
	}
	if err := ApplySQLScript(&DbConn{User: app.Database.UserName,
		Pwd:    app.Database.UserPwd,
		Server: db_host,
		Port:   db_port,
		DbName: app.Database.Name,
	},
		init_db_up_mgr); err != nil {
		return err
	}
	color.Blue("- New database is initialized from migration script: %s", init_db_fname)

	//mark last migration

	if err := mgr.SetLastMigrFileName(init_db_fname); err != nil {
		return err
	}

	//init new go module
	os.Chdir(filepath.Join(proj_dir, BUILD_DIR))
	if _, err := RunCMD("go mod init "+app.ID, false); err != nil {
		return err
	}
	color.Blue("- Golang module is initialized in: %s", proj_dir)

	fmt.Printf("\n\n")
	color.Green("Successfully initialized new project in directory: %s", proj_dir)
	fmt.Printf("\n")

	//remove backup if any
	if err := os.Remove(backup_file); err != nil {
		color.Red("Error deleting backup file: %v", err)
	}

	return nil
}

// copyProjectFiles recursively creates all project files (folders, files, symbolic links)
// from templates in sourceDir.
func copyProjectFiles(sourceDir, destDir string, params map[string]interface{}) error {
	return filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		rel_path := strings.TrimPrefix(path, sourceDir)
		rel_path = strings.ReplaceAll(rel_path, TEMPL_EXT, "")

		name := d.Name()
		is_template := (strings.HasSuffix(name, TEMPL_EXT) && rel_path != filepath.Join(BUILD_DIR, TMPL_DIR))
		if strings.Contains(name, "{{") && strings.Contains(name, "}}") {
			//may also have template in name
			rel_path = mustache.Render(rel_path, params)
		}
		new_file := filepath.Join(destDir, rel_path)
		//find out if it is a symbolic link folder and copy its contents if so
		if d.IsDir() {
			//directory
			if _, err := os.Stat(new_file); os.IsNotExist(err) {
				if err := os.Mkdir(new_file, FILE_PERMISSION); err != nil {
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
				err := os.WriteFile(new_file, new_file_data, FILE_PERMISSION)
				if err != nil {
					return fmt.Errorf("os.WriteFile() failed: %v", err)
				}
			}
			//fmt.Println("Writing to file:", new_file)
		}
		return nil
	})
}

// helper functions
func readText(prompt, defValue string, obligatory bool, promptColor color.Attribute) (string, error) {
	reader := bufio.NewReader(os.Stdin)
scanInput:
	if defValue != "" {
		color.New(promptColor).Printf("%s", prompt)
		color.New(color.FgYellow).Printf(" (default: %s)-> ", defValue)
	} else {
		color.New(promptColor).Printf("%s-> ", prompt)
	}
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.Replace(text, "\n", "", -1)
	if text == "" && defValue != "" {
		text = defValue
	}
	if obligatory && text == "" {
		color.Red("This parameter can not be empty!")
		goto scanInput
	}
	return text, nil
}

func readBool(prompt, defValue string, obligatory bool, promptColor color.Attribute) (bool, error) {
	t, err := readText(prompt, defValue, obligatory, promptColor)
	if err != nil {
		return false, err
	}
	if strings.ToUpper(t) == "Y" || strings.ToUpper(t) == "YES" {
		return true, nil
	}
	return false, nil
}
