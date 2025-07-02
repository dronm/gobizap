package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
)

const CONFIG_FILE_EXT = ".json"

type InitTimeZone struct {
	Name   string `json:"name"`
	Descr  string `json:"descr"`
	Offset string `json:"offset"`
}

type InitDbDefValues struct {
	Superuser string `json:"superuser"`
	Host      string `json:"host"`
	Schema    string `json:"schema"`
}

type InitDefEnums struct {
	UserRole   Enum `json:"userRole"`
	UserLocale Enum `json:"userLocale"`
}

type InitDefValues struct {
	Enums    InitDefEnums    `json:"enums"`
	TimeZone InitTimeZone    `json:"timeZone"`
	UserName string          `json:"userName"`
	UserPwd  string          `json:"userPwd"`
	Db       InitDbDefValues `json:"db"`
}

// Config is an Application configuration structure.
type Config struct {
	GoFormatCommand string        `json:"goFormatCommand"`
	InitDefValues   InitDefValues `json:"initDefValues"`
}

func NewConfig() (*Config, error) {
	conf := Config{}
	if err := conf.Read(); err != nil {
		return nil, err
	}
	return &conf, nil
}

// Read reads cli configuration.
func (c *Config) Read() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe = strings.TrimRight(exe, ".exe")
	data, err := os.ReadFile(exe + CONFIG_FILE_EXT)
	if err != nil {
		return err
	}
	return c.ReadData(data)
}

func (c *Config) ReadData(data []byte) error {
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	return json.Unmarshal(data, c)
}
