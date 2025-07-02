package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const CMD_INIT = "init"

const (
	PROJ_PARAM_FILE_EXT = ".prm"
)

type InitHandler struct {
}

func (h InitHandler) OnCommand(args []string) error {
	var projName string
	if len(args) > 0 {
		projName = args[0]
	}

	projDir, err := GetProjectDir()
	if err != nil {
		return fmt.Errorf("GetProjectDir() failed: %v", err)
	}

	if projName == "" {
		//dir name as the project name
		_, projName = filepath.Split(projDir)
	}

	projConf := NewProjConfig(projName)

	projArgFileName := filepath.Join(projDir, projName+PROJ_PARAM_FILE_EXT)
	if projArgFileExists, _ := FileExists(projArgFileName); projArgFileExists {
		if err := projConf.RestoreFromFile(projArgFileName); err != nil {
			return err
		}
		if err := ClearDir(projDir); err != nil {
			return err
		}

	} else {
		//file does not exist, interactive filling
		dirEmpty, err := h.IsProjDirEmpty(projDir, projArgFileName)
		if err != nil {
			return err
		}

		if !dirEmpty {
			res, err := projConf.readBool("Project directory is not empty, continue?", "No", true)
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

	return nil
}

func (c InitHandler) Help() {
	formatter := TermColors.OkColor.SprintfFunc()
	fmt.Printf("Command:\t\t\t%s\n", formatter("init"))
	fmt.Printf("Arguments:\t\t\t%s\n", formatter("<APP_NAME>"))
	fmt.Printf("%s command initializes new project in the working directory. All necessary files, folders, symlinks are created.\n", formatter("init"))
	fmt.Printf("%s is the name for a new project. The argument is not obligatory. If not defind, project folder will be used for a name.\n", formatter("<APP_NAME>"))
	fmt.Println("")
}

func (c InitHandler) IsProjDirEmpty(projDir, confFileName string) (bool, error) {
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
