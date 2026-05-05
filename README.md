# Digital EDD Logger (GO)

SDK de logging para servicios Go con soporte para PostgreSQL (desarrollo) y Google Cloud PubSub (producción).

## Instalación.

```bash
go get github.com/icastillogomar/digital-logger-edd-golang
```

## Uso Rápido

```go
package main

import (
    eddlogger "github.com/icastillogomar/digital-logger-edd-golang"
)

func main() {
	traceLogger := eddlogger.NewLogger("my-service")
    defer traceLogger.Close()

	traceLogger.Log(&eddlogger.LogOptions{
		TraceID: "bd24e7ad-2e41-4638-b129-c1dd7e125faa",
		Level: "INFO",
		Action: "",
		Context: "middleware.request_response",
		Method: "HTTP",
		Path: "/edd/fee2/facade/sbb",
		RequestBody: map[string]interface{}{},
		StatusCode: 200,
		ResponseBody: map[string]interface{}{"success":true},
		MessageInfo: "",
		MessageRaw: "",
		DurationMs: 0,
		Tags: []string{"http", "middleware", "request-response"},
		Service: "my-service",
    })
}
```

## Configuración

### Local/Dev (PostgreSQL)

```bash
DB_URL=postgresql://user:password@localhost:5432/mydb
ENV=local
```

### Producción/QA (PubSub)

```bash
ENV=prod  # o "production", "qa", "qas"
GOOGLE_CLOUD_PROJECT=my-project-id
```

## Comportamiento

| ENV                               | Driver     | Destino                 |
|-----------------------------------|------------|-------------------------|
| `local` (o vacío)                 | PostgreSQL | Tabla `LGS_EDD_SDK_HIS` |
| `prod`, `production`, `qa`, `qas` | PubSub     | Topic `digital-edd-sdk` |

Si falta configuración, usa `ConsoleDriver` como fallback.

## API

```go
type TraceLogOptions struct {
    LogID              string
    RequestID          string
    RequestType        string
    Endpoint           string
    LogAt              string
    Level              string    // DEBUG, INFO, WARNING, ERROR, CRITICAL
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
```

## Variables de Entorno

| Variable               | Descripción                    | Requerido                             |
|------------------------|--------------------------------|---------------------------------------|
| `DB_URL`               | URL de PostgreSQL              | Solo en local                         |
| `ENV`                  | `local` para forzar PostgreSQL | Opcional                              |
| `GOOGLE_CLOUD_PROJECT` | Project ID de GCP              | Solo en prod                          |
| `SDKTRACKING_PUBLISH`  | `false` para deshabilitar      | Opcional                              |
| `PUBSUB_TOPIC_NAME`    | Nombre del topic               | Opcional (default: `digital-edd-sdk`) |
