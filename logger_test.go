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

	id, err := log.SendTraceByLog(&TraceLogOptions{
		LogID:       "bd24e7ad-2e41-4638-b129-c1dd7e125faa",
		RequestType: "HTTP",
		Context:     "TestContext",
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
	if record["logId"] != "bd24e7ad-2e41-4638-b129-c1dd7e125faa" {
		t.Errorf("Expected traceId 'bd24e7ad-2e41-4638-b129-c1dd7e125faa', got '%v'", record["logId"])
	}
	if record["requestType"] != "HTTP" {
		t.Errorf("Expected action 'HTTP', got '%v'", record["requestType"])
	}
	if record["context"] != "TestContext" {
		t.Errorf("Expected context 'TestContext', got '%v'", record["context"])
	}
}

func TestLogWithRequestResponse(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	id, err := log.SendTraceByLog(&TraceLogOptions{
		LogID:              "bd24e7ad-2e41-4638-b129-c1dd7e125faa",
		RequestType:        "HTTP",
		Endpoint:           "/edd/fee2/facade/pdp",
		ServiceName:        "my-service",
		RequestMethod:      "GET",
		RequestBody:        map[string]interface{}{},
		ResponseStatusCode: 200,
		ResponseBody:       map[string]interface{}{"success": true},
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
	if request["method"] != "GET" {
		t.Errorf("Expected method 'GET', got '%v'", request["method"])
	}

	response, ok := record["response"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected response to be a map")
	}
	if response["statusCode"] != float64(200) {
		t.Errorf("Expected statusCode 200, got '%v'", response["statusCode"])
	}
}

func TestLogWithTags(t *testing.T) {
	log := NewLogger("test-service")
	mockDriver := NewMockDriver()
	log.SetDriver(mockDriver)

	tags := []string{"http", "middleware", "request-response"}
	_, err := log.SendTraceByLog(&TraceLogOptions{
		LogID:       "bd24e7ad-2e41-4638-b129-c1dd7e125faa",
		RequestType: "HTTP",
		Tags:        tags,
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

	_, err := log.SendTraceByLog(&TraceLogOptions{
		LogID: "bd24e7ad-2e41-4638-b129-c1dd7e125faa",
		Level: "INFO",
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

	_, err := log.SendTraceByLog(&TraceLogOptions{
		LogID: "bd24e7ad-2e41-4638-b129-c1dd7e125faa",
		Level: "ERROR",
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
		"LogID":    "bd24e7ad-2e41-4638-b129-c1dd7e125faa",
		"Endpoint": "/edd/fee2/facade/pdp",
	}

	id, err := driver.Send(record)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != "console-log" {
		t.Errorf("Expected id 'console-log', got '%s'", id)
	}
}

// Testing
//❯ GOCACHE=$(pwd)/.gocache go test -run TestLog -v
//=== RUN   TestLogWithMockDriver
//--- PASS: TestLogWithMockDriver (0.00s)
//=== RUN   TestLogWithRequestResponse
//--- PASS: TestLogWithRequestResponse (0.00s)
//=== RUN   TestLogWithTags
//--- PASS: TestLogWithTags (0.00s)
//=== RUN   TestLogDefaultLevel
//--- PASS: TestLogDefaultLevel (0.00s)
//=== RUN   TestLogCustomLevel
//--- PASS: TestLogCustomLevel (0.00s)
//PASS
//ok      github.com/icastillogomar/digital-logger-edd-golang     0.803s
//
//GOCACHE=$(pwd)/.gocache go test -run TestLogWithMockDriver -v
