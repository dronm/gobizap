package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	userRoleID := "admin"
	userRoleDescr := "Adminustrator"
	userLocaleID := "admin"
	userLocaleDescr := "Administrator"
	timeZoneName := "Asia/Yekaterinburg"
	timeZoneDescr := "Asia/Yekaterinburg descr"
	timeZoneOffset := "+02:00"
	userName := "admin"
	userPwd := "123456"
	dbSuperuser := "postgres"
	dbHost := "localhost:9999"
	dbSchema := "mySchema"

	confCont := `{
		"initDefValues":{
			"enums":{
				"userRole": {"id": "` + userRoleID + `", "descr": "` + userRoleDescr + `"},
				"userLocale": {"id": "` + userLocaleID + `", "descr": "` + userLocaleDescr + `"}
			},
			"timeZone": {"name": "` + timeZoneName + `", "descr": "` + timeZoneDescr + `", "offset": "` + timeZoneOffset + `"},
			"userName":"` + userName + `",
			"userPwd": "` + userPwd + `",
			"db":{
				"superuser":"` + dbSuperuser + `",
				"host":"` + dbHost + `",
				"schema":"` + dbSchema + `"	
			}
		}
	}`

	conf := &Config{}
	if err := conf.ReadData([]byte(confCont)); err != nil {
		t.Fatalf("conf.ReadData() failed: %v", err)
	}

	if conf.InitDefValues.Enums.UserRole.ID != userRoleID {
		t.Fatalf("Enum user role ID wanted: %s, got: %s", userRoleID, conf.InitDefValues.Enums.UserRole.ID)
	}
	if conf.InitDefValues.Enums.UserRole.Descr != userRoleDescr {
		t.Fatalf("userRoleID wanted: %s, got: %s", userRoleDescr, conf.InitDefValues.Enums.UserRole.Descr)
	}
	if conf.InitDefValues.Enums.UserLocale.ID != userLocaleID {
		t.Fatalf("userLocaleID wanted: %s, got: %s", userLocaleID, conf.InitDefValues.Enums.UserLocale.ID)
	}
	if conf.InitDefValues.Enums.UserLocale.Descr != userLocaleDescr {
		t.Fatalf("userLocaleDescr wanted: %s, got: %s", userLocaleDescr, conf.InitDefValues.Enums.UserLocale.Descr)
	}
	if conf.InitDefValues.TimeZone.Name != timeZoneName {
		t.Fatalf("timeZoneName wanted: %s, got: %s", timeZoneName, conf.InitDefValues.TimeZone.Name)
	}
	if conf.InitDefValues.TimeZone.Descr != timeZoneDescr {
		t.Fatalf("timeZoneDescr wanted: %s, got: %s", timeZoneDescr, conf.InitDefValues.TimeZone.Descr)
	}
	if conf.InitDefValues.TimeZone.Offset != timeZoneOffset {
		t.Fatalf("timeZoneOffset wanted: %s, got: %s", timeZoneOffset, conf.InitDefValues.TimeZone.Offset)
	}
	//db
	if conf.InitDefValues.Db.Host != dbHost {
		t.Fatalf("DbHost wanted: %s, got: %s", dbHost, conf.InitDefValues.Db.Host)
	}
	if conf.InitDefValues.Db.Schema != dbSchema {
		t.Fatalf("dbSchema wanted: %s, got: %s", dbSchema, conf.InitDefValues.Db.Schema)
	}
	if conf.InitDefValues.Db.Superuser != dbSuperuser {
		t.Fatalf("dbSuperuser wanted: %s, got: %s", dbSuperuser, conf.InitDefValues.Db.Superuser)
	}

	if conf.InitDefValues.UserName != userName {
		t.Fatalf("userName wanted: %s, got: %s", userName, conf.InitDefValues.UserName)
	}
	if conf.InitDefValues.UserPwd != userPwd {
		t.Fatalf("userPwd wanted: %s, got: %s", userPwd, conf.InitDefValues.UserPwd)
	}
}

func TestCopyProjFiles(t *testing.T) {
	//base directory
	baseDir, err := os.MkdirTemp(os.TempDir(), "gobizap")
	if err != nil {
		t.Fatalf("os.MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(baseDir) // clean up

	var dirPermis os.FileMode = 0750

	sourceDir := filepath.Join(baseDir, "source")
	destDir := filepath.Join(baseDir, "dest")

	//source
	if err := os.Mkdir(sourceDir, dirPermis); err != nil {
		t.Fatalf("os.MkdirAll failed: %v", err)
	}

	//source/dir1
	if err := os.Mkdir(filepath.Join(sourceDir, "dir1"), dirPermis); err != nil {
		t.Fatalf("os.MkdirAll failed: %v", err)
	}

	//source/dir1/sub11
	if err := os.Mkdir(filepath.Join(sourceDir, "dir1", "sub11"), dirPermis); err != nil {
		t.Fatalf("os.MkdirAll failed: %v", err)
	}

	//source/dir2
	if err := os.Mkdir(filepath.Join(sourceDir, "dir2"), dirPermis); err != nil {
		t.Fatalf("os.MkdirAll failed: %v", err)
	}

	//source/dir2/sub21
	if err := os.Mkdir(filepath.Join(sourceDir, "dir2", "sub21"), dirPermis); err != nil {
		t.Fatalf("os.MkdirAll failed: %v", err)
	}
	//file
	if err := os.WriteFile(filepath.Join(sourceDir, "dir2", "sub21", "file1.txt"), []byte("Some test data and some value for changing: {{PAR1}}"), FILE_PERMISSION); err != nil {
		t.Fatalf("os.WriteFile() failed: %v", err)
	}
	//template file
	if err := os.WriteFile(filepath.Join(sourceDir, "dir2", "sub21", "file11.txt.tmpl"), []byte("Some test data and some value for changing: {{PAR1}}"), FILE_PERMISSION); err != nil {
		t.Fatalf("os.WriteFile() failed: %v", err)
	}

	//source/dir2/sub22{{PAR1}}
	if err := os.Mkdir(filepath.Join(sourceDir, "dir2", "sub22_{{PAR1}}"), dirPermis); err != nil {
		t.Fatalf("os.MkdirAll failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "dir2", "sub22_{{PAR1}}", "file2_{{PAR1}}.txt"), []byte("Some test data: {{PAR1}}"), FILE_PERMISSION); err != nil {
		t.Fatalf("os.WriteFile() failed: %v", err)
	}
	//template
	if err := os.WriteFile(filepath.Join(sourceDir, "dir2", "sub22_{{PAR1}}", "file22_{{PAR1}}.txt.tmpl"), []byte("Some test data: {{PAR1}}"), FILE_PERMISSION); err != nil {
		t.Fatalf("os.WriteFile() failed: %v", err)
	}

	params := map[string]interface{}{"PAR1": "par1_val"}
	if err := copyProjectFiles(sourceDir, destDir, params); err != nil {
		t.Fatalf("copyProjectFiles() failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "dir1", "sub11")); os.IsNotExist(err) {
		t.Fatalf("Directory 'dir1/sub11' not found in destination: %s", destDir)
	}

	if _, err := os.Stat(filepath.Join(destDir, "dir2", "sub21")); os.IsNotExist(err) {
		t.Fatalf("Directory 'dir2/sub22' not found in destination: %s", destDir)
	}

	//parameters in dir names
	if _, err := os.Stat(filepath.Join(destDir, "dir2", "sub22_par1_val")); os.IsNotExist(err) {
		t.Fatalf("Directory 'dir2/sub22_par1_val' not found in destination: %s", destDir)
	}

	//check for ordinary file
	if _, err := os.Stat(filepath.Join(destDir, "dir2", "sub21", "file1.txt")); os.IsNotExist(err) {
		t.Fatalf("Directory 'dir2/sub21/file1.txt' not found in destination: %s", destDir)
	}

	//check parameters in file names
	if _, err := os.Stat(filepath.Join(destDir, "dir2", "sub22_par1_val", "file2_par1_val.txt")); os.IsNotExist(err) {
		t.Fatalf("Directory 'dir2/sub22_par1_val/file2_par1_val.txt' not found in destination: %s", destDir)
	}

	//check file from templates
	file2 := filepath.Join(destDir, "dir2", "sub22_par1_val", "file22_par1_val.txt")
	file2_cont, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	if !strings.Contains(string(file2_cont), "par1_val") {
		t.Fatalf("Value of PAR1 not found in file: %s", file2)
	}

	file1 := filepath.Join(destDir, "dir2", "sub21", "file11.txt")
	file1_cont, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	if !strings.Contains(string(file1_cont), "par1_val") {
		t.Fatalf("Value of PAR1 not found in file: %s", file1)
	}
}
