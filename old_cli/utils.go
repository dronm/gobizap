package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/hoisie/mustache"
)

const (
	MOD_FILE        = "go.mod"
	LOG_DATE_FORMAT = "2006-01-02T15:04:05"
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

func ApplyTemplate(source []byte, tmplParams map[string]interface{}) ([]byte, error) {
	cont := bytes.Trim([]byte(mustache.Render(string(source), tmplParams)), " ")
	cont = bytes.Trim(cont, string([]byte{0x0d, 0x0a}))
	cont = bytes.Trim(cont, string([]byte{0x0a}))
	return cont, nil
}

// RenderTemplate
func ApplyTemplateFromFile(sourceFileName string, tmplParams map[string]interface{}) ([]byte, error) {
	source_cont, err := os.ReadFile(sourceFileName)
	if err != nil {
		return []byte{}, err
	}
	return ApplyTemplate(source_cont, tmplParams)
}

// RenderTemplateFromFile renders mustache template with given parameters
func RenderTemplate(sourceFile, destFile string, tmplParams map[string]interface{}) error {
	source_cont, err := os.ReadFile(sourceFile)
	if err != nil {
		return err
	}
	cont, err := ApplyTemplate(source_cont, tmplParams)
	if err != nil {
		return err
	}
	if err := os.WriteFile(destFile, cont, FILE_PERMISSION); err != nil {
		return err
	}
	color.Yellow("Rendered template: %s", destFile)
	return nil
}

func LogWarning(process, s string, params ...interface{}) {
	if params == nil {
		color.Blue(process + " " + time.Now().Format(LOG_DATE_FORMAT) + " " + s)
	} else {
		color.Blue(process+" "+time.Now().Format(LOG_DATE_FORMAT)+" "+s, params...)
	}
}

func LogInfo(process, s string, params ...interface{}) {
	if params == nil {
		color.Green(process + " " + time.Now().Format(LOG_DATE_FORMAT) + " " + s)
	} else {
		color.Green(process+" "+time.Now().Format(LOG_DATE_FORMAT)+" "+s, params...)
	}
}

// dataType is passed as parameter because there are two type fields: dataType && oldDataType
func SQLDataType(model Model, field Field, dataType string) (string, error) {
	type_sql := ""
	if field.PrimaryKey {
		type_sql = "serial"

	} else if (dataType == "String" && field.Length != "") || dataType == "Char" || dataType == "Password" {
		type_sql = fmt.Sprintf("varchar(%s)", field.Length)

	} else if dataType == "Text" || dataType == "String" {
		type_sql = "text"

	} else if dataType == "Int" {
		type_sql = "int"

	} else if dataType == "Enum" {
		type_sql = field.EnumID

	} else if dataType == "DateTime" {
		type_sql = "timestamp without time zone"

	} else if dataType == "DateTimeTZ" {
		type_sql = "timestamp"

	} else if dataType == "Date" {
		type_sql = "date"

	} else if dataType == "Time" {
		type_sql = "time"

	} else if dataType == "Interval" {
		type_sql = "interval"

	} else if dataType == "Geometry" {
		type_sql = "geometry"

	} else if dataType == "Bool" {
		type_sql = "bool"

	} else {
		return "", fmt.Errorf("sql data type not defined for a filed: %s, model: %s", field.ID, model.ID)
	}

	return type_sql, nil
}

// golangDataType returns data type for golang based on
// metadata type
func GolangDataType(model Model, field Field) (string, error) {
	type_go := ""
	if field.DataType == "String" || field.DataType == "Text" || field.DataType == "Char" || field.DataType == "Password" {
		type_go = "fields.ValText"

	} else if field.DataType == "Int" || field.DataType == "BigInt" || field.DataType == "SmallInt" {
		type_go = "fields.ValInt"

	} else if field.DataType == "Bool" || field.DataType == "Boolean" {
		type_go = "fields.ValBool"

	} else if field.DataType == "Date" {
		type_go = "fields.ValDate"

	} else if field.DataType == "DateTime" {
		type_go = "fields.ValDateTime"

	} else if field.DataType == "Time" || field.DataType == "Interval" {
		type_go = "fields.ValTime"

	} else if field.DataType == "DateTimeTZ" {
		type_go = "fields.ValDateTimeTZ"

	} else if field.DataType == "Float" || field.DataType == "Numeric" {
		type_go = "fields.ValFloat"

	} else if field.DataType == "JSON" || field.DataType == "JSONB" {
		type_go = "fields.ValJSON"

	} else if field.DataType == "XML" {
		type_go = "fields.ValXML"

	} else if field.DataType == "Bytea" {
		type_go = "fields.ValBytea"

	} else if field.DataType == "Array" {
		type_go = "fields.ValArray"

	} else if field.DataType == "GeomPolygon" {
		type_go = "fields.ValGeomPolygon"

	} else if field.DataType == "GeomPoint" {
		type_go = "fields.ValGeomPoint"

	} else if field.DataType == "Enum" {
		type_go = "enums.ValEnum_" + field.EnumID

	} else {
		type_go = field.DataType
		return "", fmt.Errorf("golang data type not defined for a field: %s, model: %s", field.ID, model.ID)
	}
	return type_go, nil
}

func FileExists(fileName string) (bool, error) {
	if _, err := os.Stat(fileName); err != nil && errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func DeleteFileIfExists(fileName string) error {
	if exists, err := FileExists(fileName); err != nil {
		return err
	} else if exists {
		if err := os.Remove(fileName); err != nil {
			return err
		}
	}

	return nil
}

// DeleteJavascript deletes script from metadata
func DeleteJavascript(md *Metadata, relFileName string) {
	for i, scr := range md.JSScripts {
		if scr.File != relFileName {
			continue
		}
		if len(md.JSScripts) > i+1 { // not the last sctipt
			md.JSScripts = append(md.JSScripts[:i], md.JSScripts[i+1:]...)
		} else {
			md.JSScripts = md.JSScripts[:i]
		}
		break
	}
}

func FormatGoFile(fileName string) error {
	conf, err := NewConfig()
	if err != nil {
		return err
	}
	if _, err := RunCMD(conf.GoFormatCommand+" "+fileName, false); err != nil {
		return err
	}
	return nil
}
