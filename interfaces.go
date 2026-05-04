package eddlogger

type LogLevel string

const (
	DEBUG    LogLevel = "DEBUG"
	INFO     LogLevel = "INFO"
	NOTICE   LogLevel = "NOTICE"
	WARNING  LogLevel = "WARNING"
	ERROR    LogLevel = "ERROR"
	CRITICAL LogLevel = "CRITICAL"
	ALERT    LogLevel = "ALERT"
)

type RequestInfo struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body,omitempty"`
}

type ResponseInfo struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       interface{}       `json:"body,omitempty"`
}

type TraceLog struct {
	TypeStream  string        `json:"typeStream"` // "sdkHisStream"
	TraceID     string        `json:"traceId"`
	Timestamp   string        `json:"timestamp"`
	Service     string        `json:"service"`
	Level       LogLevel      `json:"level"`
	Action      string        `json:"action"`
	Context     string        `json:"context,omitempty"`
	Request     *RequestInfo  `json:"request,omitempty"`
	Response    *ResponseInfo `json:"response,omitempty"`
	MessageInfo string        `json:"messageInfo,omitempty"`
	MessageRaw  string        `json:"messageRaw,omitempty"`
	DurationMs  float64       `json:"durationMs,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
}

// TraceSDKTrn captures all detail-level data from the ms-facade algorithm pipeline.
// One row = one SKU winner route evaluated in one request.
// Relates to LGS_EDD_SDK_HIS via requestId = traceId.
type TraceSDKTrn struct {
	TypeStream string `json:"typeStream"` // "sdkTrnStream"
	RequestID  string `json:"requestId"`
	RequestType string `json:"requestType"`
	Endpoint   string `json:"endpoint"`

	CP              string `json:"cp"`
	Channel         string `json:"channel"`
	EnterpriseCode  string `json:"enterpriseCode"`
	SKU             string `json:"sku"`
	Quantity        int    `json:"quantity"`
	ProductType     string `json:"productType"`
	FulfillmentType string `json:"fulfillmentType"`

	PurchaseDateEdd1 string `json:"purchaseDateEdd1,omitempty"`
	DeliveryDateEdd2 string `json:"deliveryDateEdd2,omitempty"`
	StoreRejected    string `json:"storeRejected,omitempty"`

	OrderNumber string `json:"orderNumber,omitempty"`
	EmittedAt   string `json:"emittedAt,omitempty"`

	DeliveryDate   string  `json:"deliveryDate"`
	DeliveryMethod string  `json:"deliveryMethod"`
	Route          string  `json:"route"`
	StoreID        string  `json:"storeId"`
	StoreName      string  `json:"storeName"`
	TimeDays       int     `json:"timeDays"`
	Cost           float64 `json:"cost"`

	OhLimit        int     `json:"ohLimit"`
	OhInventory    int     `json:"ohInventory"`
	StoreCapacity  float64 `json:"storeCapacity"`
	Winner         bool    `json:"winner"`
	StoreInventory int     `json:"storeInventory"`
	SkusSent       int     `json:"skusSent"`

	IsRecalculated bool `json:"isRecalculated"`

	InventoryWeight float64 `json:"inventoryWeight"`
	TimeWeight      float64 `json:"timeWeight"`
	CostWeight      float64 `json:"costWeight"`
	NodeWeight      float64 `json:"nodeWeight"`
	Season          string  `json:"season"`
	Score           float64 `json:"score"`

	CalculatedAt    string  `json:"calculatedAt"`
	ExecutionTimeMs float64 `json:"executionTimeMs"`

	IngestedAt  string `json:"ingestedAt"`
	ServiceName string `json:"serviceName"`
}

// TraceRoute is a deprecated alias for TraceSDKTrn, kept for backward compatibility.
type TraceRoute = TraceSDKTrn

