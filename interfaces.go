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
	Method  string                 `json:"method"`
	Path    string                 `json:"path"`
	Headers map[string]string      `json:"headers,omitempty"`
	Body    interface{}            `json:"body,omitempty"`
}

type ResponseInfo struct {
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]string      `json:"headers,omitempty"`
	Body       interface{}            `json:"body,omitempty"`
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
