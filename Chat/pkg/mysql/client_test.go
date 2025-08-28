package mysql_test

import (
	"ChatRoom/chat/pkg/mysql"
	"testing"
)

func TestNewClient(t *testing.T) {
	cfg := mysql.Config{
		User:     "root",
		Password: "flyanytime0916",
		Host:     "localhost",
		Port:     3306,
		Database: "blog",
	}

	db, err := mysql.NewMysqlClient(&cfg)
	if err != nil {
		t.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("MySQL ping failed: %v", err)
	}
}
