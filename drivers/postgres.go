package drivers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	dbURL    string
	conn     *sql.DB
	migrated bool
}

const (
	tableName = "LGS_EDD_IA_LOGS_HIS"
	ddl       = `
	CREATE TABLE IF NOT EXISTS LGS_EDD_IA_LOGS_HIS (
		id SERIAL PRIMARY KEY,
		logId VARCHAR(255),
		requestId VARCHAR(255),
		requestType VARCHAR(255),
		endpoint TEXT,
		logAt TIMESTAMP NOT NULL,
		level VARCHAR(50),
		context TEXT,
		message TEXT,
		step VARCHAR(255),
		durationMs DOUBLE PRECISION,
		idTxn VARCHAR(255),
		tags TEXT,
		additionalData JSONB,
		extra JSONB,
		stacktrace TEXT,
		ingestedAt TIMESTAMP NOT NULL,
		serviceName VARCHAR(255),
		requestMethod VARCHAR(50),
		requestBody JSONB,
		responseStatusCode INTEGER,
		responseBody JSONB
	);
	CREATE INDEX IF NOT EXISTS idx_lgs_edd_ia_logs_his_request_id ON LGS_EDD_IA_LOGS_HIS(requestId);
	CREATE INDEX IF NOT EXISTS idx_lgs_edd_ia_logs_his_log_at ON LGS_EDD_IA_LOGS_HIS(logAt);
	`
)

func NewPostgresDriver(dbURL string) (*PostgresDriver, error) {
	if dbURL == "" {
		dbURL = os.Getenv("DB_URL")
	}
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL no está configurado")
	}
	return &PostgresDriver{dbURL: dbURL}, nil
}

func (d *PostgresDriver) ensureConnection() error {
	if d.conn != nil {
		return nil
	}
	conn, err := sql.Open("postgres", d.dbURL)
	if err != nil {
		return err
	}
	d.conn = conn
	fmt.Println("[digital-edd-logger] Conectado a PostgreSQL")
	return nil
}

func (d *PostgresDriver) ensureTable() error {
	if d.migrated {
		return nil
	}
	if err := d.ensureConnection(); err != nil {
		return err
	}
	if _, err := d.conn.Exec(ddl); err != nil {
		return err
	}
	d.migrated = true
	fmt.Printf("[digital-edd-logger] Tabla %s verificada/creada\n", tableName)
	return nil
}

func (d *PostgresDriver) Send(record map[string]interface{}) (string, error) {
	if err := d.ensureTable(); err != nil {
		return "", err
	}

	sqlQuery := `
		INSERT INTO LGS_EDD_IA_LOGS_HIS
			(logId, requestId, requestType, endpoint, logAt, level, context, message,
			 step, durationMs, idTxn, tags, additionalData, extra, stacktrace, ingestedAt,
			 serviceName, requestMethod, requestBody, responseStatusCode, responseBody)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id
	`

	requestMap, _ := record["request"].(map[string]interface{})
	responseMap, _ := record["response"].(map[string]interface{})

	requestMethod, _ := requestMap["method"].(string)
	requestBodyJSON := marshalNullableJSON(requestMap["body"])
	responseStatusCode := toNullableInt(responseMap["statusCode"])
	responseBodyJSON := marshalNullableJSON(responseMap["body"])

	logAt, err := parseTimestampWithDefault(record["logAt"])
	if err != nil {
		return "", fmt.Errorf("invalid logAt: %w", err)
	}
	ingestedAt, err := parseTimestampWithDefault(record["ingestedAt"])
	if err != nil {
		return "", fmt.Errorf("invalid ingestedAt: %w", err)
	}

	var tagsStr *string
	if tags := normalizeTags(record["tags"]); len(tags) > 0 {
		str := strings.Join(tags, ",")
		tagsStr = &str
	}

	var rowID int
	err = d.conn.QueryRow(
		sqlQuery,
		record["logId"],
		record["requestId"],
		record["requestType"],
		record["endpoint"],
		logAt,
		record["level"],
		record["context"],
		record["message"],
		record["step"],
		record["durationMs"],
		record["idTxn"],
		tagsStr,
		marshalNullableJSON(record["additionalData"]),
		marshalNullableJSON(record["extra"]),
		record["stacktrace"],
		ingestedAt,
		record["serviceName"],
		nullableString(requestMethod),
		requestBodyJSON,
		responseStatusCode,
		responseBodyJSON,
	).Scan(&rowID)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", rowID), nil
}

func marshalNullableJSON(value interface{}) *string {
	if value == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil || string(data) == "null" {
		return nil
	}

	str := string(data)
	return &str
}

func normalizeTags(value interface{}) []string {
	switch tags := value.(type) {
	case []string:
		return tags
	case []interface{}:
		result := make([]string, 0, len(tags))
		for _, tag := range tags {
			if str, ok := tag.(string); ok && str != "" {
				result = append(result, str)
			}
		}
		return result
	default:
		return nil
	}
}

func toNullableInt(value interface{}) interface{} {
	switch v := value.(type) {
	case nil:
		return nil
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return nil
	}
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func (d *PostgresDriver) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

func parseTimestampWithDefault(value interface{}) (time.Time, error) {
	now := time.Now()

	switch v := value.(type) {
	case nil:
		return now, nil
	case time.Time:
		if v.IsZero() {
			return now, nil
		}
		return v, nil
	case string:
		if strings.TrimSpace(v) == "" {
			return now, nil
		}
		for _, layout := range []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02 15:04:05.999999",
		} {
			if parsed, err := time.Parse(layout, v); err == nil {
				return parsed, nil
			}
		}
		return time.Time{}, fmt.Errorf("unsupported timestamp format %q", v)
	default:
		return time.Time{}, fmt.Errorf("unsupported timestamp type %T", value)
	}
}
