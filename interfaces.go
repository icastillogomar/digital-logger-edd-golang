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

// schema_input.json

type EddLine struct {
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
	ProductType string `json:"productType"`
}

type RecalculateLine struct {
	SKU              string  `json:"sku"`
	Quantity         *int    `json:"quantity,omitempty"`
	PurchaseDateEdd1 string  `json:"purchaseDateEdd1"`
	DeliveryDateEdd2 string  `json:"deliveryDateEdd2"`
	StoreRejected    *string `json:"storeRejected,omitempty"`
	CarrierRejected  *string `json:"carrierRejected,omitempty"`
}

type TraceByInput struct {
	RequestID        string           `json:"requestId"`
	RequestType      string           `json:"requestType"`
	Endpoint         string           `json:"endpoint"`
	ReceivedAt       string           `json:"receivedAt"`
	EnterpriseCode   string           `json:"enterpriseCode"`
	Cp               string           `json:"cp"`
	Channel          string           `json:"channel"`
	EddLines         *EddLine         `json:"eddLines,omitempty"`
	RecalculateLines *RecalculateLine `json:"recalculateLines,omitempty"`
	LineCount        *int             `json:"lineCount,omitempty"`
	Tags             []string         `json:"tags,omitempty"`
	AdditionalData   interface{}      `json:"additionalData,omitempty"`
	IngestedAt       string           `json:"ingestedAt"`
}

// schema_output.json

type OutputMetadata struct {
	IDTxn            *string `json:"idTxn,omitempty"`
	ProcessingTimeMs *int    `json:"processingTimeMs,omitempty"`
	IngestedAt       *string `json:"ingestedAt,omitempty"`
	RecalculateOrder *string `json:"recalculateOrder,omitempty"`
}

type AlgorithmWeights struct {
	Inventory  *float64 `json:"inventory,omitempty"`
	LeadTime   *float64 `json:"leadTime,omitempty"`
	Cost       *float64 `json:"cost,omitempty"`
	Node       *float64 `json:"node,omitempty"`
	Path       *float64 `json:"path,omitempty"`
	Difference *float64 `json:"difference,omitempty"`
	Splits     *float64 `json:"splits,omitempty"`
}

type EddCalculatedSummary struct {
	Split           *bool    `json:"split,omitempty"`
	ProductType     *string  `json:"productType,omitempty"`
	MaxDeliveryDays *int     `json:"maxDeliveryDays,omitempty"`
	UsedRoutes      *int     `json:"usedRoutes,omitempty"`
	StoreSelected   *string  `json:"storeSelected,omitempty"`
	StoreName       *string  `json:"storeName,omitempty"`
	Edd1            *string  `json:"edd1,omitempty"`
	TotalCost       *float64 `json:"totalCost,omitempty"`
	Plan            *string  `json:"plan,omitempty"`
	ErrorCode       *string  `json:"errorCode,omitempty"`
	ErrorMessage    *string  `json:"errorMessage,omitempty"`
}

type EddCalculatedRoute struct {
	IDRoute         *string  `json:"idRoute,omitempty"`
	Quantity        *int     `json:"quantity,omitempty"`
	DeliveryDate    *string  `json:"deliveryDate,omitempty"`
	TimeDays        *int     `json:"timeDays,omitempty"`
	Cost            *float64 `json:"cost,omitempty"`
	DeliveryMethod  *string  `json:"deliveryMethod,omitempty"`
	IDCarrier       *int     `json:"idCarrier,omitempty"`
	StoreID         *int     `json:"storeId,omitempty"`
	StoreName       *string  `json:"storeName,omitempty"`
	StoreCapacity   *int     `json:"storeCapacity,omitempty"`
	Inventory       *int     `json:"inventory,omitempty"`
	AvailableStores *int     `json:"availableStores,omitempty"`
}

type EddCalculated struct {
	SKU     string                 `json:"sku"`
	Summary []EddCalculatedSummary `json:"summary"`
	Routes  []EddCalculatedRoute   `json:"routes"`
}

type TraceByOutput struct {
	RequestID           string            `json:"requestId"`
	RequestType         string            `json:"requestType"`
	Endpoint            string            `json:"endpoint"`
	RespondedAt         string            `json:"respondedAt"`
	HTTPStatusCode      int               `json:"httpStatusCode"`
	StatusFamily        *int              `json:"statusFamily,omitempty"`
	IsError             *bool             `json:"isError,omitempty"`
	Metadata            *OutputMetadata   `json:"metadata,omitempty"`
	AlgorithmModelState *string           `json:"algorithmModelState,omitempty"`
	AlgorithmWeights    *AlgorithmWeights `json:"algorithmWeights,omitempty"`
	EddCalculated       *EddCalculated    `json:"eddCalculated,omitempty"`
	StoreIDs            []int             `json:"storeIds,omitempty"`
	ErrorCode           *string           `json:"errorCode,omitempty"`
	ErrorMessage        *string           `json:"errorMessage,omitempty"`
	Tags                []string          `json:"tags,omitempty"`
	AdditionalData      interface{}       `json:"additionalData,omitempty"`
	IngestedAt          string            `json:"ingestedAt"`
}

// schema_logs.json

type LogRequestPayload struct {
	Method string      `json:"method"`
	Body   interface{} `json:"body,omitempty"`
}

type LogResponsePayload struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body,omitempty"`
}

type TraceByLog struct {
	LogID          string              `json:"logId"`
	RequestID      string              `json:"requestId"`
	RequestType    string              `json:"requestType"`
	Endpoint       string              `json:"endpoint"`
	LogAt          string              `json:"logAt"`
	Level          *string             `json:"level,omitempty"`
	Context        *string             `json:"context,omitempty"`
	Message        *string             `json:"message,omitempty"`
	Step           *string             `json:"step,omitempty"`
	DurationMs     *float64            `json:"durationMs,omitempty"`
	IDTxn          string              `json:"idTxn"`
	Tags           []string            `json:"tags,omitempty"`
	AdditionalData interface{}         `json:"additionalData,omitempty"`
	Extra          interface{}         `json:"extra,omitempty"`
	Stacktrace     *string             `json:"stacktrace,omitempty"`
	IngestedAt     string              `json:"ingestedAt"`
	ServiceName    *string             `json:"serviceName,omitempty"`
	Request        *LogRequestPayload  `json:"request,omitempty"`
	Response       *LogResponsePayload `json:"response,omitempty"`
}

//map[string]interface{}
