package health

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/rs/zerolog"
)

func TestHealthHandler_Healthy(t *testing.T) {
	// Mock DB
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectPing()

	// Mock Redis
	rdb, mockRedis := redismock.NewClientMock()
	mockRedis.ExpectPing().SetVal("PONG")

	logger := zerolog.Nop()
	handler := NewHandler(db, rdb, logger)

	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Status != "healthy" {
		t.Errorf("expected overall status healthy, got %s", resp.Status)
	}
	if resp.Checks["postgres"].Status != "healthy" {
		t.Errorf("expected postgres status healthy, got %s", resp.Checks["postgres"].Status)
	}
	if resp.Checks["redis"].Status != "healthy" {
		t.Errorf("expected redis status healthy, got %s", resp.Checks["redis"].Status)
	}
}

func TestHealthHandler_Unhealthy(t *testing.T) {
	// Mock DB
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectPing().WillReturnError(errors.New("db down"))

	// Mock Redis
	rdb, mockRedis := redismock.NewClientMock()
	mockRedis.ExpectPing().SetErr(errors.New("redis down"))

	logger := zerolog.Nop()
	handler := NewHandler(db, rdb, logger)

	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	var resp Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Status != "unhealthy" {
		t.Errorf("expected overall status unhealthy, got %s", resp.Status)
	}
	if resp.Checks["postgres"].Status != "unhealthy" {
		t.Errorf("expected postgres status unhealthy, got %s", resp.Checks["postgres"].Status)
	}
	if resp.Checks["redis"].Status != "unhealthy" {
		t.Errorf("expected redis status unhealthy, got %s", resp.Checks["redis"].Status)
	}
}
