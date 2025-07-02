package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/dronm/gobizapp/config"
)

const (
	LAST_SQL_SCRIPT  = "last_update.sql"
	ACCUM_SQL_SCRIPT = "update.sql"

	CONN_INVALID_FORMAT = "invalid connection string format"

	NO_NAME_ACTION = "<no_name>"

	MIG_FILE_DATE_FORMAT = "20060102150405"
)

type MigType int

func (m MigType) Descr() string {
	if m == MG_UP {
		return "up"
	} else {
		return "down"
	}
}

const (
	MG_UP MigType = iota
	MG_DOWN
)

type DbConn struct {
	User   string
	Pwd    string
	Server string
	Port   int
	DbName string
}

type MigrationFile struct {
	FileName string
	FilePos  time.Time
	Action   string
}

type MigrationFileList []MigrationFile

func (f MigrationFileList) Len() int {
	return len(f)
}

func (f MigrationFileList) Less(i, j int) bool {
	d1, _, err := migrationFileParams(f[i].FileName)
	if err != nil {
		return false
	}
	d2, _, err := migrationFileParams(f[j].FileName)
	if err != nil {
		return false
	}
	return d1.Before(d2)
}

func (f MigrationFileList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func GetMigrationAbsDir(relDir string) string {
	return filepath.Join(BUILD_DIR, SQL_DIR, relDir)
}

// parseDbHost parses host in SERVER:PORT format,
// splits SERVER and PORT to different variables.
func parseDbHost(host string) (string, int, error) {
	srv_port := strings.Split(host, ":")
	if len(srv_port) < 2 {
		return "", 0, fmt.Errorf(CONN_INVALID_FORMAT)
	}
	port, err := strconv.Atoi(srv_port[1])
	if err != nil {
		return "", 0, fmt.Errorf(CONN_INVALID_FORMAT)
	}
	return srv_port[0], port, nil
}

// parseDbConn is a helper function to parse db connection string
// from template postgresql://USER:PWD@SERVER:PORT/DB_NAME
// It makes a structure where all fields are separate.
func parseDbConn(conn string) (*DbConn, error) {
	conn_parts := strings.Split(conn, "://")
	if len(conn_parts) < 2 || conn_parts[0] != "postgresql" {
		return nil, fmt.Errorf(CONN_INVALID_FORMAT)
	}

	conn_parts1_parts := strings.Split(conn_parts[1], "@")
	if len(conn_parts1_parts) < 2 {
		return nil, fmt.Errorf(CONN_INVALID_FORMAT)
	}
	user_pwd := strings.Split(conn_parts1_parts[0], ":")
	if len(user_pwd) < 2 {
		return nil, fmt.Errorf(CONN_INVALID_FORMAT)
	}

	host_db := strings.Split(conn_parts1_parts[1], "/")
	if len(host_db) < 2 {
		return nil, fmt.Errorf(CONN_INVALID_FORMAT)
	}

	host, port, err := parseDbHost(host_db[0])
	if err != nil {
		return nil, err
	}

	return &DbConn{User: user_pwd[0],
		Pwd:    user_pwd[1],
		DbName: host_db[1],
		Server: host,
		Port:   port,
	}, nil
}

// RunSQL executes migrations: up, down, add. If no argument supplied then up
// migration is expected
func RunSQL(args []string) error {
	if len(args) <= 1 || args[1] == "up" {
		//up migration
		only_first := true
		if len(args) > 2 && args[2] == "all" {
			only_first = false
		}
		return runUpMigration(only_first)

	} else if len(args) > 1 && args[1] == "down" {
		only_first := true
		if len(args) > 2 && args[2] == "all" {
			only_first = false
		}
		return runDownMigration(only_first)

	} else if len(args) > 1 && args[1] == "add" {
		//check for action
		if len(args) < 3 {
			//name
			return errors.New("action should be specified")
		}
		act := args[2]
		file_pref, err := runAddMigration(act)
		if err != nil {
			return err
		}
		formatter := color.New(color.FgGreen, color.Bold).SprintFunc()
		fmt.Printf("Generated new migration file: %s\n", formatter(file_pref))

	} else if len(args) > 1 && args[1] == "pos" {
		return runShowMigration()

	} else {
		return errors.New("not supported argument")
	}
	return nil
}

// runMigrations runs up or down migration
func runMigrations(migDir string, mgType MigType, onlyFirst bool) error {
	mig_pos_file, err := lastMigrationPos()
	if err != nil {
		return err
	}
	var cur_mig_pos time.Time
	if mig_pos_file != "" {
		cur_mig_pos, _, err = migrationFileParams(mig_pos_file)
		if err != nil {
			return err
		}
	}

	mig_abs_dir := GetMigrationAbsDir(migDir)
	mig_files, err := makeMigrationFileList(cur_mig_pos, mig_abs_dir, mgType)
	if err != nil {
		return err
	}

	formatter := color.New(color.FgGreen, color.Bold).SprintFunc()
	if len(mig_files) == 0 {
		return errors.New("no migration found")
	}
	fmt.Printf("found migrations: %s\n", formatter(len(mig_files)))

	proj_dir, err := GetProjectDir()
	if err != nil {
		return err
	}
	//current database connection
	conn, err := GetProjDbConn(proj_dir)
	if err != nil {
		return err
	}

	last_mig_ind := 0
	for _, fl := range mig_files {
		if err := ApplySQLScript(&DbConn{User: conn.User,
			Pwd:    conn.Pwd,
			Server: conn.Server,
			Port:   conn.Port,
			DbName: conn.DbName,
		}, fl.FileName); err != nil {
			return err
		}
		fmt.Printf("applied %s migration: %s\n", mgType.Descr(), formatter(fl.Action))

		if onlyFirst {
			break
		}
		last_mig_ind++
	}
	last_mig_file_name := ""
	if mgType == MG_DOWN {
		last_mig_ind++
	}
	if last_mig_ind < len(mig_files) {
		last_mig_file_name = filepath.Base(mig_files[last_mig_ind].FileName)
	}
	if err := updateMigrationPos(last_mig_file_name); err != nil {
		return err
	}
	return nil
}

// applySQL is a helper function.
// It runs all sql scripts from updates directory on master.
// Scripts are sorted on filemtime
func applySQL(projDir, sqlDir string) error {
	db_conn, err := GetProjDbConn(projDir)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(sqlDir)
	if err != nil {
		return fmt.Errorf("os.ReadDir() failed: %v", err)
	}
	// Create a slice to store file names
	var file_names []string
	for _, file := range files {
		if !file.IsDir() {
			file_names = append(file_names, file.Name())
		}
	}

	// Sort file names by modified time
	sort.Slice(file_names, func(i, j int) bool {
		file1Path := filepath.Join(sqlDir, file_names[i])
		file2Path := filepath.Join(sqlDir, file_names[j])

		file1Info, err := os.Stat(file1Path)
		if err != nil {
			return false
		}

		file2Info, err := os.Stat(file2Path)
		if err != nil {
			return false
		}

		return file1Info.ModTime().Before(file2Info.ModTime())
	})

	for _, file := range file_names {
		if err := ApplySQLScript(db_conn, file); err != nil {
			return err
		}
	}

	return nil
}

// runAddMigration creates an empty migtation file
func runAddMigration(act string) (string, error) {
	proj_dir, err := GetProjectDir()
	if err != nil {
		return "", fmt.Errorf("GetProjectDir() failed: %v", err)
	}
	//up migration
	file_pref := dbMigrationName(time.Now(), act)
	//up
	mig_up_rel_file := GetMigrationAbsDir(filepath.Join(SQL_MIG_UP_DIR, file_pref))
	mig_up_file := filepath.Join(proj_dir, mig_up_rel_file)
	err = os.WriteFile(mig_up_file, []byte{}, NEW_FILE_PERMIS)
	if err != nil {
		return "", err
	}
	//down
	mig_down_rel_file := GetMigrationAbsDir(filepath.Join(SQL_MIG_DOWN_DIR, file_pref))
	mig_down_file := filepath.Join(proj_dir, mig_down_rel_file)
	err = os.WriteFile(mig_down_file, []byte{}, NEW_FILE_PERMIS)
	if err != nil {
		return "", err
	}
	return file_pref, nil
}

// runUpMigration executes migration up sql
func runUpMigration(onlyFirst bool) error {
	return runMigrations(SQL_MIG_UP_DIR, MG_UP, onlyFirst)
}

// runDownMigration executes migration down sql
func runDownMigration(onlyFirst bool) error {
	return runMigrations(SQL_MIG_DOWN_DIR, MG_DOWN, onlyFirst)
}

func runShowMigration() error {
	last_mig, err := lastMigrationPos()
	if err != nil {
		return err
	}
	mig_text := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("Last migration position: %s\n", mig_text(last_mig))
	return nil
}

// GetProjDbConn fetches DB connection parameters from project module file go.mod
func GetProjDbConn(projDir string) (*DbConn, error) {
	module_name, err := GetProjectModuleName(projDir)
	if err != nil {
		// return nil, err
		//if no module file, then dir name is the conf name
		dirs := strings.Split(projDir, "/")
		if len(dirs) > 0 {
			module_name = dirs[len(dirs)-1]
		} else {
			module_name = projDir
		}
	}

	//file module_name.json must exist
	config_file := module_name + CONFIG_FILE_EXT
	if _, err := os.Stat(config_file); os.IsNotExist(err) {
		return nil, fmt.Errorf("application config file %s not found", config_file)
	}

	conf := config.AppConfig{}
	if err := config.ReadConf(config_file, &conf); err != nil {
		return nil, fmt.Errorf("config.ReadConf() failed: %v", err)
	}
	db_conn, err := parseDbConn(conf.Db.Primary)
	if err != nil {
		return nil, fmt.Errorf("parseDbConn() failed: %v", err)
	}

	return db_conn, nil
}

// ApplySQLScript is a helper function.
// It runs one sql script with psql and given
// connection parameters.
func ApplySQLScript(dbConn *DbConn, scriptFile string) error {
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

// lastMigrationFileName returns last migration position
func lastMigrationFileName() string {
	return GetMigrationAbsDir(MIG_POS_FILE)
}

// lastMigrationPos returns last position
func lastMigrationPos() (string, error) {
	mig_pos_file := lastMigrationFileName()
	if exists, err := FileExists(mig_pos_file); err != nil {
		return "", err
	} else if exists {
		cur_mig, err := os.ReadFile(mig_pos_file)
		if err != nil {
			return "", err
		}
		return string(cur_mig), nil
	}
	//no position file
	return "", nil
}

// updateMigrationPos writes current migration postion to migration postion file.
func updateMigrationPos(curMigFile string) error {
	mig_pos_file := lastMigrationFileName()
	if err := os.WriteFile(mig_pos_file, []byte(curMigFile), NEW_FILE_PERMIS); err != nil {
		return err
	}
	return nil
}

// makeMigrationFileList returns migration files sorted by dates
// If fromDate parameter is set then files with
// bigger dates are returned.
func makeMigrationFileList(fromDate time.Time, migDir string, migType MigType) (MigrationFileList, error) {
	var files MigrationFileList
	err := filepath.Walk(migDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file_date, file_act, err := migrationFileParams(info.Name())
			if err != nil {
				return err
			}
			empty_time := time.Time{}
			if fromDate != empty_time &&
				(migType == MG_UP && (file_date.Equal(fromDate) || file_date.Before(fromDate))) ||
				(migType == MG_DOWN && file_date.After(fromDate)) {

				return nil
			}
			file_mig := MigrationFile{FileName: path,
				FilePos: file_date,
				Action:  file_act,
			}
			files = append(files, file_mig)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	sort.Sort(files)

	if migType == MG_DOWN {
		files_rev := make(MigrationFileList, len(files))
		copy(files_rev, files)
		for i, j := 0, len(files_rev)-1; i < j; i, j = i+1, j-1 {
			files_rev[i], files_rev[j] = files_rev[j], files_rev[i]
		}
		return files_rev, nil
	}

	return files, nil
}

// migrationFileParams returns migration position file date.
// Migration file name format: date_action.sql
// Where date is in format 20060102150405
func migrationFileParams(fileName string) (time.Time, string, error) {
	//20240101153010_ACT - take first 14 chars
	act := ""
	if len(fileName) >= 14 {
		d, err := time.Parse(MIG_FILE_DATE_FORMAT, fileName[:14])
		if err != nil {
			return time.Time{}, "", err
		}
		if len(fileName) >= 16 {
			acts := strings.Split(fileName[15:], ".")
			if len(acts) > 0 {
				act = acts[0]
			} else {
				act = fileName[15:]
			}
		} else {
			act = NO_NAME_ACTION
		}
		return d, act, nil
	}
	return time.Time{}, act, errors.New("migration file name structure error")
}
