package main

import "fmt"

const (
	MD_CMD_ADD = "add"
	MD_CMD_DEL = "del"
	MD_CMD_ALT = "alt"
)

// ModelFieldParam is a template parameter for a field
type ModelFieldTemplateParam struct {
	NAME          string
	ID            string
	TYPE          string
	NOT_FIRST     bool
	NOT_REF_TABLE bool
	SYS_COL       bool
	NAME_ALIAS    string
}

type ModelTemplateParam struct {
	OBJECT_DATA_TABLE   string
	DATA_SCHEMA         string
	FIELDS              map[string]string
	VIRTUAL             bool
	NOT_VIRTUAL         bool
	KEYS                map[string]string
	ENUMS_EXIST         bool
	APP_PACKAGE         string
	BASE_DATA_TABLE     string
	AGG_FUNCTIONS_EXIST bool
}

type ModelRefFieldTemplateParam struct {
	NOT_FIRST bool
	ALIAS     string
	EXPR      string
}

// Build builds project.
func Build(args []string) error {
	md, err := NewMetadata()
	if err != nil {
		return fmt.Errorf("NewMetadata() failed: %v", err)
	}

	proj_dir, err := GetProjectDir()
	if err != nil {
		return fmt.Errorf("GetProjectDir() failed: %v", err)
	}

	if err := BuildControllers(md, proj_dir); err != nil {
		return err
	}

	if err := BuildModels(md, proj_dir); err != nil {
		return err
	}

	if err := BuildEnums(md, proj_dir); err != nil {
		return err
	}

	if err := BuildConstants(md, proj_dir); err != nil {
		return err
	}

	return nil
}
