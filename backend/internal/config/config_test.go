package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear any existing env vars for the test
	os.Clearenv()

	cfg := Load()

	if cfg.BackendPort != "8080" {
		t.Errorf("Expected default BackendPort '8080', got '%s'", cfg.BackendPort)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("Expected default DBHost 'localhost', got '%s'", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("Expected default DBPort '5432', got '%s'", cfg.DBPort)
	}
	if cfg.DBUser != "flagmgmt" {
		t.Errorf("Expected default DBUser 'flagmgmt', got '%s'", cfg.DBUser)
	}
	if cfg.DBPassword != "flagmgmt_dev" {
		t.Errorf("Expected default DBPassword 'flagmgmt_dev', got '%s'", cfg.DBPassword)
	}
	if cfg.DBName != "flagmanagment" {
		t.Errorf("Expected default DBName 'flagmanagment', got '%s'", cfg.DBName)
	}
	if cfg.RedisHost != "localhost" {
		t.Errorf("Expected default RedisHost 'localhost', got '%s'", cfg.RedisHost)
	}
	if cfg.RedisPort != "6379" {
		t.Errorf("Expected default RedisPort '6379', got '%s'", cfg.RedisPort)
	}
	if cfg.Env != "development" {
		t.Errorf("Expected default Env 'development', got '%s'", cfg.Env)
	}
	if cfg.LogFormat != "auto" {
		t.Errorf("Expected default LogFormat 'auto', got '%s'", cfg.LogFormat)
	}
}

func TestLoad_Overrides(t *testing.T) {
	os.Clearenv()
	t.Setenv("FM_BACKEND_PORT", "9090")
	t.Setenv("FM_DB_NAME", "testdb")

	cfg := Load()

	if cfg.BackendPort != "9090" {
		t.Errorf("Expected overridden BackendPort '9090', got '%s'", cfg.BackendPort)
	}
	if cfg.DBName != "testdb" {
		t.Errorf("Expected overridden DBName 'testdb', got '%s'", cfg.DBName)
	}
	// Verify another default is intact
	if cfg.DBHost != "localhost" {
		t.Errorf("Expected default DBHost 'localhost', got '%s'", cfg.DBHost)
	}
}
