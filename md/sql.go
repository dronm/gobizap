package md

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	CONN_INVALID_FORMAT = "invalid connection string format"
)

type DbConn struct {
	User   string
	Pwd    string
	Server string
	Port   int
	DbName string
}

// NewDbConn is a helper function to parse db connection string
// from template postgresql://USER:PWD@SERVER:PORT/DB_NAME
// It makes a structure where all fields are separate.
func NewDbConn(conn string) (*DbConn, error) {
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

	host, port, err := ParseDbHost(host_db[0])
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

// ApplySQLScript is a helper function.
// It runs one sql script with psql and given
// connection parameters.
func (db *DbConn) ApplySQLScript(scriptFile string) error {
	if db.Pwd != "" {
		if err := os.Setenv("PGPASSWORD", db.Pwd); err != nil {
			return fmt.Errorf("os.Setenv() failed: %v", err)
		}
	}
	bash_cmd := fmt.Sprintf("psql -h %s -p %d -d %s -U %s -f %s",
		db.Server,
		db.Port,
		db.DbName,
		db.User,
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

// ParseDbHost parses host in SERVER:PORT format,
// splits SERVER and PORT to different variables.
func ParseDbHost(host string) (string, int, error) {
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
