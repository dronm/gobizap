package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dronm/gobizap/md"
)

const (
	CONFIG_FILE_EXT      = ".json"
	MOD_FILE             = "go.mod"
	NEW_FILE_PERMISSIONS = 0775
)

func FileExists(fileName string) (bool, error) {
	if _, err := os.Stat(fileName); err != nil && errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func GetProjectDir() (string, error) {
	return os.Getwd()
}

// GetProjectModuleName returns project name from
// mod file inside project directory
func GetProjectModuleName(projDir string) (string, error) {
	modFile := filepath.Join(projDir, MOD_FILE)
	modFileCont, err := os.ReadFile(modFile)
	if err != nil {
		return "", fmt.Errorf("os.ReadFile() failed: %v", err)
	}
	if len(modFileCont) == 0 {
		return "", fmt.Errorf("%s file is empty", modFile)
	}

	//mod_file_cont has at least one line
	module := bytes.Split(bytes.ReplaceAll(modFileCont, []byte("\r\n"), []byte("\n")), []byte("\n"))[0]
	moduleName := bytes.ReplaceAll(module, []byte("module "), []byte(""))
	return string(moduleName), nil
}

// GetProjectConfigFile returns the name of the project configuration file.
// Configuration is retrieved from project mod.go's module name
// If configuration file not found error is returned.
func GetProjectConfigFile() (string, error) {
	projDir, err := GetProjectDir()
	if err != nil {
		return "", err
	}

	if config_name, err := GetProjectModuleName(projDir); err == nil {
		return config_name + CONFIG_FILE_EXT, nil
	}

	return "", fmt.Errorf("configuration file not found, search path: %s", projDir)
}

// ClearDir deletes all files in the given directory.
func ClearDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// ApplySQLScript is a helper function.
// It runs one sql script with psql and given
// connection parameters.
func ApplySQLScript(dbConn *md.DbConn, scriptFile string) error {
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
