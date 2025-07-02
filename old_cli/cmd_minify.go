package main

import (
	"fmt"
	"github.com/dchest/jsmin"
	"os"
	"path/filepath"
)

const (
	MINIFIED_JS_NAME  = "lib.js"
	MINIFIED_CSS_NAME = "style.css"
)

type MdScript interface {
	GetFile() string
	GetCompressed() bool
	GetStandalone() bool
}

func RunMinify(args []string) error {
	conf := Config{}
	if err := conf.Read(); err != nil {
		return fmt.Errorf("conf.ReadProjectConfig() failed: %v", err)
	}

	md, err := NewMetadata()
	if err != nil {
		return fmt.Errorf("NewMetadata() failed: %v", err)
	}

	proj_dir, err := GetProjectDir()
	if err != nil {
		return fmt.Errorf("GetProjectDir() failed: %v", err)
	}
	//fmt.Printf("md:%+v\n", md)

	//js
	if err := minifyScriptList(md.JSScripts,
		filepath.Join(proj_dir, JS_DIR),
		MINIFIED_JS_NAME,
	); err != nil {
		return err
	}

	//css
	if err := minifyScriptList(md.CSSScripts,
		filepath.Join(proj_dir, CSS_DIR),
		MINIFIED_CSS_NAME,
	); err != nil {
		return err
	}

	return nil
}

func minifyScriptList[T MdScript](scripts []T, scriptDir string, minName string) error {
	var list_data []byte
	for _, script := range scripts {
		if err := minifyScript(script, scriptDir, list_data); err != nil {
			return err
		}
	}
	if len(list_data) == 0 {
		return nil
	}
	min_file := filepath.Join(scriptDir, minName)
	if err := os.WriteFile(min_file, list_data, FILE_PERMISSION); err != nil {
		return fmt.Errorf("os.WriteFile() failed: %v", err)
	}
	return nil
}

func minifyScript(script MdScript, scriptDir string, scriptData []byte) error {
	if script.GetStandalone() {
		return nil
	}
	f_name := filepath.Join(scriptDir, script.GetFile())
	if !script.GetCompressed() {
		f, err := os.ReadFile(f_name)
		if err != nil {
			//could be missing
			fmt.Printf("os.ReadFile() failed: %v\n", err)
			return nil
		}
		js_data_min, err := jsmin.Minify(f)
		if err != nil {
			return fmt.Errorf("jsmin.Minify() failed: %v", err)
		}
		scriptData = append(scriptData, js_data_min...)

	} else {
		f, err := os.ReadFile(f_name)
		if err != nil {
			fmt.Printf("os.ReadFile() failed: %v\n", err)
			return nil
		}
		//could be missing
		scriptData = append(scriptData, f...)
	}
	return nil
}
