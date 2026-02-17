# Digital EDD Logger (Go)

SDK de logging para servicios Go con soporte para PostgreSQL (desarrollo) y Google Cloud PubSub (producción).

## Instalación

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
    log := eddlogger.NewLogger("my-service")
    defer log.Close()
    
    log.Log(&eddlogger.LogOptions{
        TraceID:      "abc-123",
        Action:       "ORDER_CREATED",
        Context:      "OrderService",
        Method:       "POST",
        Path:         "/api/orders",
        RequestBody:  map[string]interface{}{"product": "ABC", "qty": 2},
        StatusCode:   200,
        ResponseBody: map[string]interface{}{"order_id": "12345"},
        DurationMs:   150.5,
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

| ENV | Driver | Destino |
|-----|--------|---------|
| `local` (o vacío) | PostgreSQL | Tabla `LGS_EDD_SDK_HIS` |
| `prod`, `production`, `qa`, `qas` | PubSub | Topic `digital-edd-sdk` |

Si falta configuración, usa `ConsoleDriver` como fallback.

## API

```go
type LogOptions struct {
    TraceID         string
    Level           string                 // DEBUG, INFO, WARNING, ERROR, CRITICAL
    Action          string
    Context         string
    Method          string
    Path            string
    RequestHeaders  map[string]string
    RequestBody     interface{}
    StatusCode      int
    ResponseHeaders map[string]string
    ResponseBody    interface{}
    DurationMs      float64
    MessageInfo     string
    MessageRaw      string
    Tags            []string
    Service         string
}
```

## Variables de Entorno

| Variable | Descripción | Requerido |
|----------|-------------|-----------|
| `DB_URL` | URL de PostgreSQL | Solo en local |
| `ENV` | `local` para forzar PostgreSQL | Opcional |
| `GOOGLE_CLOUD_PROJECT` | Project ID de GCP | Solo en prod |
| `SDKTRACKING_PUBLISH` | `false` para deshabilitar | Opcional |
| `PUBSUB_TOPIC_NAME` | Nombre del topic | Opcional (default: `digital-edd-sdk`) |
