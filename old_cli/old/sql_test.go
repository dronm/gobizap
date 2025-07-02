package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	USER_DATE_FORMAT = "2006-01-02T15:04:05"
)

type MigFile struct {
	FileName  string
	ListIndex int
	MigDate   time.Time
	Action    string
}

func TestParseDbHost(t *testing.T) {
	type SplitVal struct {
		Host string
		Port int
	}
	tests := []struct {
		Host     string
		SplitVal SplitVal
	}{
		{"localhost:5432", SplitVal{"localhost", 5432}},
		{"localhost:9999", SplitVal{"localhost", 9999}},
		{"127.0.0.1:5432", SplitVal{"127.0.0.1", 5432}},
		{"www.someresource:8888", SplitVal{"www.someresource", 8888}},
	}

	for _, tt := range tests {
		host, port, err := parseDbHost(tt.Host)
		if err != nil {
			t.Fatalf("parseDbHost() failed: %v", err)
		}
		if host != tt.SplitVal.Host {
			t.Fatalf("host expected %s, got %s", tt.SplitVal.Host, host)
		}
		if port != tt.SplitVal.Port {
			t.Fatalf("port expected %d, got %d", tt.SplitVal.Port, port)
		}
	}
}

func TestParseDbConn(t *testing.T) {
	tests := []struct {
		ConnStr string
		Conn    DbConn
	}{
		{"postgresql://userName:userPassword@localhost:5432/dbName",
			DbConn{Server: "localhost",
				Port:   5432,
				User:   "userName",
				Pwd:    "userPassword",
				DbName: "dbName",
			},
		},
	}
	for _, tt := range tests {
		conn, err := parseDbConn(tt.ConnStr)
		if err != nil {
			t.Fatalf("parseDbConn() failed: %v", err)
		}
		if conn.Server != tt.Conn.Server {
			t.Fatalf("server expected %s, got %s", tt.Conn.Server, conn.Server)
		}
		if conn.Port != tt.Conn.Port {
			t.Fatalf("Port expected %d, got %d", tt.Conn.Port, conn.Port)
		}
		if conn.User != tt.Conn.User {
			t.Fatalf("User expected %s, got %s", tt.Conn.User, conn.User)
		}
		if conn.Pwd != tt.Conn.Pwd {
			t.Fatalf("Pwd expected %s, got %s", tt.Conn.Pwd, conn.Pwd)
		}
		if conn.DbName != tt.Conn.DbName {
			t.Fatalf("DbName expected %s, got %s", tt.Conn.DbName, conn.DbName)
		}
	}
}

func TestMigFileParams(t *testing.T) {
	tests := []struct {
		FileName string
		MigDate  time.Time
		Action   string
	}{
		{"20240731122345_test.sql",
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20240731122345")
				return d
			}(),
			"test",
		},
		{"20240501235959_createTable.sql",
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20240501235959")
				return d
			}(),
			"createTable",
		},
		{"20200501235959_createTable.sql",
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20200501235959")
				return d
			}(),
			"createTable",
		},
	}
	for _, tt := range tests {
		exp_date := tt.MigDate
		dd, act, err := migrationFileParams(tt.FileName)
		if err != nil {
			t.Fatalf("migrationFileParams() failed: %v", err)
		}
		if act != tt.Action {
			t.Fatalf("action expected to be %s, got %s", tt.Action, act)
		}
		if dd != exp_date {
			t.Fatalf("date expected to be %s, got %s", exp_date.Format(USER_DATE_FORMAT), dd.Format(USER_DATE_FORMAT))
		}
	}
}

func MigrationTest(tests []MigFile, mgType MigType, expListCount int, fromDate time.Time, t *testing.T) {
	baseDir, err := os.MkdirTemp(os.TempDir(), "gobizap")
	if err != nil {
		t.Fatalf("os.MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(baseDir) // clean up

	//create migration files
	for _, tt := range tests {
		fileName := filepath.Join(baseDir, tt.FileName)
		if err := os.WriteFile(fileName, []byte("test"), FILE_PERMISSION); err != nil {
			t.Fatalf("os.WriteFile failed: %v", err)
		}
	}

	migList, err := makeMigrationFileList(fromDate, baseDir, mgType)
	if err != nil {
		t.Fatalf("makeMigrationFileList failed: %v", err)
	}
	if len(migList) != expListCount {
		t.Fatalf("Expected list count %d, got %d", expListCount, len(migList))
	}
	for i, f := range migList {
		//find with i index
		for _, tt := range tests {
			if tt.ListIndex != i {
				continue
			}
			if !f.FilePos.Equal(tt.MigDate) {
				t.Fatalf("file %s at index %d, position expected to be %s, got %s", f.FileName, i, f.FilePos.Format(USER_DATE_FORMAT), tt.MigDate.Format(USER_DATE_FORMAT))
			}
			if f.Action != tt.Action {
				t.Fatalf("file %s action expected to be %s, got %s", f.FileName, f.Action, tt.Action)
			}
			break
		}
	}
}

func TestMigUp(t *testing.T) {
	//creating files, checking list
	tests := []MigFile{
		{"20240731122345_act1.sql",
			1,
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20240731122345")
				return d
			}(),
			"act1",
		},
		{"20240501130000_act2.sql",
			0,
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20240501130000")
				return d
			}(),
			"act2",
		},
		//these files is not in the list because date is less then fromD
		{"20210731122345_act0.sql",
			-1,
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20210731122345")
				return d
			}(),
			"act0",
		},
		{"20200731122345_act0.sql",
			-1,
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20200731122345")
				return d
			}(),
			"act0",
		},
	}
	expListCount := 2
	fromDate, err := time.Parse(MIG_FILE_DATE_FORMAT, "20210731122345")
	if err != nil {
		t.Fatalf("time.Parse failed: %v", err)
	}

	MigrationTest(tests, MG_UP, expListCount, fromDate, t)
}

func TestMigDown(t *testing.T) {
	//creating files, checking list
	tests := []MigFile{
		{"20240731122345_act1.sql",
			-1,
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20240731122345")
				return d
			}(),
			"act1",
		},
		{"20240501130000_act2.sql",
			0,
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20240501130000")
				return d
			}(),
			"act2",
		},
		//these files is not in the list because date is less then fromD
		{"20210731122345_act3.sql",
			1,
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20210731122345")
				return d
			}(),
			"act3",
		},
		{"20200731122345_act4.sql",
			2,
			func() time.Time {
				d, _ := time.Parse(MIG_FILE_DATE_FORMAT, "20200731122345")
				return d
			}(),
			"act4",
		},
	}
	expListCount := 3
	fromDate, err := time.Parse(MIG_FILE_DATE_FORMAT, "20240731122345")
	if err != nil {
		t.Fatalf("time.Parse failed: %v", err)
	}

	MigrationTest(tests, MG_DOWN, expListCount, fromDate, t)
}
