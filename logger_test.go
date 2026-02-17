package eddlogger

import (
	"testing"

	"github.com/icastillogomar/digital-logger-edd-golang/drivers"
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
		TraceID: "test-trace-001",
		Action:  "TEST_ACTION",
		Context: "TestContext",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != "mock-id" {
		t.Errorf("Expected id 'mock-id', got '%s'", id)
	}

	if len(mockDriver.records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(mockDriver.records))
	}

	record := mockDriver.records[0]
	if record["traceId"] != "test-trace-001" {
		t.Errorf("Expected traceId 'test-trace-001', got '%v'", record["traceId"])
	}
	if record["action"] != "TEST_ACTION" {
		t.Errorf("Expected action 'TEST_ACTION', got '%v'", record["action"])
	}
	if record["context"] != "TestContext" {
		t.Errorf("Expected context 'TestContext', got '%v'", record["context"])
	}
}

func TestLogWithRequestResponse(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	id, err := log.Log(&LogOptions{
		TraceID:      "test-trace-002",
		Action:       "API_CALL",
		Method:       "POST",
		Path:         "/api/test",
		RequestBody:  map[string]interface{}{"key": "value"},
		StatusCode:   200,
		ResponseBody: map[string]interface{}{"result": "success"},
		DurationMs:   123.45,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != "mock-id" {
		t.Errorf("Expected id 'mock-id', got '%s'", id)
	}

	record := mockDriver.records[0]
	
	request, ok := record["request"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected request to be a map")
	}
	if request["method"] != "POST" {
		t.Errorf("Expected method 'POST', got '%v'", request["method"])
	}
	if request["path"] != "/api/test" {
		t.Errorf("Expected path '/api/test', got '%v'", request["path"])
	}

	response, ok := record["response"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected response to be a map")
	}
	if response["statusCode"] != float64(200) {
		t.Errorf("Expected statusCode 200, got '%v'", response["statusCode"])
	}

	if record["durationMs"] != 123.45 {
		t.Errorf("Expected durationMs 123.45, got '%v'", record["durationMs"])
	}
}

func TestLogWithTags(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	tags := []string{"tag1", "tag2", "tag3"}
	_, err := log.Log(&LogOptions{
		TraceID: "test-trace-003",
		Action:  "TAGGED_ACTION",
		Tags:    tags,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	record := mockDriver.records[0]
	recordTags, ok := record["tags"].([]interface{})
	if !ok {
		t.Fatal("Expected tags to be an array")
	}

	if len(recordTags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(recordTags))
	}
}

func TestLogDefaultLevel(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	_, err := log.Log(&LogOptions{
		TraceID: "test-trace-004",
		Action:  "DEFAULT_LEVEL",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	record := mockDriver.records[0]
	if record["level"] != "INFO" {
		t.Errorf("Expected default level 'INFO', got '%v'", record["level"])
	}
}

func TestLogCustomLevel(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	_, err := log.Log(&LogOptions{
		TraceID: "test-trace-005",
		Action:  "ERROR_ACTION",
		Level:   "ERROR",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	record := mockDriver.records[0]
	if record["level"] != "ERROR" {
		t.Errorf("Expected level 'ERROR', got '%v'", record["level"])
	}
}

func TestConsoleDriver(t *testing.T) {
	driver := drivers.NewConsoleDriver()
	defer driver.Close()

	record := map[string]interface{}{
		"traceId": "console-test",
		"action":  "CONSOLE_TEST",
	}

	id, err := driver.Send(record)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != "console-log" {
		t.Errorf("Expected id 'console-log', got '%s'", id)
	}
}
