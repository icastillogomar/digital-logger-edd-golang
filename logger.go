package eddlogger

import (
	"encoding/json"
	"sync"

	"github.com/icastillogomar/digital-logger-edd-golang/drivers"
)

type EddLogger struct {
	service string
	driver  drivers.BaseDriver
}

type LogOptions struct {
	TraceID         string
	Level           string
	Action          string
	Context         string
	Method          string
	Path            string
	RequestHeaders  map[string]string
	RequestBody     interface{}
	StatusCode      int
	ResponseHeaders map[string]string
	ResponseBody    interface{}
	MessageInfo     string
	MessageRaw      string
	DurationMs      float64
	Tags            []string
	Service         string
}

func NewLogger(service string) *EddLogger {
	if service == "" {
		service = "digital-edd"
	}
	return &EddLogger{
		service: service,
	}
}

func (l *EddLogger) getDriver() drivers.BaseDriver {
	if l.driver != nil {
		return l.driver
	}
	l.driver = l.createDriver()
	return l.driver
}

func (l *EddLogger) createDriver() drivers.BaseDriver {
	if IsProduction() {
		driver, err := drivers.NewPubSubDriver("", "")
		if err != nil {
			LogError("No se pudo inicializar PubSubDriver: " + err.Error())
			LogWarning("Usando ConsoleDriver como fallback")
			return drivers.NewConsoleDriver()
		}
		return driver
	}

	driver, err := drivers.NewPostgresDriver("")
	if err != nil {
		LogError("No se pudo inicializar PostgresDriver: " + err.Error())
		LogWarning("Usando ConsoleDriver como fallback")
		return drivers.NewConsoleDriver()
	}
	return driver
}

func (l *EddLogger) SetDriver(driver drivers.BaseDriver) {
	l.driver = driver
}

func (l *EddLogger) SendTraceLog(trace *TraceLog) (string, error) {
	data, err := json.Marshal(trace)
	if err != nil {
		return "", err
	}

	var record map[string]interface{}
	if err := json.Unmarshal(data, &record); err != nil {
		return "", err
	}

	return l.getDriver().Send(record)
}

func (l *EddLogger) Log(opts *LogOptions) (string, error) {
	if opts == nil {
		opts = &LogOptions{}
	}

	level := opts.Level
	if level == "" {
		level = string(INFO)
	}

	var request *RequestInfo
	if opts.Method != "" && opts.Path != "" {
		request = &RequestInfo{
			Method:  opts.Method,
			Path:    opts.Path,
			Headers: opts.RequestHeaders,
			Body:    opts.RequestBody,
		}
	}

	var response *ResponseInfo
	if opts.StatusCode != 0 {
		response = &ResponseInfo{
			StatusCode: opts.StatusCode,
			Headers:    opts.ResponseHeaders,
			Body:       opts.ResponseBody,
		}
	}

	service := opts.Service
	if service == "" {
		service = l.service
	}

	trace := &TraceLog{
		TypeStream:  "sdkHisStream",
		TraceID:     opts.TraceID,
		Timestamp:   GetMexicoTimeAsUTC(),
		Service:     service,
		Level:       LogLevel(level),
		Action:      opts.Action,
		Context:     opts.Context,
		Request:     request,
		Response:    response,
		MessageInfo: opts.MessageInfo,
		MessageRaw:  opts.MessageRaw,
		DurationMs:  opts.DurationMs,
		Tags:        opts.Tags,
	}
	return l.SendTraceLog(trace)
}

func (l *EddLogger) Close() error {
	if l.driver != nil {
		return l.driver.Close()
	}
	return nil
}

func (l *EddLogger) sendTraceAll(trace interface{}) (string, error) {
	data, err := json.Marshal(trace)
	if err != nil {
		return "", err
	}

	var record map[string]interface{}
	if err := json.Unmarshal(data, &record); err != nil {
		return "", err
	}

	return l.getDriver().Send(record)
}

// SDKTrnOptions collects detail-level data from the ms-facade algorithm pipeline.
// One Add() call = one SKU winner route. Writes to LGS_EDD_SDK_TRN.
type SDKTrnOptions struct {
	RequestType string
	Endpoint    string

	CP              string
	Channel         string
	EnterpriseCode  string
	SKU             string
	Quantity        int
	ProductType     string
	FulfillmentType string

	PurchaseDateEdd1 string
	DeliveryDateEdd2 string
	StoreRejected    string

	OrderNumber string
	EmittedAt   string

	DeliveryDate   string
	DeliveryMethod string
	Route          string
	StoreID        string
	StoreName      string
	TimeDays       int
	Cost           float64

	OhLimit        int
	OhInventory    int
	StoreCapacity  float64
	StoreInventory int
	SkusSent       int

	IsRecalculated bool

	InventoryWeight float64
	TimeWeight      float64
	CostWeight      float64
	NodeWeight      float64
	Season          string
	Score           float64

	CalculatedAt    string
	ExecutionTimeMs float64

	ServiceName string
}

// SDKTrnCollector accumulates SDK_TRN rows during the ms-facade algorithm lifecycle.
// One Add() call = one winner route for one SKU. Call Save() once at the end.
// Safe for concurrent use.
type SDKTrnCollector struct {
	mu          sync.Mutex
	logger      *EddLogger
	routes      []TraceSDKTrn
	requestID   string
	serviceName string
	startedAt   string
	saved       bool
}

// NewSDKTrnCollector creates a collector bound to a single request lifecycle.
func (l *EddLogger) NewSDKTrnCollector(requestID string) *SDKTrnCollector {
	return &SDKTrnCollector{
		logger:    l,
		requestID: requestID,
		startedAt: GetMexicoTimeAsUTC(),
	}
}

// Add appends the winner route for one SKU. Thread-safe.
func (rc *SDKTrnCollector) Add(opts *SDKTrnOptions) {
	if opts == nil {
		return
	}

	serviceName := opts.ServiceName
	if serviceName == "" {
		serviceName = rc.logger.service
	}

	rc.mu.Lock()
	if rc.serviceName == "" {
		rc.serviceName = serviceName
	}
	rc.mu.Unlock()

	now := GetMexicoTimeAsUTC()

	trace := TraceSDKTrn{
		TypeStream: "sdkTrnStream",

		RequestID:   rc.requestID,
		RequestType: opts.RequestType,
		Endpoint:    opts.Endpoint,

		CP:              opts.CP,
		Channel:         opts.Channel,
		EnterpriseCode:  opts.EnterpriseCode,
		SKU:             opts.SKU,
		Quantity:        opts.Quantity,
		ProductType:     opts.ProductType,
		FulfillmentType: opts.FulfillmentType,

		PurchaseDateEdd1: opts.PurchaseDateEdd1,
		DeliveryDateEdd2: opts.DeliveryDateEdd2,
		StoreRejected:    opts.StoreRejected,

		OrderNumber: opts.OrderNumber,
		EmittedAt:   opts.EmittedAt,

		DeliveryDate:   opts.DeliveryDate,
		DeliveryMethod: opts.DeliveryMethod,
		Route:          opts.Route,
		StoreID:        opts.StoreID,
		StoreName:      opts.StoreName,
		TimeDays:       opts.TimeDays,
		Cost:           opts.Cost,

		OhLimit:        opts.OhLimit,
		OhInventory:    opts.OhInventory,
		StoreCapacity:  opts.StoreCapacity,
		Winner:         true,
		StoreInventory: opts.StoreInventory,
		SkusSent:       opts.SkusSent,

		IsRecalculated: opts.IsRecalculated,

		InventoryWeight: opts.InventoryWeight,
		TimeWeight:      opts.TimeWeight,
		CostWeight:      opts.CostWeight,
		NodeWeight:      opts.NodeWeight,
		Season:          opts.Season,
		Score:           opts.Score,

		CalculatedAt:    opts.CalculatedAt,
		ExecutionTimeMs: opts.ExecutionTimeMs,

		IngestedAt:  now,
		ServiceName: rc.serviceName,
	}

	rc.mu.Lock()
	rc.routes = append(rc.routes, trace)
	rc.mu.Unlock()
}

// Save flushes all collected records through the configured driver. Idempotent.
func (rc *SDKTrnCollector) Save() (int, error) {
	rc.mu.Lock()
	if rc.saved {
		rc.mu.Unlock()
		return 0, nil
	}
	rc.saved = true
	routes := rc.routes
	rc.mu.Unlock()

	var lastErr error
	for i := range routes {
		if _, err := rc.logger.sendTraceAll(&routes[i]); err != nil {
			lastErr = err
		}
	}
	return len(routes), lastErr
}

// Len returns the number of routes collected so far.
func (rc *SDKTrnCollector) Len() int {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return len(rc.routes)
}

// Deprecated aliases for backward compatibility.
type TraceRouteOptions = SDKTrnOptions
type RouteCollector = SDKTrnCollector

func (l *EddLogger) NewRouteCollector(requestID string) *SDKTrnCollector {
	return l.NewSDKTrnCollector(requestID)
}
