package md

import (
	"testing"
)

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
		host, port, err := ParseDbHost(tt.Host)
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

func TestDbConn(t *testing.T) {
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
		conn, err := NewDbConn(tt.ConnStr)
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
