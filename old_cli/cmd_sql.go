package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/dronm/gobizapp/config"
)

const (
	CONN_INVALID_FORMAT = "invalid connection string format"
)

type DbConn struct {
	User   string
	Pwd    string
	Server string
	Port   int
	DbName string
}

// ParseDbHost parses host in SERVER:PORT format,
// splits SERVER and PORT to different variables.
func ParseDbHost(host string) (string, int, error) {
	srv_port := strings.Split(host, ":")
	if len(srv_port) < 2 {
		return "", 0, fmt.Errorf(CONN_INVALID_FORMAT)
	}
	port, err := strconv.Atoi(srv_port[1])
	if err != nil {
		return "", 0, fmt.Errorf(CONN_INVALID_FORMAT)
	}
	return srv_port[0], port, nil
}

// ParseDbConn is a helper function to parse db connection string
// from template postgresql://USER:PWD@SERVER:PORT/DB_NAME
// It makes a structure where all fields are separate.
func ParseDbConn(conn string) (*DbConn, error) {
	conn_parts := strings.Split(conn, "://")
	if len(conn_parts) < 2 || conn_parts[0] != "postgresql" {
		return nil, fmt.Errorf(CONN_INVALID_FORMAT)
	}

	conn_parts1_parts := strings.Split(conn_parts[1], "@")
	if len(conn_parts1_parts) < 2 {
		return nil, fmt.Errorf(CONN_INVALID_FORMAT)
	}
	user_pwd := strings.Split(conn_parts1_parts[0], ":")
	if len(user_pwd) < 2 {
		return nil, fmt.Errorf(CONN_INVALID_FORMAT)
	}

	host_db := strings.Split(conn_parts1_parts[1], "/")
	if len(host_db) < 2 {
		return nil, fmt.Errorf(CONN_INVALID_FORMAT)
	}

	host, port, err := parseDbHost(host_db[0])
	if err != nil {
		return nil, err
	}

	return &DbConn{User: user_pwd[0],
		Pwd:    user_pwd[1],
		DbName: host_db[1],
		Server: host,
		Port:   port,
	}, nil
}

// RunSQL executes migrations: up, down, add. If no argument supplied then up
// migration is expected
func RunSQL(args []string) error {
	if len(args) <= 1 || args[1] == "up" {
		//up migration
		if len(args) > 2 && args[2] == "all" {
			return RunAllUp()
		}
		return RunUp()

	} else if len(args) > 1 && args[1] == "down" {
		if len(args) > 2 && args[2] == "all" {
			return RunAllDown()
		}
		return RunDown()

	} else if len(args) > 1 && args[1] == "add" {
		//check for action
		if len(args) < 3 {
			//name
			return errors.New("action should be specified")
		}
		act := args[2]
		file_pref, err := RunAdd(act)
		if err != nil {
			return err
		}
		formatter := color.New(color.FgGreen, color.Bold).SprintFunc()
		fmt.Printf("Generated new migration file: %s\n", formatter(file_pref))

	} else if len(args) > 1 && args[1] == "pos" {
		return RunPos()

	} else {
		return errors.New("not supported argument")
	}
	return nil
}

func RunAllUp() error {
	return nil
}

func RunUp() error {
	return nil
}

func RunDown() error {
	return nil
}

func RunAllDown() error {
	return nil
}

func RunPos() error {
	return nil
}

func RunAdd(act string) error {
	return nil
}

// applySQL is a helper function.
// It runs all sql scripts from updates directory on master.
// Scripts are sorted on filemtime
func applySQL(projDir, sqlDir string) error {
	db_conn, err := GetProjDbConn(projDir)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(sqlDir)
	if err != nil {
		return fmt.Errorf("os.ReadDir() failed: %v", err)
	}
	// Create a slice to store file names
	var file_names []string
	for _, file := range files {
		if !file.IsDir() {
			file_names = append(file_names, file.Name())
		}
	}

	// Sort file names by modified time
	sort.Slice(file_names, func(i, j int) bool {
		file1Path := filepath.Join(sqlDir, file_names[i])
		file2Path := filepath.Join(sqlDir, file_names[j])

		file1Info, err := os.Stat(file1Path)
		if err != nil {
			return false
		}

		file2Info, err := os.Stat(file2Path)
		if err != nil {
			return false
		}

		return file1Info.ModTime().Before(file2Info.ModTime())
	})

	for _, file := range file_names {
		if err := ApplySQLScript(db_conn, file); err != nil {
			return err
		}
	}

	return nil
}

// GetProjDbConn fetches DB connection parameters from project module file go.mod
func GetProjDbConn(projDir string) (*DbConn, error) {
	module_name, err := GetProjectModuleName(projDir)
	if err != nil {
		// return nil, err
		//if no module file, then dir name is the conf name
		dirs := strings.Split(projDir, "/")
		if len(dirs) > 0 {
			module_name = dirs[len(dirs)-1]
		} else {
			module_name = projDir
		}
	}

	//file module_name.json must exist
	config_file := module_name + CONFIG_FILE_EXT
	if _, err := os.Stat(config_file); os.IsNotExist(err) {
		return nil, fmt.Errorf("application config file %s not found", config_file)
	}

	conf := config.AppConfig{}
	if err := config.ReadConf(config_file, &conf); err != nil {
		return nil, fmt.Errorf("config.ReadConf() failed: %v", err)
	}
	db_conn, err := parseDbConn(conf.Db.Primary)
	if err != nil {
		return nil, fmt.Errorf("parseDbConn() failed: %v", err)
	}

	return db_conn, nil
}

// ApplySQLScript is a helper function.
// It runs one sql script with psql and given
// connection parameters.
func ApplySQLScript(dbConn *DbConn, scriptFile string) error {
	if dbConn.Pwd != "" {
		if err := os.Setenv("PGPASSWORD", dbConn.Pwd); err != nil {
			return fmt.Errorf("os.Setenv() failed: %v", err)
		}
	}
	bash_cmd := fmt.Sprintf("psql -h %s -p %d -d %s -U %s -f %s",
		dbConn.Server,
		dbConn.Port,
		dbConn.DbName,
		dbConn.User,
		scriptFile,
	)
	out_text, err := RunCMD(bash_cmd, false)
	if err != nil {
		return err
	}
	if len(out_text) > 0 {
		fmt.Println(out_text)
	}

	return nil
}
