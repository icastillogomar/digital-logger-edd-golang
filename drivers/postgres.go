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
	tableName = "LGS_EDD_SDK_HIS"
	ddl       = `
	CREATE TABLE IF NOT EXISTS LGS_EDD_SDK_HIS (
		id SERIAL PRIMARY KEY,
		traceId VARCHAR(255) NOT NULL,
		timeLocal TIMESTAMP NOT NULL,
		timeUTC TIMESTAMP NOT NULL,
		service VARCHAR(255) NOT NULL,
		level VARCHAR(50) NOT NULL,
		"user" VARCHAR(255),
		action VARCHAR(255),
		context VARCHAR(255),
		request JSONB,
		response JSONB,
		durationMs FLOAT,
		tags TEXT,
		messageInfo TEXT,
		messageRaw TEXT,
		flagSummary INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_lgs_edd_sdk_his_trace_id ON LGS_EDD_SDK_HIS(traceId);
	CREATE INDEX IF NOT EXISTS idx_lgs_edd_sdk_his_time_utc ON LGS_EDD_SDK_HIS(timeUTC);
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
		INSERT INTO LGS_EDD_SDK_HIS 
			(traceId, timeLocal, timeUTC, service, level, "user", action, context, 
			 request, response, durationMs, tags, messageInfo, messageRaw, flagSummary)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`

	now := time.Now()
	nowUTC := time.Now().UTC()

	var requestJSON, responseJSON *string
	if req, ok := record["request"]; ok && req != nil {
		data, _ := json.Marshal(req)
		str := string(data)
		requestJSON = &str
	}
	if resp, ok := record["response"]; ok && resp != nil {
		data, _ := json.Marshal(resp)
		str := string(data)
		responseJSON = &str
	}

	var tagsStr *string
	if tags, ok := record["tags"].([]string); ok && len(tags) > 0 {
		str := strings.Join(tags, ",")
		tagsStr = &str
	}

	var rowID int
	err := d.conn.QueryRow(
		sqlQuery,
		record["traceId"],
		now,
		nowUTC,
		record["service"],
		record["level"],
		record["user"],
		record["action"],
		record["context"],
		requestJSON,
		responseJSON,
		record["durationMs"],
		tagsStr,
		record["messageInfo"],
		record["messageRaw"],
		0,
	).Scan(&rowID)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", rowID), nil
}

func (d *PostgresDriver) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
