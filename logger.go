package eddlogger

import (
	"encoding/json"

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

type TraceInputOptions struct {
	RequestID                    string
	RequestType                  string
	Endpoint                     string
	ReceivedAt                   string
	EnterpriseCode               string
	Cp                           string
	Channel                      string
	EddLineSKU                   string
	EddLineQuantity              int
	EddLineProductType           string
	RecalculateLineSKU           string
	RecalculateLineQuantity      *int
	RecalculateLinePurchaseDate  string
	RecalculateLineDeliveryDate  string
	RecalculateLineStoreRejected string
	RecalculateLineCarrierReject string
	LineCount                    *int
	Tags                         []string
	AdditionalData               interface{}
	IngestedAt                   string
}

type TraceOutputOptions struct {
	RequestID                  string
	RequestType                string
	Endpoint                   string
	RespondedAt                string
	HTTPStatusCode             int
	StatusFamily               *int
	IsError                    *bool
	MetadataIDTxn              string
	MetadataProcessingTimeMs   *int
	MetadataIngestedAt         string
	MetadataRecalculateOrder   string
	AlgorithmModelState        string
	AlgorithmWeightsInventory  *float64
	AlgorithmWeightsLeadTime   *float64
	AlgorithmWeightsCost       *float64
	AlgorithmWeightsNode       *float64
	AlgorithmWeightsPath       *float64
	AlgorithmWeightsDifference *float64
	AlgorithmWeightsSplits     *float64
	EddCalculated              *EddCalculated
	StoreIDs                   []int
	ErrorCode                  string
	ErrorMessage               string
	Tags                       []string
	AdditionalData             interface{}
	IngestedAt                 string
}

type TraceLogOptions struct {
	LogID              string
	RequestID          string
	RequestType        string
	Endpoint           string
	LogAt              string
	Level              string
	Context            string
	Message            string
	Step               string
	DurationMs         *float64
	IDTxn              string
	Tags               []string
	AdditionalData     interface{}
	Extra              interface{}
	Stacktrace         string
	IngestedAt         string
	ServiceName        string
	RequestMethod      string
	RequestBody        interface{}
	ResponseStatusCode int
	ResponseBody       interface{}
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

// TODO: Mar9,2026 nuevas funciones para los schema de fee2 sdk

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

// SendTraceByInput is kept for backward compatibility and delegates to SendTraceAll.
func (l *EddLogger) SendTraceByInput(opts *TraceInputOptions) (string, error) {
	if opts == nil {
		opts = &TraceInputOptions{}
	}

	var eddLine *EddLine
	if opts.EddLineSKU != "" || opts.EddLineQuantity != 0 || opts.EddLineProductType != "" {
		eddLine = &EddLine{
			SKU:         opts.EddLineSKU,
			Quantity:    opts.EddLineQuantity,
			ProductType: opts.EddLineProductType,
		}
	}

	var recalculateLine *RecalculateLine
	if opts.RecalculateLineSKU != "" ||
		opts.RecalculateLineQuantity != nil ||
		opts.RecalculateLinePurchaseDate != "" ||
		opts.RecalculateLineDeliveryDate != "" ||
		opts.RecalculateLineStoreRejected != "" ||
		opts.RecalculateLineCarrierReject != "" {
		recalculateLine = &RecalculateLine{
			SKU:              opts.RecalculateLineSKU,
			Quantity:         opts.RecalculateLineQuantity,
			PurchaseDateEdd1: opts.RecalculateLinePurchaseDate,
			DeliveryDateEdd2: opts.RecalculateLineDeliveryDate,
		}
		if opts.RecalculateLineStoreRejected != "" {
			recalculateLine.StoreRejected = &opts.RecalculateLineStoreRejected
		}
		if opts.RecalculateLineCarrierReject != "" {
			recalculateLine.CarrierRejected = &opts.RecalculateLineCarrierReject
		}
	}

	trace := &TraceByInput{
		RequestID:        opts.RequestID,
		RequestType:      opts.RequestType,
		Endpoint:         opts.Endpoint,
		ReceivedAt:       opts.ReceivedAt,
		EnterpriseCode:   opts.EnterpriseCode,
		Cp:               opts.Cp,
		Channel:          opts.Channel,
		EddLines:         eddLine,
		RecalculateLines: recalculateLine,
		LineCount:        opts.LineCount,
		Tags:             opts.Tags,
		AdditionalData:   opts.AdditionalData,
		IngestedAt:       opts.IngestedAt,
	}
	return l.sendTraceAll(trace)
}

// SendTraceByOutput is kept for backward compatibility and delegates to SendTraceAll.
func (l *EddLogger) SendTraceByOutput(opts *TraceOutputOptions) (string, error) {
	if opts == nil {
		opts = &TraceOutputOptions{}
	}

	var metadata *OutputMetadata
	if opts.MetadataIDTxn != "" ||
		opts.MetadataProcessingTimeMs != nil ||
		opts.MetadataIngestedAt != "" ||
		opts.MetadataRecalculateOrder != "" {
		metadata = &OutputMetadata{
			ProcessingTimeMs: opts.MetadataProcessingTimeMs,
		}
		if opts.MetadataIDTxn != "" {
			metadata.IDTxn = &opts.MetadataIDTxn
		}
		if opts.MetadataIngestedAt != "" {
			metadata.IngestedAt = &opts.MetadataIngestedAt
		}
		if opts.MetadataRecalculateOrder != "" {
			metadata.RecalculateOrder = &opts.MetadataRecalculateOrder
		}
	}

	var algorithmModelState *string
	if opts.AlgorithmModelState != "" {
		algorithmModelState = &opts.AlgorithmModelState
	}

	var algorithmWeights *AlgorithmWeights
	if opts.AlgorithmWeightsInventory != nil ||
		opts.AlgorithmWeightsLeadTime != nil ||
		opts.AlgorithmWeightsCost != nil ||
		opts.AlgorithmWeightsNode != nil ||
		opts.AlgorithmWeightsPath != nil ||
		opts.AlgorithmWeightsDifference != nil ||
		opts.AlgorithmWeightsSplits != nil {
		algorithmWeights = &AlgorithmWeights{
			Inventory:  opts.AlgorithmWeightsInventory,
			LeadTime:   opts.AlgorithmWeightsLeadTime,
			Cost:       opts.AlgorithmWeightsCost,
			Node:       opts.AlgorithmWeightsNode,
			Path:       opts.AlgorithmWeightsPath,
			Difference: opts.AlgorithmWeightsDifference,
			Splits:     opts.AlgorithmWeightsSplits,
		}
	}

	var errorCode *string
	if opts.ErrorCode != "" {
		errorCode = &opts.ErrorCode
	}

	var errorMessage *string
	if opts.ErrorMessage != "" {
		errorMessage = &opts.ErrorMessage
	}

	trace := &TraceByOutput{
		RequestID:           opts.RequestID,
		RequestType:         opts.RequestType,
		Endpoint:            opts.Endpoint,
		RespondedAt:         opts.RespondedAt,
		HTTPStatusCode:      opts.HTTPStatusCode,
		StatusFamily:        opts.StatusFamily,
		IsError:             opts.IsError,
		Metadata:            metadata,
		AlgorithmModelState: algorithmModelState,
		AlgorithmWeights:    algorithmWeights,
		EddCalculated:       opts.EddCalculated,
		StoreIDs:            opts.StoreIDs,
		ErrorCode:           errorCode,
		ErrorMessage:        errorMessage,
		Tags:                opts.Tags,
		AdditionalData:      opts.AdditionalData,
		IngestedAt:          opts.IngestedAt,
	}
	return l.sendTraceAll(trace)
}

// SendTraceByLog is kept for backward compatibility and delegates to SendTraceAll.
func (l *EddLogger) SendTraceByLog(opts *TraceLogOptions) (string, error) {
	if opts == nil {
		opts = &TraceLogOptions{}
	}

	var level *string
	if opts.Level != "" {
		level = &opts.Level
	}

	var context *string
	if opts.Context != "" {
		context = &opts.Context
	}

	var message *string
	if opts.Message != "" {
		message = &opts.Message
	}

	var step *string
	if opts.Step != "" {
		step = &opts.Step
	}

	var stacktrace *string
	if opts.Stacktrace != "" {
		stacktrace = &opts.Stacktrace
	}

	var serviceName *string
	if opts.ServiceName != "" {
		serviceName = &opts.ServiceName
	}

	var request *LogRequestPayload
	if opts.RequestMethod != "" || opts.RequestBody != nil {
		request = &LogRequestPayload{
			Method: opts.RequestMethod,
			Body:   opts.RequestBody,
		}
	}

	var response *LogResponsePayload
	if opts.ResponseStatusCode != 0 || opts.ResponseBody != nil {
		response = &LogResponsePayload{
			StatusCode: opts.ResponseStatusCode,
			Body:       opts.ResponseBody,
		}
	}

	trace := &TraceByLog{
		LogID:          opts.LogID,
		RequestID:      opts.RequestID,
		RequestType:    opts.RequestType,
		Endpoint:       opts.Endpoint,
		LogAt:          opts.LogAt,
		Level:          level,
		Context:        context,
		Message:        message,
		Step:           step,
		DurationMs:     opts.DurationMs,
		IDTxn:          opts.IDTxn,
		Tags:           opts.Tags,
		AdditionalData: opts.AdditionalData,
		Extra:          opts.Extra,
		Stacktrace:     stacktrace,
		IngestedAt:     opts.IngestedAt,
		ServiceName:    serviceName,
		Request:        request,
		Response:       response,
	}
	return l.sendTraceAll(trace)
}
