package main

// This application initializes new project in the working directory. It accepts one parameter -
// path to a file containing new project parameters. If there is not such file in the project
// directory, a new file is created, parameters are questioned interactively.
// If there is no start arguments then project directory name is used as project name and
// project paramere file name.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dronm/gobizap/md"
)

const (
	PROJ_PARAM_FILE_EXT = ".prm"
	FILE_PERMISSION     = 0775
)

func main() {
	var projName string

	if len(os.Args) > 1 {
		projName = os.Args[1]
	}
	if err := InitProject(projName); err != nil {
		TermColors.ErrorColor.Println(err)
		os.Exit(1)
	}
}

// InitProject initializes new project in the working directory.
// If projArgFile is an empty string then new file will be created.
func InitProject(projName string) error {
	//new project directory
	projDir, err := md.GetProjectDir()
	if err != nil {
		return fmt.Errorf("md.GetProjectDir() failed: %v", err)
	}

	if projName == "" {
		//dir name as the project name
		_, projName = filepath.Split(projDir)
	}

	projArgFileName := filepath.Join(projDir, projName+PROJ_PARAM_FILE_EXT)

	projConf := NewProjConfig(projName)

	if projArgFileExists, _ := FileExists(projArgFileName); projArgFileExists {
		if err := projConf.RestoreFromFile(projArgFileName); err != nil {
			return err
		}
		if err := ClearDir(projDir); err != nil {
			return err
		}

	} else {
		//file does not exist, interactive filling
		dirEmpty, err := IsProjDirEmpty(projDir, projArgFileName)
		if err != nil {
			return err
		}

		if !dirEmpty {
			res, err := readBool("Project directory is not empty, continue?", "No", true, TermColors)
			if err != nil {
				return err
			}
			if !res {
				//clear exit
				return nil
			}
			if err := ClearDir(projDir); err != nil {
				return err
			}
		}
		if err := projConf.InteractiveFillig(projDir); err != nil {
			return err
		}
		if errSaving := projConf.SaveToFile(projName + PROJ_PARAM_FILE_EXT); errSaving != nil {
			fmt.Printf("error saving configuration to a file: %v\n", err)
		}
	}

	if err := projConf.Init(projDir); err != nil {
		return err
	}

	// color.Green("Successfully initialized new project in directory: %s", proj_dir)
	//
	return nil
}

func showUsage() {

}
func IsProjDirEmpty(projDir, confFileName string) (bool, error) {
	projDirEmpty := true
	if err := filepath.Walk(projDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if projDir != path && info.Name() != confFileName {
			projDirEmpty = false
		}
		return nil

	}); err != nil {
		return false, err
	}
	return projDirEmpty, nil
}

func ClearDir(projDir string) error {
	//delete all files
	d, err := os.Open(projDir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(projDir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
