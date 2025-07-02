package md

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

const (
	MD_FILE_NAME = "metadata.xml"

	BUILD_DIR    = "src"
	SQL_DIR      = "sql"
	CONTR_DIR    = "controllers"
	MODEL_DIR    = "models"
	ENUM_DIR     = "enums"
	CONSTANT_DIR = "constants"
	CONTROLS_DIR = "custom_controls"
	CSS_DIR      = "css"
	JS_DIR       = "js"
	UPDATES_DIR  = "updates"

	FILE_PERMISSION = 0775

	PROJ_DIR = "project"
	TMPL_DIR = "templates"

	MOD_FILE        = "go.mod"
	CONFIG_FILE_EXT = ".json"
)

func GetProjectDir() (string, error) {
	return os.Getwd()
}

// GetProjectModuleName returns project name from
// mod file inside project directory
func GetProjectModuleName(projDir string) (string, error) {
	mod_file := filepath.Join(projDir, MOD_FILE)
	mod_file_cont, err := os.ReadFile(mod_file)
	if err != nil {
		return "", fmt.Errorf("os.ReadFile() failed: %v", err)
	}
	if len(mod_file_cont) == 0 {
		return "", fmt.Errorf("%s file is empty", mod_file)
	}

	//mod_file_cont has at least one line
	module := bytes.Split(bytes.ReplaceAll(mod_file_cont, []byte("\r\n"), []byte("\n")), []byte("\n"))[0]
	module_name := bytes.ReplaceAll(module, []byte("module "), []byte(""))
	return string(module_name), nil
}

// GetProjectConfigFile returns the name of the project configuration file.
// Configuration is retrieved from project mod.go's module name
// If configuration file not found error is returned.
func GetProjectConfigFile() (string, error) {
	proj_dir, err := GetProjectDir()
	if err != nil {
		return "", err
	}

	if config_name, err := GetProjectModuleName(proj_dir); err == nil {
		return config_name + CONFIG_FILE_EXT, nil
	}

	return "", fmt.Errorf("configuration file not found, search path: %s", proj_dir)
}
