package md

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/dronm/sqlmigr"
	"github.com/hoisie/mustache"
)

type ProjMigration struct {
	Migrator        *sqlmigr.Migrator
	MigrationTime   time.Time
	MigrationAction string
	TemplateDir     string
	TemplateParams  map[string]interface{}
	TemplateUp      string
	TemplateDown    string
}

func (p ProjMigration) ApplyMigration() error {
	tmplFileUp := filepath.Join(p.TemplateDir, BUILD_DIR, TMPL_DIR, p.TemplateUp)
	dataUp, err := ApplyTemplateFromFile(tmplFileUp, p.TemplateParams)
	if err != nil {
		return err
	}
	if err := p.Migrator.Add(p.MigrationTime, p.MigrationAction, sqlmigr.MG_UP, dataUp); err != nil {
		return err
	}
	tmplFileDown := filepath.Join(p.TemplateDir, BUILD_DIR, TMPL_DIR, p.TemplateDown)
	dataDown, err := ApplyTemplateFromFile(tmplFileDown, p.TemplateParams)
	if err != nil {
		return err
	}
	if err := p.Migrator.Add(p.MigrationTime, p.MigrationAction, sqlmigr.MG_DOWN, dataDown); err != nil {
		return err
	}

	return nil
}

// ApplyTemplate apply tmplParams to source data, returns the result byte array.
func ApplyTemplate(source []byte, tmplParams map[string]interface{}) ([]byte, error) {
	cont := bytes.Trim([]byte(mustache.Render(string(source), tmplParams)), " ")
	cont = bytes.Trim(cont, string([]byte{0x0d, 0x0a}))
	cont = bytes.Trim(cont, string([]byte{0x0a}))
	return cont, nil
}

// ApplyTemplateFromFile reads data from sourceFileName and applies tmplParams, returns the result byte array.
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

	return nil
}
