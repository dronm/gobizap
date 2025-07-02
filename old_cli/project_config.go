package main

import (
	"bytes"
	"encoding/json"
	"os"
)

type ProductionHost struct {
	User     string `json:"user"`
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	AppDir   string `json:"appDir"`
	AppStop  string `json:"appStop"`
	AppStart string `json:"appStart"`
}

type Production struct {
	Compile string           `json:"compile"`
	Files   []string         `json:"files"`
	Hosts   []ProductionHost `json:"hosts"`
}

//	type ProjectDb struct {
//		Primary   string
//		Secondary map[string]string
//	}
//
// ProjectConfig
type ProjectConfig struct {
	Production Production `json:"production"`
	// Db         ProjectDb
}

// Read reads project configuration.
func (c *ProjectConfig) Read() error {
	config_name, err := GetProjectConfigFile()
	if err != nil {
		return err
	}

	file, err := os.ReadFile(config_name + CONFIG_FILE_EXT)
	if err != nil {
		return err
	}
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	return json.Unmarshal([]byte(file), c)
}
