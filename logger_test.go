package eddlogger

import (
	"testing"
)

type MockDriver struct {
	records []map[string]interface{}
}

func NewMockDriver() *MockDriver {
	return &MockDriver{
		records: make([]map[string]interface{}, 0),
	}
}

func (m *MockDriver) Send(record map[string]interface{}) (string, error) {
	m.records = append(m.records, record)
	return "mock-id", nil
}

func (m *MockDriver) Close() error {
	return nil
}

func TestNewLogger(t *testing.T) {
	log := NewLogger("test-service")
	if log.service != "test-service" {
		t.Errorf("Expected service 'test-service', got '%s'", log.service)
	}
}

func TestLogWithMockDriver(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	id, err := log.Log(&LogOptions{
		TraceID: "bd24e7ad-2e41-4638-b129-c1dd7e125faa",
		Action:  "TestAction",
		Context: "TestContext",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if id != "mock-id" {
		t.Errorf("Expected ID 'mock-id', got '%s'", id)
	}
	if len(mockDriver.records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(mockDriver.records))
	}
}

func TestLogInfoLevel(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	log.Log(&LogOptions{
		TraceID:  "test-trace-id",
		Level:    "INFO",
		Action:   "test_action",
		Context:  "TestContext",
		DurationMs: 123.45,
	})

	record := mockDriver.records[0]
	if record["traceId"] != "test-trace-id" {
		t.Errorf("Expected traceId 'test-trace-id', got '%v'", record["traceId"])
	}
	if record["level"] != "INFO" {
		t.Errorf("Expected level 'INFO', got '%v'", record["level"])
	}
}

func TestLogDefaultInfoLevel(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	log.Log(&LogOptions{
		TraceID: "test-trace-id",
		Action:  "test_action",
	})

	record := mockDriver.records[0]
	if record["level"] != "INFO" {
		t.Errorf("Expected default level 'INFO', got '%v'", record["level"])
	}
}

func TestLogWithRequestResponse(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	log.Log(&LogOptions{
		TraceID:     "req-trace-id",
		Method:      "POST",
		Path:        "/api/test",
		RequestHeaders: map[string]string{"Content-Type": "application/json"},
		RequestBody: map[string]interface{}{"key": "value"},
		StatusCode:  200,
	})

	record := mockDriver.records[0]
	if record["traceId"] != "req-trace-id" {
		t.Errorf("Expected traceId 'req-trace-id', got '%v'", record["traceId"])
	}
}

func TestLogServiceOverride(t *testing.T) {
	log := NewLogger("default-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	log.Log(&LogOptions{
		TraceID: "svc-trace-id",
		Service: "overridden-service",
	})

	record := mockDriver.records[0]
	if record["service"] != "overridden-service" {
		t.Errorf("Expected service 'overridden-service', got '%v'", record["service"])
	}
}

func TestSetDriver(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	if log.driver != mockDriver {
		t.Errorf("Expected driver to be set")
	}
}

func TestCloseWithDriver(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	err := log.Close()
	if err != nil {
		t.Errorf("Expected no error on close, got %v", err)
	}
}

func TestNewLoggerEmptyService(t *testing.T) {
	log := NewLogger("")
	if log.service != "digital-edd" {
		t.Errorf("Expected default service 'digital-edd', got '%s'", log.service)
	}

	defer log.Close()
	driver := log.getDriver()
	if driver == nil {
		t.Fatal("Expected a driver to be created")
	}
}

func TestSendTraceLog(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	id, err := log.SendTraceLog(&TraceLog{
		TraceID: "test-trace-id",
		Service: "test-service",
		Level:   INFO,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if id != "mock-id" {
		t.Errorf("Expected ID 'mock-id', got '%s'", id)
	}
	record := mockDriver.records[0]
	if record["traceId"] != "test-trace-id" {
		t.Errorf("Expected traceId 'test-trace-id', got '%v'", record["traceId"])
	}
}

func TestGetMexicoTimeAsUTC(t *testing.T) {
	now := GetMexicoTimeAsUTC()
	if now == "" {
		t.Errorf("Expected non-empty time string")
	}
}

func TestProductionCheck(t *testing.T) {
	result := IsProduction()
	if result {
		t.Skip("Running in production mode (PUBSUB_ENABLED set), skipping local only test")
	}
}
